# Network Namespace Linker (nslinker)
This project automatically binds containerd containers to Linux network namespaces by creating symbolic links in /var/run/netns/. This allows easy network namespace management using ip netns exec commands.

### Features
- Lists running containerd containers.
- Extracts network namespace from each container.
- Creates symlinks in /var/run/netns/ for each container.
- Enables using ip netns exec <pod-name> for network operations.

### Prerequisites
- Linux system (Ubuntu/Debian recommended)
- containerd installed and running
- Go 1.18+ installed

### Installation & Setup
#### 1. Clone the Repository
    git clone https://github.com/yourusername/nslinker.git
    cd nslinker
    
#### 1. Build and Run the Project
    go run main.go

### How It Works
- Fetch running containerd containers using containerd.New().
- Extract PID and Pod Name from container config files.
- Create a symlink in /var/run/netns/<pod-name> pointing to the container’s network namespace.
- Allow network inspection via ip netns exec <pod-name> ....

Project idea gathered from: https://github.com/aerosouund/nslinker
