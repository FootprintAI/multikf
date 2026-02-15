#!/usr/bin/env python3
"""
Image Mirror Tool - Creates commands to mirror container images to private registries.

Usage:
    python mirror_image.py --image docker.io/kubeflow/training-operator:v1.5.0 --mirror reg.footprint-ai.com/kubeflow-mirror
    
    Optional arguments:
    --username USERNAME - Registry username for authentication
    --password PASSWORD - Registry password for authentication
    --batch-file FILE - File containing list of images to mirror (one per line)
    --execute - Execute the generated commands instead of just printing them
"""

import argparse
import subprocess
import sys
import os
from typing import List, Optional, Tuple

def parse_image_name(image_path: str) -> Tuple[str, str, str, str]:
    """
    Parse the image path into its components.
    
    Args:
        image_path: Full image path (e.g., docker.io/kubeflow/training-operator:v1.5.0)
        
    Returns:
        Tuple of (registry, repository, image, tag)
    """
    # Handle registry with port if present
    if '@' in image_path:
        # Handle digest format
        path_parts, digest = image_path.split('@')
        tag = f"@{digest}"
    elif ':' in image_path.split('/')[-1]:
        # Handle normal tag format
        path_parts, tag = image_path.rsplit(':', 1)
        tag = f":{tag}"
    else:
        # Default to latest if no tag is specified
        path_parts = image_path
        tag = ":latest"
    
    # Split the remaining path
    parts = path_parts.split('/')
    
    # Handle different path formats
    if len(parts) == 1:
        # Image only, assume docker.io/library
        registry = "docker.io"
        repo = "library"
        image = parts[0]
    elif len(parts) == 2:
        # Could be either registry/image or namespace/image
        if '.' in parts[0] or parts[0] == "localhost":
            # It's registry/image
            registry = parts[0]
            repo = ""
            image = parts[1]
        else:
            # It's namespace/image on default registry
            registry = "docker.io"
            repo = parts[0]
            image = parts[1]
    else:
        # registry/namespace/image or registry/namespace/subnamespace/image
        registry = parts[0]
        image = parts[-1]
        repo = '/'.join(parts[1:-1])
    
    return registry, repo, image, tag

def generate_mirror_commands(image_path: str, mirror_registry: str, username: Optional[str] = None, 
                             password: Optional[str] = None) -> List[str]:
    """
    Generate commands to pull, tag and push an image to a mirror registry.
    
    Args:
        image_path: Full image path
        mirror_registry: Target mirror registry
        username: Optional registry username
        password: Optional registry password
        
    Returns:
        List of shell commands to execute
    """
    commands = []
    
    # Parse the image name
    source_registry, repo, image, tag = parse_image_name(image_path)
    
    # Build the full source path
    if repo:
        source_path = f"{source_registry}/{repo}/{image}{tag}"
    else:
        source_path = f"{source_registry}/{image}{tag}"
    
    # Build the destination path
    if repo:
        dest_path = f"{mirror_registry}/{repo}/{image}{tag}"
    else:
        dest_path = f"{mirror_registry}/{image}{tag}"
    
    # Add login command if credentials are provided
    if username and password:
        commands.append(f"echo {password} | docker login {mirror_registry} --username {username} --password-stdin")
    
    # Add pull, tag and push commands
    commands.append(f"docker pull {source_path}")
    commands.append(f"docker tag {source_path} {dest_path}")
    commands.append(f"docker push {dest_path}")
    
    return commands

def process_batch_file(batch_file: str, mirror_registry: str, username: Optional[str] = None, 
                       password: Optional[str] = None, execute: bool = False) -> None:
    """
    Process a batch file containing one image per line.
    
    Args:
        batch_file: Path to file with image list
        mirror_registry: Target mirror registry
        username: Optional registry username
        password: Optional registry password
        execute: Whether to execute commands
    """
    if not os.path.exists(batch_file):
        print(f"Error: Batch file {batch_file} not found.")
        sys.exit(1)
        
    with open(batch_file, 'r') as f:
        images = [line.strip() for line in f if line.strip() and not line.strip().startswith('#')]
    
    print(f"Processing {len(images)} images from {batch_file}...")
    
    # Add login command once if credentials are provided
    if username and password and execute:
        login_cmd = f"echo {password} | docker login {mirror_registry} --username {username} --password-stdin"
        print(f"Executing: {login_cmd}")
        subprocess.run(login_cmd, shell=True, check=True)
    
    for image in images:
        print(f"\nMirroring: {image}")
        commands = generate_mirror_commands(image, mirror_registry, None, None)  # Skip login command for batch
        
        if execute:
            for cmd in commands:
                print(f"Executing: {cmd}")
                try:
                    subprocess.run(cmd, shell=True, check=True)
                except subprocess.CalledProcessError as e:
                    print(f"Error executing command: {e}")
        else:
            for cmd in commands:
                print(cmd)

def main() -> None:
    parser = argparse.ArgumentParser(description="Generate commands to mirror container images")
    parser.add_argument("--image", help="Source image path (e.g., docker.io/kubeflow/training-operator:v1.5.0)")
    parser.add_argument("--mirror", required=True, help="Mirror registry (e.g., reg.footprint-ai.com/kubeflow-mirror)")
    parser.add_argument("--username", help="Registry username for authentication")
    parser.add_argument("--password", help="Registry password for authentication")
    parser.add_argument("--batch-file", help="File containing list of images to mirror (one per line)")
    parser.add_argument("--execute", action="store_true", help="Execute the commands instead of printing them")
    
    args = parser.parse_args()
    
    if args.batch_file:
        process_batch_file(args.batch_file, args.mirror, args.username, args.password, args.execute)
    elif args.image:
        commands = generate_mirror_commands(args.image, args.mirror, args.username, args.password)
        
        if args.execute:
            for cmd in commands:
                print(f"Executing: {cmd}")
                try:
                    subprocess.run(cmd, shell=True, check=True)
                except subprocess.CalledProcessError as e:
                    print(f"Error executing command: {e}")
        else:
            print("\n".join(commands))
    else:
        parser.print_help()
        print("\nError: Either --image or --batch-file must be specified.")
        sys.exit(1)

if __name__ == "__main__":
    main()
