use std::path::Path;

use super::{ASSETS_LOC, POSTS_LOC, PROJECTS_LOC, TAGS_LOC};

/// Slugify returns the slugified path of a document or media asset.
/// This removes the extension and directory prefix
pub fn slugify(path: &str) -> String {
    if let Some(stripped) = path.strip_prefix(ASSETS_LOC) {
        return stripped.to_string();
    }

    pathify(path)
        .trim_end_matches(['/', '\\'])
        .split('.')
        .next()
        .unwrap_or("")
        .to_string()
}

/// Pathify returns the path to the document page or media asset page by stripping prefixes
pub fn pathify(path: &str) -> String {
    if let Some(stripped) = path.strip_prefix(POSTS_LOC) {
        return stripped.to_string();
    }
    if let Some(stripped) = path.strip_prefix(PROJECTS_LOC) {
        return stripped.to_string();
    }
    if let Some(stripped) = path.strip_prefix(TAGS_LOC) {
        return stripped.to_string();
    }
    if let Some(stripped) = path.strip_prefix(ASSETS_LOC) {
        return stripped.to_string();
    }
    
    panic!("Failed to pathify {}", path);
}

/// Check if file extension indicates it's a valid media type
pub fn is_allowed_media_type(path: &str) -> bool {
    let allowed_extensions = [".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif", ".tiff", ".svg", ".pdf"];
    
    if let Some(extension) = Path::new(path).extension() {
        if let Some(ext_str) = extension.to_str() {
            return allowed_extensions.contains(&format!(".{}", ext_str).as_str());
        }
    }
    
    false
}

/// Check if file extension indicates it's a valid document type
pub fn is_allowed_document_type(path: &str) -> bool {
    path.ends_with(".md")
}

/// Check if the file is either a media type or document type
pub fn is_allowed_asset(path: &str) -> bool {
    is_allowed_media_type(path) || is_allowed_document_type(path)
}

/// Get the content type for a file based on its extension
pub fn get_content_type(path: &str) -> String {
    match mime_guess::from_path(path).first_raw() {
        Some(mime) => mime.to_string(),
        None => "application/octet-stream".to_string(),
    }
}