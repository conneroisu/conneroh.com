#!/usr/bin/env python3
"""
SVG Internal Links Removal Script

This script removes internal ID references from SVG files to prevent conflicts
when SVGs are used in contexts where IDs might clash with other elements.

Patterns removed:
- xlink:href="#[id]" attributes 
- href="#[id]" attributes
- url(#[id]) references in fill, mask, filter attributes
- Corresponding definition elements (<linearGradient>, <mask>, <filter>, etc.)

Usage:
    python fix_svg_internal_links.py [options]
    
Options:
    --dry-run    Show what would be changed without modifying files
    --backup     Create .bak files before making changes
    --verbose    Show detailed processing information
"""

import xml.etree.ElementTree as ET
import re
import argparse
import os
import shutil
from pathlib import Path
from typing import Set, Dict, List, Tuple
import sys

class SVGProcessor:
    """Processes SVG files to remove internal link references."""
    
    # Namespaces used in SVG files
    NAMESPACES = {
        'svg': 'http://www.w3.org/2000/svg',
        'xlink': 'http://www.w3.org/1999/xlink'
    }
    
    # Elements that typically contain ID references that should be removed
    DEFINITION_ELEMENTS = {
        'linearGradient', 'radialGradient', 'pattern', 'mask', 'clipPath', 
        'filter', 'marker', 'symbol', 'use'
    }
    
    # Attributes that may contain url(#id) references
    URL_ATTRIBUTES = {
        'fill', 'stroke', 'mask', 'clip-path', 'filter', 'marker-start', 
        'marker-mid', 'marker-end', 'opacity-mask'
    }
    
    # Fallback colors for removed gradients
    FALLBACK_COLORS = {
        'gradient': '#666666',  # Neutral gray
        'pattern': '#cccccc',   # Light gray
        'default': '#000000'    # Black
    }
    
    def __init__(self, verbose: bool = False):
        self.verbose = verbose
        self.stats = {
            'files_processed': 0,
            'files_modified': 0,
            'references_removed': 0,
            'elements_removed': 0
        }
    
    def log(self, message: str, level: str = 'INFO'):
        """Log message if verbose mode is enabled."""
        if self.verbose:
            print(f"[{level}] {message}")
    
    def find_target_files(self, directory: str) -> List[Path]:
        """Find all SVG files that need processing based on the known problematic files."""
        target_filenames = {
            'apache-original.svg', 'clion-original.svg', 'd3js-original.svg', 
            'eclipse-original.svg', 'dropwizard-original.svg', 'gentoo-original.svg',
            'gimp-original.svg', 'goland-original.svg', 'json-original.svg',
            'jira-original.svg', 'moodle-original.svg', 'nextjs-original.svg',
            'ocaml-original.svg', 'poetry-original.svg', 'prolog-original.svg',
            'rollup-original.svg', 'ruby-original.svg', 'webstorm-original.svg',
            'xcode-original.svg', 'maven-original.svg', 'renpy-original.svg',
            'vscode-original.svg'
        }
        
        directory_path = Path(directory)
        found_files = []
        
        for filename in target_filenames:
            file_path = directory_path / filename
            if file_path.exists():
                found_files.append(file_path)
                self.log(f"Found target file: {filename}")
            else:
                self.log(f"Target file not found: {filename}", 'WARNING')
        
        return found_files
    
    def extract_internal_ids(self, root: ET.Element) -> Set[str]:
        """Extract all internal IDs that are referenced by href or url() patterns."""
        referenced_ids = set()
        
        # Register namespaces for XPath queries
        for prefix, uri in self.NAMESPACES.items():
            ET.register_namespace(prefix, uri)
        
        # Find xlink:href and href references
        for elem in root.iter():
            # Check xlink:href attributes
            href = elem.get(f'{{{self.NAMESPACES["xlink"]}}}href')
            if href and href.startswith('#'):
                referenced_ids.add(href[1:])
                self.log(f"Found xlink:href reference: {href}")
            
            # Check href attributes
            href = elem.get('href')
            if href and href.startswith('#'):
                referenced_ids.add(href[1:])
                self.log(f"Found href reference: {href}")
            
            # Check url(#id) patterns in various attributes
            for attr_name in self.URL_ATTRIBUTES:
                attr_value = elem.get(attr_name)
                if attr_value:
                    # Match url(#id) patterns
                    url_matches = re.findall(r'url\(#([^)]+)\)', attr_value)
                    for match in url_matches:
                        referenced_ids.add(match)
                        self.log(f"Found url() reference in {attr_name}: #{match}")
        
        return referenced_ids
    
    def remove_internal_references(self, root: ET.Element, referenced_ids: Set[str]) -> int:
        """Remove internal references and return count of removed references."""
        removed_count = 0
        
        for elem in root.iter():
            # Remove xlink:href attributes
            href_attr = f'{{{self.NAMESPACES["xlink"]}}}href'
            if elem.get(href_attr) and elem.get(href_attr).startswith('#'):
                elem.attrib.pop(href_attr)
                removed_count += 1
                self.log(f"Removed xlink:href from {elem.tag}")
            
            # Remove href attributes
            if elem.get('href') and elem.get('href').startswith('#'):
                elem.attrib.pop('href')
                removed_count += 1
                self.log(f"Removed href from {elem.tag}")
            
            # Remove url(#id) patterns from attributes
            for attr_name in self.URL_ATTRIBUTES:
                attr_value = elem.get(attr_name)
                if attr_value:
                    # Replace url(#id) with fallback colors
                    original_value = attr_value
                    
                    # For fill attributes, use appropriate fallback colors
                    if attr_name == 'fill':
                        attr_value = re.sub(r'url\(#[^)]+\)', self.FALLBACK_COLORS['gradient'], attr_value)
                    else:
                        # For other attributes, remove the url() reference entirely
                        attr_value = re.sub(r'url\(#[^)]+\)', '', attr_value)
                    
                    if attr_value != original_value:
                        if attr_value.strip():
                            elem.set(attr_name, attr_value)
                        else:
                            # Remove empty attributes
                            elem.attrib.pop(attr_name, None)
                        removed_count += 1
                        self.log(f"Modified {attr_name} attribute in {elem.tag}")
        
        return removed_count
    
    def remove_definition_elements(self, root: ET.Element, referenced_ids: Set[str]) -> int:
        """Remove definition elements that were referenced by internal links."""
        removed_count = 0
        elements_to_remove = []
        
        # Find all elements with IDs that were referenced
        for elem in root.iter():
            elem_id = elem.get('id')
            if elem_id and elem_id in referenced_ids:
                # Check if this is a definition element we should remove
                tag_name = elem.tag.split('}')[-1] if '}' in elem.tag else elem.tag
                
                if tag_name in self.DEFINITION_ELEMENTS:
                    elements_to_remove.append((elem, elem.getparent()))
                    self.log(f"Marking for removal: {tag_name}#{elem_id}")
        
        # Remove the elements (do this after iteration to avoid modifying during iteration)
        for elem, parent in elements_to_remove:
            if parent is not None:
                parent.remove(elem)
                removed_count += 1
                self.log(f"Removed element: {elem.tag}#{elem.get('id')}")
        
        return removed_count
    
    def process_svg_file(self, file_path: Path, dry_run: bool = False, backup: bool = False) -> Tuple[bool, str]:
        """Process a single SVG file. Returns (success, message)."""
        try:
            self.log(f"Processing: {file_path}")
            
            # Parse the SVG file
            tree = ET.parse(file_path)
            root = tree.getroot()
            
            # Extract internal IDs that are referenced
            referenced_ids = self.extract_internal_ids(root)
            
            if not referenced_ids:
                self.log(f"No internal references found in {file_path.name}")
                return True, "No internal references found"
            
            self.log(f"Found {len(referenced_ids)} internal references: {', '.join(referenced_ids)}")
            
            if dry_run:
                return True, f"Would remove {len(referenced_ids)} internal references"
            
            # Create backup if requested
            if backup:
                backup_path = file_path.with_suffix(file_path.suffix + '.bak')
                shutil.copy2(file_path, backup_path)
                self.log(f"Created backup: {backup_path}")
            
            # Remove internal references
            refs_removed = self.remove_internal_references(root, referenced_ids)
            
            # Remove corresponding definition elements
            elems_removed = self.remove_definition_elements(root, referenced_ids)
            
            # Write the modified SVG back to file
            # Preserve the XML declaration and encoding
            tree.write(file_path, encoding='utf-8', xml_declaration=True)
            
            self.stats['files_modified'] += 1
            self.stats['references_removed'] += refs_removed
            self.stats['elements_removed'] += elems_removed
            
            message = f"Removed {refs_removed} references and {elems_removed} elements"
            self.log(f"Successfully processed {file_path.name}: {message}")
            
            return True, message
            
        except ET.ParseError as e:
            error_msg = f"XML parsing error: {e}"
            self.log(error_msg, 'ERROR')
            return False, error_msg
        except Exception as e:
            error_msg = f"Unexpected error: {e}"
            self.log(error_msg, 'ERROR')
            return False, error_msg
    
    def process_directory(self, directory: str, dry_run: bool = False, backup: bool = False):
        """Process all target SVG files in the directory."""
        target_files = self.find_target_files(directory)
        
        if not target_files:
            print("No target SVG files found in the directory.")
            return
        
        print(f"Found {len(target_files)} target SVG files to process")
        if dry_run:
            print("DRY RUN MODE - No files will be modified")
        
        print("-" * 60)
        
        success_count = 0
        for file_path in target_files:
            self.stats['files_processed'] += 1
            success, message = self.process_svg_file(file_path, dry_run, backup)
            
            status = "✓" if success else "✗"
            print(f"{status} {file_path.name}: {message}")
            
            if success:
                success_count += 1
        
        print("-" * 60)
        print(f"Processing complete!")
        print(f"Files processed: {self.stats['files_processed']}")
        print(f"Files successfully processed: {success_count}")
        
        if not dry_run:
            print(f"Files modified: {self.stats['files_modified']}")
            print(f"Total references removed: {self.stats['references_removed']}")
            print(f"Total elements removed: {self.stats['elements_removed']}")

def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Remove internal links from SVG files to prevent ID conflicts",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python fix_svg_internal_links.py                    # Process current directory
  python fix_svg_internal_links.py --dry-run          # Show what would be changed 
  python fix_svg_internal_links.py --backup --verbose # Create backups and show details
        """
    )
    
    parser.add_argument(
        '--directory', '-d',
        default='.',
        help='Directory containing SVG files (default: current directory)'
    )
    
    parser.add_argument(
        '--dry-run',
        action='store_true',
        help='Show what would be changed without modifying files'
    )
    
    parser.add_argument(
        '--backup',
        action='store_true',
        help='Create .bak files before making changes'
    )
    
    parser.add_argument(
        '--verbose', '-v',
        action='store_true',
        help='Show detailed processing information'
    )
    
    args = parser.parse_args()
    
    # Verify directory exists
    if not os.path.isdir(args.directory):
        print(f"Error: Directory '{args.directory}' does not exist")
        sys.exit(1)
    
    # Create processor and run
    processor = SVGProcessor(verbose=args.verbose)
    processor.process_directory(
        directory=args.directory,
        dry_run=args.dry_run,
        backup=args.backup
    )

if __name__ == '__main__':
    main()