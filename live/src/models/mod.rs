use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::fmt;

pub mod paths;

/// Constants for document locations
pub const VAULT_LOC: &str = "internal/data/docs/";
pub const ASSETS_LOC: &str = "assets/";
pub const POSTS_LOC: &str = "posts/";
pub const TAGS_LOC: &str = "tags/";
pub const PROJECTS_LOC: &str = "projects/";

/// CustomTime allows for flexible YAML time parsing
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CustomTime(pub DateTime<Utc>);

impl Default for CustomTime {
    fn default() -> Self {
        Self(Utc::now())
    }
}

impl fmt::Display for CustomTime {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0.format("%Y-%m-%d %H:%M:%S"))
    }
}

/// Base document containing common fields for all document types
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Doc {
    pub title: String,
    #[serde(skip)]
    pub path: String,
    pub slug: String,
    pub description: String,
    #[serde(skip)]
    pub content: String,
    pub banner_path: String,
    pub icon: String,
    pub created_at: CustomTime,
    pub updated_at: CustomTime,
    #[serde(default)]
    pub tags: Vec<String>,
    #[serde(default)]
    pub posts: Vec<String>,
    #[serde(default)]
    pub projects: Vec<String>,
    #[serde(skip)]
    pub hash: String,
    pub x: f64,
    pub y: f64,
    pub z: f64,
}

/// Cache represents an asset or document hash for change detection
#[derive(Debug, Clone)]
pub struct Cache {
    pub id: Option<i64>,
    pub path: String,
    pub hash: String,
    pub x: f64,
    pub y: f64,
    pub z: f64,
}

/// Post represents a blog post
#[derive(Debug, Clone)]
pub struct Post {
    pub id: Option<i64>,
    pub title: String,
    pub slug: String,
    pub description: String,
    pub content: String,
    pub banner_path: String,
    pub created_at: i64, // Unix timestamp for SQLite compatibility
    pub tags: Vec<String>,
    pub posts: Vec<String>,
    pub projects: Vec<String>,
    pub x: f64,
    pub y: f64,
    pub z: f64,
}

/// Project represents a portfolio project
#[derive(Debug, Clone)]
pub struct Project {
    pub id: Option<i64>,
    pub title: String,
    pub slug: String,
    pub description: String,
    pub content: String,
    pub banner_path: String,
    pub created_at: i64, // Unix timestamp for SQLite compatibility
    pub tags: Vec<String>,
    pub posts: Vec<String>,
    pub projects: Vec<String>,
    pub x: f64,
    pub y: f64,
    pub z: f64,
}

/// Tag represents a category or skill
#[derive(Debug, Clone)]
pub struct Tag {
    pub id: Option<i64>,
    pub title: String,
    pub slug: String,
    pub description: String,
    pub content: String,
    pub banner_path: String,
    pub icon: String,
    pub created_at: i64, // Unix timestamp for SQLite compatibility
    pub tags: Vec<String>,
    pub posts: Vec<String>,
    pub projects: Vec<String>,
    pub x: f64,
    pub y: f64,
    pub z: f64,
}

/// Many-to-many relationship structs
#[derive(Debug, Clone)]
pub struct PostToTag {
    pub post_id: i64,
    pub tag_id: i64,
}

#[derive(Debug, Clone)]
pub struct PostToPost {
    pub source_post_id: i64,
    pub target_post_id: i64,
}

#[derive(Debug, Clone)]
pub struct PostToProject {
    pub post_id: i64,
    pub project_id: i64,
}

#[derive(Debug, Clone)]
pub struct ProjectToTag {
    pub project_id: i64,
    pub tag_id: i64,
}

#[derive(Debug, Clone)]
pub struct ProjectToProject {
    pub source_project_id: i64,
    pub target_project_id: i64,
}

#[derive(Debug, Clone)]
pub struct TagToTag {
    pub source_tag_id: i64,
    pub target_tag_id: i64,
}

impl From<Doc> for Post {
    fn from(doc: Doc) -> Self {
        Self {
            id: None,
            title: doc.title,
            slug: doc.slug,
            description: doc.description,
            content: doc.content,
            banner_path: doc.banner_path,
            created_at: doc.created_at.0.timestamp(),
            tags: doc.tags,
            posts: doc.posts,
            projects: doc.projects,
            x: doc.x,
            y: doc.y,
            z: doc.z,
        }
    }
}

impl From<Doc> for Project {
    fn from(doc: Doc) -> Self {
        Self {
            id: None,
            title: doc.title,
            slug: doc.slug,
            description: doc.description,
            content: doc.content,
            banner_path: doc.banner_path,
            created_at: doc.created_at.0.timestamp(),
            tags: doc.tags,
            posts: doc.posts,
            projects: doc.projects,
            x: doc.x,
            y: doc.y,
            z: doc.z,
        }
    }
}

impl From<Doc> for Tag {
    fn from(doc: Doc) -> Self {
        Self {
            id: None,
            title: doc.title,
            slug: doc.slug,
            description: doc.description,
            content: doc.content,
            banner_path: doc.banner_path,
            icon: doc.icon,
            created_at: doc.created_at.0.timestamp(),
            tags: doc.tags,
            posts: doc.posts,
            projects: doc.projects,
            x: doc.x,
            y: doc.y,
            z: doc.z,
        }
    }
}

/// Calculate an MD5 hash for content
pub fn hash(content: &[u8]) -> String {
    format!("{:x}", md5::compute(content))
}