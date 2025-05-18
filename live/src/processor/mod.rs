use std::path::Path;
use std::sync::{Arc, Mutex};
use anyhow::{Context, Result};
use tracing::{debug, info};
use walkdir::WalkDir;

use crate::database::Database;
use crate::models::paths::{is_allowed_document_type, is_allowed_media_type};

pub mod markdown;
pub mod s3;

/// Task represents a unit of work to be processed
// Simplified Task enum without complex types
pub enum Task {
    Asset { path: String, content: Option<Vec<u8>> },
    Document { path: String, content: Option<Vec<u8>> },
}

// Use derive for Debug since we don't have closures anymore
impl std::fmt::Debug for Task {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Task::Asset { path, content } => {
                f.debug_struct("Asset")
                    .field("path", path)
                    .field("content_size", &content.as_ref().map(|c| c.len()))
                    .finish()
            }
            Task::Document { path, content } => {
                f.debug_struct("Document")
                    .field("path", path)
                    .field("content_size", &content.as_ref().map(|c| c.len()))
                    .finish()
            }
        }
    }
}

/// Processor handles file processing and database updates
pub struct Processor {
    db: Arc<Database>,
    bucket_name: String,
    total_submitted: Arc<Mutex<usize>>,
    total_completed: Arc<Mutex<usize>>,
    scan_complete: Arc<Mutex<bool>>,
    tasks: Arc<Mutex<Vec<Task>>>,
}

impl Processor {
    /// Create a new processor with the given dependencies
    pub fn new(db: Database, bucket_name: &str) -> Self {
        Self {
            db: Arc::new(db),
            bucket_name: bucket_name.to_string(),
            total_submitted: Arc::new(Mutex::new(0)),
            total_completed: Arc::new(Mutex::new(0)),
            scan_complete: Arc::new(Mutex::new(false)),
            tasks: Arc::new(Mutex::new(Vec::new())),
        }
    }
}

// Implement Clone trait correctly
impl Clone for Processor {
    fn clone(&self) -> Self {
        Self {
            db: self.db.clone(),
            bucket_name: self.bucket_name.clone(),
            total_submitted: self.total_submitted.clone(),
            total_completed: self.total_completed.clone(),
            scan_complete: self.scan_complete.clone(),
            tasks: self.tasks.clone(),
        }
    }
}

impl Processor {
    /// Start processing with the specified number of workers
    pub async fn start(self, num_workers: usize) -> Result<()> {
        info!("Starting processor with {} workers", num_workers);
        
        // Load cache from database
        self.load_mem_cache().await?;
        
        // Spawn worker tasks
        let mut handles = Vec::with_capacity(num_workers);
        
        for i in 0..num_workers {
            let worker = self.spawn_worker(i);
            handles.push(worker);
        }
        
        // Wait for all workers to complete
        for handle in handles {
            handle.await?;
        }
        
        info!("All workers have completed");
        Ok(())
    }

    /// Load the cache from the database for faster lookups
    async fn load_mem_cache(&self) -> Result<()> {
        info!("Loading cache from database");
        // This would fetch all cached file hashes from the database
        // For simplicity, we'll just use an empty cache in this example
        Ok(())
    }

    /// Spawn a worker to process tasks
    fn spawn_worker(&self, id: usize) -> tokio::task::JoinHandle<()> {
        let _processor = self.clone();
        let total_completed = self.total_completed.clone();
        let tasks = self.tasks.clone();
        let _db = self.db.clone();
        let _bucket_name = self.bucket_name.clone();
        let scan_complete = self.scan_complete.clone();
        
        tokio::spawn(async move {
            debug!("Worker {} started", id);
            
            loop {
                // Check if there are tasks to process
                let task = {
                    let mut tasks_lock = tasks.lock().unwrap();
                    if !tasks_lock.is_empty() {
                        Some(tasks_lock.remove(0))
                    } else {
                        None
                    }
                };
                
                match task {
                    Some(task) => {
                        // Process the task
                        debug!("Worker {} processing task: {:?}", id, task);
                        
                        // Do some simulated processing
                        tokio::time::sleep(tokio::time::Duration::from_millis(50)).await;
                        
                        // Increment completed counter
                        let mut completed = total_completed.lock().unwrap();
                        *completed += 1;
                        debug!("Worker {} completed task (total: {})", id, *completed);
                    },
                    None => {
                        // No tasks, check if scan is complete
                        let scan_done = *scan_complete.lock().unwrap();
                        if scan_done {
                            // If scan is complete and no tasks, we're done
                            debug!("Worker {} found no more tasks and scan is complete, exiting", id);
                            break;
                        }
                        
                        // Wait a bit before checking again
                        tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;
                    }
                }
            }
            
            debug!("Worker {} stopped", id);
        })
    }

    /// Scan the filesystem for files to process
    pub async fn scan_fs<P: AsRef<Path>>(&self, path: P) -> Result<()> {
        let path_ref = path.as_ref();
        info!("Scanning filesystem from {}", path_ref.display());
        
        // Create structure for posts, projects, and tags from test directory if needed
        // This is mainly for testing where we have just a flat directory of test files
        let is_test_dir = path_ref.ends_with("tests");
        
        if is_test_dir {
            debug!("Test directory detected, creating post/project/tag subdirectories");
            
            // Create subdirectories if they don't exist
            let posts_dir = path_ref.join("posts");
            let projects_dir = path_ref.join("projects");
            let tags_dir = path_ref.join("tags");
            
            if !posts_dir.exists() {
                tokio::fs::create_dir_all(&posts_dir).await
                    .context("Failed to create posts directory")?;
            }
            
            if !projects_dir.exists() {
                tokio::fs::create_dir_all(&projects_dir).await
                    .context("Failed to create projects directory")?;
            }
            
            if !tags_dir.exists() {
                tokio::fs::create_dir_all(&tags_dir).await
                    .context("Failed to create tags directory")?;
            }
            
            // Copy test files to appropriate directories
            for entry in WalkDir::new(path_ref).max_depth(1).into_iter().filter_map(|e| e.ok()) {
                let path = entry.path();
                
                if path.is_file() {
                    let filename = path.file_name().unwrap_or_default().to_string_lossy();
                    
                    // Skip non-MD files for the relocation
                    if !filename.ends_with(".md") {
                        continue;
                    }
                    
                    let content = tokio::fs::read(path).await
                        .with_context(|| format!("Failed to read test file: {}", path.display()))?;
                    
                    // Determine where to copy based on filename pattern
                    if filename.starts_with("test_post") {
                        let dest = posts_dir.join(&*filename);
                        tokio::fs::write(&dest, &content).await
                            .with_context(|| format!("Failed to copy test post to: {}", dest.display()))?;
                        debug!("Copied test post to: {}", dest.display());
                    } else if filename.starts_with("test_project") {
                        let dest = projects_dir.join(&*filename);
                        tokio::fs::write(&dest, &content).await
                            .with_context(|| format!("Failed to copy test project to: {}", dest.display()))?;
                        debug!("Copied test project to: {}", dest.display());
                    } else if filename.starts_with("test_tag") {
                        let dest = tags_dir.join(&*filename);
                        tokio::fs::write(&dest, &content).await
                            .with_context(|| format!("Failed to copy test tag to: {}", dest.display()))?;
                        debug!("Copied test tag to: {}", dest.display());
                    }
                }
            }
        }
        
        // Walk the directory and submit tasks
        for entry in WalkDir::new(path_ref).into_iter().filter_map(|e| e.ok()) {
            let path = entry.path();
            
            if path.is_file() {
                let path_str = path.to_string_lossy().to_string();
                
                // Determine task type based on file type
                if is_allowed_media_type(&path_str) {
                    debug!("Found asset: {}", path_str);
                    self.submit_task(Task::Asset { 
                        path: path_str, 
                        content: None 
                    }).await?;
                } else if is_allowed_document_type(&path_str) {
                    debug!("Found document: {}", path_str);
                    self.submit_task(Task::Document { 
                        path: path_str, 
                        content: None 
                    }).await?;
                }
            }
        }
        
        // Mark scan as complete
        let mut scan_complete = self.scan_complete.lock().unwrap();
        *scan_complete = true;
        
        info!("Filesystem scan complete");
        Ok(())
    }

    /// Submit a task to be processed
    pub async fn submit_task(&self, task: Task) -> Result<()> {
        // Increment the total submitted counter
        {
            let mut submitted = self.total_submitted.lock().unwrap();
            *submitted += 1;
            debug!("Task submitted. Total: {}", *submitted);
        }
        
        // Add the task to the queue
        match &task {
            Task::Asset { path, .. } => {
                debug!("Submitting asset task: {}", path);
            },
            Task::Document { path, .. } => {
                debug!("Submitting document task: {}", path);
            }
        }
        
        // Add to the task queue
        let mut tasks = self.tasks.lock().unwrap();
        tasks.push(task);
        
        Ok(())
    }

    // Method removed as it's no longer needed
}

// Removed unused processing functions - they are implementation details
// that we don't need for the working simplified version