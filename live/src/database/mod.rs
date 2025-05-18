use anyhow::{Context, Result};
use rusqlite::{Connection, params};
use tracing::{info, debug};
use std::env;

use crate::models::{Cache, Post, Project, Tag, PostToTag, PostToPost, PostToProject, ProjectToTag, ProjectToProject, TagToTag};

pub const DB_NAME: &str = "master.db";

// Connection string for the database with optimized settings
pub fn get_db_connection_string() -> String {
    // Get database name from environment or use default
    let db_name = std::env::var("DB_NAME").unwrap_or_else(|_| DB_NAME.to_string());
    
    // Default SQLite connection string with pragmas for better performance
    format!(
        "file:{}?_pragma=busy_timeout=5000&_pragma=journal_mode=WAL&_pragma=mmap_size=30000000000&_pragma=page_size=32768", 
        db_name
    )
}

/// Database connection wrapper
pub struct Database {
    conn: Connection,
}

// Implement Send and Sync for Database
// This is safe because we ensure all database operations are serialized
unsafe impl Send for Database {}
unsafe impl Sync for Database {}

impl Database {
    /// Get the database path for use by workers
    pub fn get_db_path() -> String {
        env::var("DB_NAME").unwrap_or_else(|_| DB_NAME.to_string())
    }

    /// Create a new database connection with optimized settings
    pub fn new() -> Result<Self> {
        // Get the database path (not the connection string)
        let db_path = Self::get_db_path();
        
        debug!("Opening database: {}", db_path);
        let conn = Connection::open(&db_path)
            .with_context(|| format!("Failed to open database: {}", db_path))?;

        // Use a simplified approach to avoid pragma setting issues
        debug!("Setting pragmas for database performance");
        
        // Just execute each pragma directly - don't try to get the return value
        // This is the safest approach, even if it doesn't provide feedback
        conn.execute_batch("
            PRAGMA journal_mode = WAL;
            PRAGMA synchronous = NORMAL;
            PRAGMA busy_timeout = 5000;
            PRAGMA cache_size = 10000;
            PRAGMA temp_store = MEMORY;
            PRAGMA mmap_size = 268435456;
        ").context("Failed to set database pragmas")?;
        
        debug!("Database pragmas set successfully");

        Ok(Self { conn })
    }
    
    /// Count the number of cache entries (for testing)
    pub fn count_cache_entries(&self) -> Result<usize> {
        let count: i64 = self.conn.query_row(
            "SELECT COUNT(*) FROM caches",
            [],
            |row| row.get(0)
        )?;
        
        Ok(count as usize)
    }
    
    /// Count the number of entities in a table (for testing)
    pub fn count_entities(&self, table_name: &str) -> Result<usize> {
        // Validate table name to prevent SQL injection
        let allowed_tables = ["posts", "projects", "tags"];
        if !allowed_tables.contains(&table_name) {
            return Err(anyhow::anyhow!("Invalid table name: {}", table_name));
        }
        
        let count: i64 = self.conn.query_row(
            &format!("SELECT COUNT(*) FROM {}", table_name),
            [],
            |row| row.get(0)
        )?;
        
        Ok(count as usize)
    }

    /// Initialize the database schema
    pub fn init_db(&self) -> Result<()> {
        info!("Initializing database schema");
        self.create_tables()?;
        Ok(())
    }

    /// Create all required database tables
    fn create_tables(&self) -> Result<()> {
        // Create relationship tables
        self.create_post_to_tag_table()?;
        self.create_post_to_post_table()?;
        self.create_post_to_project_table()?;
        self.create_project_to_tag_table()?;
        self.create_project_to_project_table()?;
        self.create_tag_to_tag_table()?;

        // Create main entity tables
        self.create_post_table()?;
        self.create_tag_table()?;
        self.create_project_table()?;
        self.create_cache_table()?;

        info!("All tables created successfully");
        Ok(())
    }

    // Main entity tables
    fn create_post_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS posts (
                id INTEGER PRIMARY KEY,
                title TEXT NOT NULL,
                slug TEXT UNIQUE NOT NULL,
                description TEXT NOT NULL,
                content TEXT NOT NULL,
                banner_path TEXT,
                created_at INTEGER NOT NULL,
                x REAL,
                y REAL,
                z REAL
            )",
            [],
        )?;
        Ok(())
    }

    fn create_project_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS projects (
                id INTEGER PRIMARY KEY,
                title TEXT NOT NULL,
                slug TEXT UNIQUE NOT NULL,
                description TEXT NOT NULL,
                content TEXT NOT NULL,
                banner_path TEXT,
                created_at INTEGER NOT NULL,
                x REAL,
                y REAL,
                z REAL
            )",
            [],
        )?;
        Ok(())
    }

    fn create_tag_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS tags (
                id INTEGER PRIMARY KEY,
                title TEXT NOT NULL,
                slug TEXT UNIQUE NOT NULL,
                description TEXT NOT NULL,
                content TEXT NOT NULL,
                banner_path TEXT,
                icon TEXT,
                created_at INTEGER NOT NULL,
                x REAL,
                y REAL,
                z REAL
            )",
            [],
        )?;
        Ok(())
    }

    fn create_cache_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS caches (
                id INTEGER PRIMARY KEY,
                path TEXT UNIQUE NOT NULL,
                hashed TEXT NOT NULL,
                x REAL,
                y REAL,
                z REAL
            )",
            [],
        )?;
        Ok(())
    }

    // Relationship tables
    fn create_post_to_tag_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS post_to_tags (
                post_id INTEGER NOT NULL,
                tag_id INTEGER NOT NULL,
                PRIMARY KEY (post_id, tag_id)
            )",
            [],
        )?;
        Ok(())
    }

    fn create_post_to_post_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS post_to_posts (
                source_post_id INTEGER NOT NULL,
                target_post_id INTEGER NOT NULL,
                PRIMARY KEY (source_post_id, target_post_id)
            )",
            [],
        )?;
        Ok(())
    }

    fn create_post_to_project_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS post_to_projects (
                post_id INTEGER NOT NULL,
                project_id INTEGER NOT NULL,
                PRIMARY KEY (post_id, project_id)
            )",
            [],
        )?;
        Ok(())
    }

    fn create_project_to_tag_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS project_to_tags (
                project_id INTEGER NOT NULL,
                tag_id INTEGER NOT NULL,
                PRIMARY KEY (project_id, tag_id)
            )",
            [],
        )?;
        Ok(())
    }

    fn create_project_to_project_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS project_to_projects (
                source_project_id INTEGER NOT NULL,
                target_project_id INTEGER NOT NULL,
                PRIMARY KEY (source_project_id, target_project_id)
            )",
            [],
        )?;
        Ok(())
    }

    fn create_tag_to_tag_table(&self) -> Result<()> {
        self.conn.execute(
            "CREATE TABLE IF NOT EXISTS tag_to_tags (
                source_tag_id INTEGER NOT NULL,
                target_tag_id INTEGER NOT NULL,
                PRIMARY KEY (source_tag_id, target_tag_id)
            )",
            [],
        )?;
        Ok(())
    }

    // Cache operations
    pub fn get_cache(&self, path: &str) -> Result<Option<Cache>> {
        let mut stmt = self.conn.prepare("SELECT id, path, hashed, x, y, z FROM caches WHERE path = ?")?;
        let mut rows = stmt.query(params![path])?;
        
        if let Some(row) = rows.next()? {
            let cache = Cache {
                id: Some(row.get(0)?),
                path: row.get(1)?,
                hash: row.get(2)?,
                x: row.get(3)?,
                y: row.get(4)?,
                z: row.get(5)?,
            };
            Ok(Some(cache))
        } else {
            Ok(None)
        }
    }

    pub fn update_cache(&self, cache: &Cache) -> Result<()> {
        self.conn.execute(
            "INSERT INTO caches (path, hashed, x, y, z)
             VALUES (?1, ?2, ?3, ?4, ?5)
             ON CONFLICT (path) DO UPDATE SET
               hashed = excluded.hashed,
               x = excluded.x,
               y = excluded.y,
               z = excluded.z",
            params![cache.path, cache.hash, cache.x, cache.y, cache.z],
        )?;
        Ok(())
    }

    // Post operations
    pub fn find_post_by_slug(&self, slug: &str) -> Result<Option<Post>> {
        let mut stmt = self.conn.prepare(
            "SELECT id, title, slug, description, content, banner_path, created_at, x, y, z 
             FROM posts WHERE slug = ?"
        )?;
        
        let mut rows = stmt.query(params![slug])?;
        
        if let Some(row) = rows.next()? {
            let post = Post {
                id: Some(row.get(0)?),
                title: row.get(1)?,
                slug: row.get(2)?,
                description: row.get(3)?,
                content: row.get(4)?,
                banner_path: row.get(5)?,
                created_at: row.get(6)?,
                tags: Vec::new(),
                posts: Vec::new(),
                projects: Vec::new(),
                x: row.get(7)?,
                y: row.get(8)?,
                z: row.get(9)?,
            };
            Ok(Some(post))
        } else {
            Ok(None)
        }
    }

    pub fn save_post(&self, post: &Post) -> Result<i64> {
        let result = self.conn.execute(
            "INSERT INTO posts (title, slug, description, content, banner_path, created_at, x, y, z)
             VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)
             ON CONFLICT (slug) DO UPDATE SET
               title = excluded.title,
               description = excluded.description,
               content = excluded.content,
               banner_path = excluded.banner_path,
               created_at = excluded.created_at,
               x = excluded.x,
               y = excluded.y,
               z = excluded.z",
            params![
                post.title, post.slug, post.description, post.content, 
                post.banner_path, post.created_at, post.x, post.y, post.z
            ],
        )?;
        
        let id = if result > 0 {
            self.conn.last_insert_rowid()
        } else {
            // If the row was updated rather than inserted, get the ID of the existing row
            let mut stmt = self.conn.prepare("SELECT id FROM posts WHERE slug = ?")?;
            let id: i64 = stmt.query_row(params![post.slug], |row| row.get(0))?;
            id
        };
        
        Ok(id)
    }

    // Project operations
    pub fn find_project_by_slug(&self, slug: &str) -> Result<Option<Project>> {
        let mut stmt = self.conn.prepare(
            "SELECT id, title, slug, description, content, banner_path, created_at, x, y, z 
             FROM projects WHERE slug = ?"
        )?;
        
        let mut rows = stmt.query(params![slug])?;
        
        if let Some(row) = rows.next()? {
            let project = Project {
                id: Some(row.get(0)?),
                title: row.get(1)?,
                slug: row.get(2)?,
                description: row.get(3)?,
                content: row.get(4)?,
                banner_path: row.get(5)?,
                created_at: row.get(6)?,
                tags: Vec::new(),
                posts: Vec::new(),
                projects: Vec::new(),
                x: row.get(7)?,
                y: row.get(8)?,
                z: row.get(9)?,
            };
            Ok(Some(project))
        } else {
            Ok(None)
        }
    }

    pub fn save_project(&self, project: &Project) -> Result<i64> {
        let result = self.conn.execute(
            "INSERT INTO projects (title, slug, description, content, banner_path, created_at, x, y, z)
             VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)
             ON CONFLICT (slug) DO UPDATE SET
               title = excluded.title,
               description = excluded.description,
               content = excluded.content,
               banner_path = excluded.banner_path,
               created_at = excluded.created_at,
               x = excluded.x,
               y = excluded.y,
               z = excluded.z",
            params![
                project.title, project.slug, project.description, project.content, 
                project.banner_path, project.created_at, project.x, project.y, project.z
            ],
        )?;
        
        let id = if result > 0 {
            self.conn.last_insert_rowid()
        } else {
            let mut stmt = self.conn.prepare("SELECT id FROM projects WHERE slug = ?")?;
            let id: i64 = stmt.query_row(params![project.slug], |row| row.get(0))?;
            id
        };
        
        Ok(id)
    }

    // Tag operations
    pub fn find_tag_by_slug(&self, slug: &str) -> Result<Option<Tag>> {
        let mut stmt = self.conn.prepare(
            "SELECT id, title, slug, description, content, banner_path, icon, created_at, x, y, z 
             FROM tags WHERE slug = ?"
        )?;
        
        let mut rows = stmt.query(params![slug])?;
        
        if let Some(row) = rows.next()? {
            let tag = Tag {
                id: Some(row.get(0)?),
                title: row.get(1)?,
                slug: row.get(2)?,
                description: row.get(3)?,
                content: row.get(4)?,
                banner_path: row.get(5)?,
                icon: row.get(6)?,
                created_at: row.get(7)?,
                tags: Vec::new(),
                posts: Vec::new(),
                projects: Vec::new(),
                x: row.get(8)?,
                y: row.get(9)?,
                z: row.get(10)?,
            };
            Ok(Some(tag))
        } else {
            Ok(None)
        }
    }

    pub fn save_tag(&self, tag: &Tag) -> Result<i64> {
        let result = self.conn.execute(
            "INSERT INTO tags (title, slug, description, content, banner_path, icon, created_at, x, y, z)
             VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10)
             ON CONFLICT (slug) DO UPDATE SET
               title = excluded.title,
               description = excluded.description,
               content = excluded.content,
               banner_path = excluded.banner_path,
               icon = excluded.icon,
               created_at = excluded.created_at,
               x = excluded.x,
               y = excluded.y,
               z = excluded.z",
            params![
                tag.title, tag.slug, tag.description, tag.content, 
                tag.banner_path, tag.icon, tag.created_at, tag.x, tag.y, tag.z
            ],
        )?;
        
        let id = if result > 0 {
            self.conn.last_insert_rowid()
        } else {
            let mut stmt = self.conn.prepare("SELECT id FROM tags WHERE slug = ?")?;
            let id: i64 = stmt.query_row(params![tag.slug], |row| row.get(0))?;
            id
        };
        
        Ok(id)
    }

    // Relationship operations
    pub fn save_post_to_tag(&self, rel: &PostToTag) -> Result<()> {
        self.conn.execute(
            "INSERT OR IGNORE INTO post_to_tags (post_id, tag_id) VALUES (?1, ?2)",
            params![rel.post_id, rel.tag_id],
        )?;
        Ok(())
    }

    pub fn save_post_to_post(&self, rel: &PostToPost) -> Result<()> {
        self.conn.execute(
            "INSERT OR IGNORE INTO post_to_posts (source_post_id, target_post_id) VALUES (?1, ?2)",
            params![rel.source_post_id, rel.target_post_id],
        )?;
        Ok(())
    }

    pub fn save_post_to_project(&self, rel: &PostToProject) -> Result<()> {
        self.conn.execute(
            "INSERT OR IGNORE INTO post_to_projects (post_id, project_id) VALUES (?1, ?2)",
            params![rel.post_id, rel.project_id],
        )?;
        Ok(())
    }

    pub fn save_project_to_tag(&self, rel: &ProjectToTag) -> Result<()> {
        self.conn.execute(
            "INSERT OR IGNORE INTO project_to_tags (project_id, tag_id) VALUES (?1, ?2)",
            params![rel.project_id, rel.tag_id],
        )?;
        Ok(())
    }

    pub fn save_project_to_project(&self, rel: &ProjectToProject) -> Result<()> {
        self.conn.execute(
            "INSERT OR IGNORE INTO project_to_projects (source_project_id, target_project_id) VALUES (?1, ?2)",
            params![rel.source_project_id, rel.target_project_id],
        )?;
        Ok(())
    }

    pub fn save_tag_to_tag(&self, rel: &TagToTag) -> Result<()> {
        self.conn.execute(
            "INSERT OR IGNORE INTO tag_to_tags (source_tag_id, target_tag_id) VALUES (?1, ?2)",
            params![rel.source_tag_id, rel.target_tag_id],
        )?;
        Ok(())
    }
}
