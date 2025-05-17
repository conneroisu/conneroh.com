use std::env;
use std::future::Future;
use std::path::{Path, PathBuf};
use std::pin::Pin;
use tokio::fs;
use tokio::io;

/// Recursively visits a directory and counts the number of files.
///
/// * `path`: The path to the directory to visit.
fn visit(path: PathBuf) -> Pin<Box<dyn Future<Output = io::Result<i32>> + Send>> {
    Box::pin(async move {
        let mut entries = fs::read_dir(path).await?;
        let mut count = 0;
        let mutex = tokio::sync::Mutex::new(());

        while let Some(entry) = entries.next_entry().await? {
            let path = entry.path();
            if path.is_dir() {
                count += visit(path).await?;
            } else {
                let _guard = mutex.lock().await;
                count += 1;
                println!("{}", path.display());
            } // _guard is dropped here
        }
        Ok(count)
    })
}

#[tokio::main]
async fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("Usage: {} <path>", args[0]);
        std::process::exit(1);
    }
    let path = Path::new(&args[1]);
    let start = std::time::Instant::now();
    let num = visit(PathBuf::from(path)).await?;
    eprintln!("Finished in {:?}", start.elapsed());
    eprintln!("{} files", num);
    Ok(())
}
