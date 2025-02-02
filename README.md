# Karakuri
## Introduction
`karakuri` is a platform for managing and running containers.  
This platform consists of three components.
1. `karakuri`  
   A CLI tool that provides the ability to send requests to high-level container runtime.
1. `hitoha`  
   High-level container runtime.  
   Runs as a daemon process and performs container networking, lifecycle management and image management.  
   `hitoha` provide a REST API as an interface and executes low-level container runtime in response to received requests.
1. `futaba`  
   Low-level container runtime.  
   Actual container operation, including namespace isolation, mounts, root filesystem changes, etc.