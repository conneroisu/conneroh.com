use anyhow::{Context, Result};
use serde_yaml::Value;
use tracing::debug;

use crate::models::Doc;

/// Process markdown content with frontmatter
pub fn process_markdown(content: &[u8], doc: &mut Doc) -> Result<()> {
    let content_str = std::str::from_utf8(content)
        .context("Failed to convert markdown content to UTF-8")?;
    
    // Extract frontmatter and content
    let (frontmatter, markdown_content) = extract_frontmatter(content_str)
        .context("Failed to extract frontmatter")?;
    
    // Parse frontmatter
    let parsed: Value = serde_yaml::from_str(&frontmatter)
        .context("Failed to parse frontmatter as YAML")?;
    
    // Convert Value to Doc
    populate_doc_from_yaml(doc, &parsed)
        .context("Failed to populate Doc from YAML")?;
    
    // Set content
    doc.content = markdown_content.to_string();
    
    debug!("Processed markdown document: {}", doc.title);
    Ok(())
}

/// Extract frontmatter and content from markdown
fn extract_frontmatter(content: &str) -> Result<(String, &str)> {
    // Look for triple-dash frontmatter
    let mut lines = content.lines();
    
    // Check for starting ---
    if let Some(first_line) = lines.next() {
        if first_line.trim() != "---" {
            return Err(anyhow::anyhow!("Frontmatter must start with ---"));
        }
    } else {
        return Err(anyhow::anyhow!("Empty content"));
    }
    
    // Collect frontmatter lines until ending ---
    let mut frontmatter = Vec::new();
    let mut found_end = false;
    
    for line in &mut lines {
        if line.trim() == "---" {
            found_end = true;
            break;
        }
        frontmatter.push(line);
    }
    
    if !found_end {
        return Err(anyhow::anyhow!("Frontmatter must end with ---"));
    }
    
    // Join frontmatter lines
    let frontmatter = frontmatter.join("\n");
    
    // The rest is markdown content
    let content_start = frontmatter.len() + 8; // 8 = 2 * "---\n" (assuming \n line endings)
    let content = &content[content_start.min(content.len())..];
    
    Ok((frontmatter, content))
}

/// Populate Doc fields from YAML Value
fn populate_doc_from_yaml(doc: &mut Doc, yaml: &Value) -> Result<()> {
    // Extract values from YAML
    if let Some(title) = yaml.get("title").and_then(|v| v.as_str()) {
        doc.title = title.to_string();
    }
    
    if let Some(slug) = yaml.get("slug").and_then(|v| v.as_str()) {
        doc.slug = slug.to_string();
    }
    
    if let Some(description) = yaml.get("description").and_then(|v| v.as_str()) {
        doc.description = description.to_string();
    }
    
    if let Some(banner_path) = yaml.get("banner_path").and_then(|v| v.as_str()) {
        doc.banner_path = banner_path.to_string();
    }
    
    if let Some(icon) = yaml.get("icon").and_then(|v| v.as_str()) {
        doc.icon = icon.to_string();
    }
    
    // Parse tags
    if let Some(tags) = yaml.get("tags").and_then(|v| v.as_sequence()) {
        doc.tags = tags.iter()
            .filter_map(|v| v.as_str().map(|s| s.to_string()))
            .collect();
    }
    
    // Parse posts
    if let Some(posts) = yaml.get("posts").and_then(|v| v.as_sequence()) {
        doc.posts = posts.iter()
            .filter_map(|v| v.as_str().map(|s| s.to_string()))
            .collect();
    }
    
    // Parse projects
    if let Some(projects) = yaml.get("projects").and_then(|v| v.as_sequence()) {
        doc.projects = projects.iter()
            .filter_map(|v| v.as_str().map(|s| s.to_string()))
            .collect();
    }
    
    // Additional numeric values
    if let Some(x) = yaml.get("x").and_then(|v| v.as_f64()) {
        doc.x = x;
    }
    
    if let Some(y) = yaml.get("y").and_then(|v| v.as_f64()) {
        doc.y = y;
    }
    
    if let Some(z) = yaml.get("z").and_then(|v| v.as_f64()) {
        doc.z = z;
    }
    
    Ok(())
}