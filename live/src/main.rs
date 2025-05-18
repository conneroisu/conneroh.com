mod database;
mod models;
mod processor;

use std::env;
use std::path::Path;
use std::process::exit;
use std::sync::Arc;
use anyhow::{Context, Result};
use tokio::sync::oneshot;
use tokio::time::Duration;
use tracing::{info, error, warn, debug, Level};
use tracing_subscriber::FmtSubscriber;

use crate::database::Database;
use crate::models::VAULT_LOC;
use crate::processor::Processor;

/// Command line arguments for the program
struct Args {
    workers: usize,
    vault_path: String,
    timeout_seconds: u64,
    bucket_name: String,
    debug: bool,
}

impl Default for Args {
    fn default() -> Self {
        Self {
            workers: 8,
            vault_path: VAULT_LOC.to_string(),
            timeout_seconds: 300, // 5 minutes
            bucket_name: "conneroh-com".to_string(),
            debug: false,
        }
    }
}

/// Parse command line arguments
fn parse_args() -> Args {
    let mut args = Args::default();
    
    let args_vec: Vec<String> = env::args().collect();
    
    // Parse command-line arguments
    for i in 1..args_vec.len() {
        match args_vec[i].as_str() {
            "--workers" | "-w" => {
                if i + 1 < args_vec.len() {
                    if let Ok(workers) = args_vec[i + 1].parse::<usize>() {
                        args.workers = workers;
                    }
                }
            },
            "--path" | "-p" => {
                if i + 1 < args_vec.len() {
                    args.vault_path = args_vec[i + 1].clone();
                }
            },
            "--timeout" | "-t" => {
                if i + 1 < args_vec.len() {
                    if let Ok(timeout) = args_vec[i + 1].parse::<u64>() {
                        args.timeout_seconds = timeout;
                    }
                }
            },
            "--bucket" | "-b" => {
                if i + 1 < args_vec.len() {
                    args.bucket_name = args_vec[i + 1].clone();
                }
            },
            "--debug" | "-d" => {
                args.debug = true;
            },
            "--help" | "-h" => {
                print_help();
                exit(0);
            },
            _ => {}
        }
    }
    
    // Check environment variables for values not set by command line
    if let Ok(bucket) = env::var("BUCKET_NAME") {
        args.bucket_name = bucket;
    }
    
    if let Ok(workers) = env::var("WORKERS").and_then(|s| s.parse::<usize>().map_err(|_| std::env::VarError::NotPresent)) {
        args.workers = workers;
    }
    
    if let Ok(debug) = env::var("DEBUG") {
        args.debug = debug == "true" || debug == "1";
    }
    
    args
}

/// Print help message
fn print_help() {
    println!("Usage: live [OPTIONS]");
    println!("Initialize the database from markdown files and assets");
    println!();
    println!("Options:");
    println!("  -w, --workers NUM    Number of worker threads (default: 8)");
    println!("  -p, --path PATH      Path to the vault directory (default: {})", VAULT_LOC);
    println!("  -t, --timeout SECS   Timeout in seconds (default: 300)");
    println!("  -b, --bucket NAME    S3 bucket name (default: conneroh-com)");
    println!("  -d, --debug          Enable debug logging");
    println!("  -h, --help           Print this help message");
    println!();
    println!("Environment variables:");
    println!("  BUCKET_NAME          S3 bucket name");
    println!("  WORKERS              Number of worker threads");
    println!("  DEBUG                Enable debug logging (set to 'true' or '1')");
}

/// Initialize the database from markdown files and assets.
///
/// This program scans the filesystem for markdown files and assets,
/// processes them, and updates the database.
#[tokio::main]
async fn main() -> Result<()> {
    // Parse command line arguments
    let args = parse_args();
    
    // Initialize logging
    let subscriber = FmtSubscriber::builder()
        .with_max_level(if args.debug { Level::DEBUG } else { Level::INFO })
        .finish();
    tracing::subscriber::set_global_default(subscriber)
        .context("Failed to set global default subscriber")?;

    info!("Starting database initialization");
    debug!("Arguments: workers={}, path={}, timeout={}s, bucket={}, debug={}", 
           args.workers, args.vault_path, args.timeout_seconds, args.bucket_name, args.debug);
    
    // Create shutdown channel
    let (shutdown_tx, shutdown_rx) = oneshot::channel::<()>();
    let shutdown_tx = Arc::new(tokio::sync::Mutex::new(Some(shutdown_tx)));
    
    // Set up Ctrl+C handler
    let shutdown_tx_clone = shutdown_tx.clone();
    tokio::spawn(async move {
        if let Err(e) = tokio::signal::ctrl_c().await {
            error!("Failed to listen for Ctrl+C: {}", e);
            return;
        }
        
        info!("Received Ctrl+C, shutting down gracefully...");
        if let Some(tx) = shutdown_tx_clone.lock().await.take() {
            let _ = tx.send(());
        }
    });
    
    // Create database connection
    let db = Database::new()
        .context("Failed to create database connection")?;
    
    // Initialize database schema
    db.init_db()
        .context("Failed to initialize database schema")?;
    
    // Create processor
    let processor = Processor::new(db, &args.bucket_name);
    
    // Scan filesystem
    let vault_path = Path::new(&args.vault_path);
    if !vault_path.exists() {
        return Err(anyhow::anyhow!("Vault path does not exist: {}", args.vault_path));
    }
    
    // Scan the filesystem directly
    let scan_result = match processor.scan_fs(vault_path).await {
        Ok(_) => Ok(()),
        Err(e) => {
            error!("Error scanning filesystem: {}", e);
            Err(e)
        }
    };
    
    // Start processing with timeout
    let start = std::time::Instant::now();
    let process_timeout = Duration::from_secs(args.timeout_seconds);
    
    // Process with timeout - using a simplified approach to avoid type issues
    let process_future = processor.start(args.workers);
    
    // First check scan result
    if let Err(e) = scan_result {
        return Err(e).context("Filesystem scan error");
    }
    
    let process_result = tokio::select! {
        // Run the processor with timeout
        result = tokio::time::timeout(process_timeout, process_future) => {
            match result {
                Ok(inner_result) => inner_result,
                Err(_) => {
                    error!("Processing timed out after {:?}", process_timeout);
                    Err(anyhow::anyhow!("Processing timed out"))
                }
            }
        },
        // Shutdown if requested
        _ = shutdown_rx => {
            warn!("Received shutdown signal during processing");
            Ok(())
        }
    };
    
    // Handle errors
    if let Err(e) = process_result {
        error!("Error in processor: {}", e);
        return Err(e).context("Processor error");
    }
    
    let elapsed = start.elapsed();
    info!("Database initialization completed in {:?}", elapsed);
    
    Ok(())
}
