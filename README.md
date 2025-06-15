# 🚀 Gwencat-LAN

![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-purple.svg)
<a href="https://github.com/gigachad80/TxtRipper/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
 

<img src="https://github.com/user-attachments/assets/a0c1b3af-5485-4ead-b53f-73c7e43cf723" alt="Description" width="200" height="200">

> **A high-performance, cross-platform file transfer utility with auto-mode detection, progress monitoring, integrity verification, and robust network handling for seamless file transfers between Linux and Windows & Mac systems in local networks**

---

### 📌 Overview

**Gwencat-LAN** is a Go-based file transfer utility designed for simple, fast, and reliable file transfers over TCP networks within local area networks. It automatically detects whether you want to send or receive files based on file existence, supports both explicit and implicit modes, includes SHA256 integrity verification, real-time progress monitoring, and handles network timeouts gracefully. Perfect for transferring files between different operating systems in local networks.

### 🌟 Key Features

✅ **Auto-mode detection** - No need to remember sender/receiver flags  
✅ **Cross-platform** - Works on Linux, Windows, macOS  
✅ **Progress monitoring** - Real-time transfer statistics  
✅ **Integrity verification** - SHA256 checksum validation  
✅ **Network resilience** - Configurable timeouts and error handling  
✅ **IP filtering** - Security through connection restrictions  
✅ **Human-readable output** - Formatted file sizes and transfer speeds  
✅ **Zero dependencies** - Single binary, no external requirements  

### 📚 Requirements & Dependencies

- **Go 1.16+** (for compilation)
- **Network connectivity** between source and destination machines
- **Open firewall port** on the receiving machine (default: any port you choose)
- **Target machines**: Linux, Windows, macOS (any platform Go supports)

### ⚡ Quick Installation & Usage

### Installation

1. **Clone and build:**
   ```bash
   git clone https://github.com/gigachad80/Gwencat-LAN
 
   go build -o gwencat-lan ( For Linux / macOS )
   go build -o gwencat-lan.exe ( For Windows)
   ```
   
  Or you can build by yourself with your own custom name based kn your architecture 

2. **Or download pre-built binaries** from releases page

### Basic Usage

Transfer files in two simple steps:

```bash
# Step 1: Start receiver first
./gwencat-lan -p 4444 received_file.txt

# Step 2: Send file from another machine
./gwencat-lan -r 192.168.1.100 -p 4444 myfile.txt
```

### 🚀 Demo Syntax & Examples

### Quick Reference

**On Listener/Receiver's side:**
```bash
./gwencat-lan -r {sender_ip} -p {port} {received_filename}
```

**On Sender's side:**
```bash
./gwencat-lan -r {listener_ip} -p {port} {filename_to_send}
```

> [!IMPORTANT]
> First you have to start the listener/receiver, then the sender's IP

### Basic File Transfer

```bash
# Send file to remote machine (auto-detected sender mode)
gwencat-lan -r 192.168.1.100 -p 4444 document.pdf

# Receive file (auto-detected receiver mode)
gwencat-lan -p 4444 received_document.pdf
```

### Cross-Platform Examples

#### Windows to Linux
```cmd
# Windows (Sender): 192.168.181.1
gwencat-lan.exe -r 192.168.181.128 -p 4444 "C:\Users\John\report.docx"

# Linux (Receiver): 192.168.181.128  
./gwencat-lan -p 4444 /home/user/report.docx
```

#### Linux to Windows
```bash
# Linux (Sender): 192.168.181.128
./gwencat-lan -r 192.168.181.1 -p 4444 /home/user/backup.tar.gz

# Windows (Receiver): 192.168.181.1
gwencat-lan.exe -p 4444 "C:\Downloads\backup.tar.gz"
```

### Advanced Usage with All Features

```bash
# Send with progress, verification, and custom timeout
gwencat-lan -r 192.168.1.100 -p 4444 -progress -v -t 60 largefile.zip

# Receive with IP filtering and verification
gwencat-lan -r 192.168.1.50 -p 4444 -v -progress received.zip

# Explicit mode specification
gwencat-lan -mode send -r 192.168.1.100 -p 4444 -progress myfile.txt
gwencat-lan -mode receive -p 4444 -a -v newfile.txt
```

### Security & Filtering

```bash
# Only accept connections from specific IP
gwencat-lan -r 192.168.1.50 -p 4444 secure_file.txt

# Listen on all interfaces (less secure)
gwencat-lan -a -p 4444 public_file.txt
```

### 📖 Command Line Options

| Option | Description | Example |
|--------|-------------|---------|
| `-mode` | Explicit mode: `send` or `receive` | `-mode send` |
| `-r` | Remote IP address | `-r 192.168.1.100` |
| `-p` | Port number (required) | `-p 4444` |
| `-t` | Connection timeout in seconds | `-t 60` |
| `-v` | Enable SHA256 integrity verification | `-v` |
| `-progress` | Show real-time transfer progress | `-progress` |
| `-a` | Listen on all interfaces (receiver only) | `-a` |

### 🔧 Usage Scenarios

### Scenario 1: Quick File Share
```bash
# Person A wants to send a file to Person B
# Person B runs: gwencat-lan -p 8080 incoming.zip
# Person A runs: gwencat-lan -r <PersonB_IP> -p 8080 myfile.zip
```

### Scenario 2: Backup Transfer
```bash
# Server backup to local machine with verification
# Local: gwencat-lan -p 9999 -v backup_$(date +%Y%m%d).tar.gz
# Server: gwencat-lan -r <local_ip> -p 9999 -v -progress /backup/full.tar.gz
```

### Scenario 3: Development File Sync
```bash
# Sync build artifacts between dev machines
# Target: gwencat-lan -r 192.168.1.10 -p 3000 -progress build.zip
# Source: gwencat-lan -p 3000 -v new_build.zip
```

### 🛡️ Security Considerations

- **IP Filtering**: Use `-r` parameter to restrict connections to specific IPs
- **Local Networks**: Designed for trusted network environments
- **Firewall**: Ensure receiving machine allows connections on chosen port
- **Verification**: Always use `-v` flag for important files to verify integrity

### 🐛 Troubleshooting

### Common Issues

#### "bind: cannot assign requested address"
```bash
# Wrong: trying to bind to remote IP
gwencat-lan -l 192.168.1.100 -p 4444 file.txt  # ❌

# Correct: specify remote IP as target
gwencat-lan -r 192.168.1.100 -p 4444 file.txt  # ✅
```

#### "connection refused"
```bash
# Make sure receiver is started first and port is open
# Check firewall: sudo ufw allow 4444  (Linux)
# Check Windows Firewall settings
```

#### "timeout"
```bash
# Increase timeout for large files or slow networks
gwencat-lan -r 192.168.1.100 -p 4444 -t 120 largefile.iso
```

### 📊 Performance Tips

- **Buffer Size**: Optimized 32KB buffer for best performance
- **Progress Updates**: Limited to 100ms intervals to avoid overhead  
- **Network Timeout**: Auto-scales to 10x connection timeout for transfers
- **Large Files**: Use `-progress` flag to monitor long transfers

### 🔄 Auto-Detection Logic

Gwencat-LAN automatically determines mode based on file existence:
- **File exists** → Send mode (connects to remote IP)
- **File doesn't exist** → Receive mode (listens for connections)
- **Override**: Use `-mode` flag to force specific behavior

### 🤔 Why This Name?

Initially, I thought of "goncat" (combining "Go" with "netcat"), but that name was already taken. During development, I somehow remembered the "Lucky Girl" episode featuring Gwendolyn from Ben 10, and since this program is designed specifically for LAN (Local Area Network) usage, I named it **Gwencat-LAN**. The name reflects both the Go based implementation of netcat but only for file transfer and the tool's LAN-focused functionality. I also have a WAN version of this utility, so the LAN suffix helps distinguish between the two variants. Also logo of this is similar to lucky girl mask 😉



### 🙃 Why I Created This

I developed Gwencat-LAN to solve the common frustration of transferring files between different systems (especially Windows to Linux running inside VM) without relying on external services, complex setup procedures, or slow protocols. The development was driven by a few specific reasons:

1.  Personal Project Needs 😖: 
I had my own project which required a file transfer utility. Initially, I decided to develop a full netcat implementation, but I didn't succeed with that ambitious goal, so I stuck with creating a simple yet effective file transfer utility that met my immediate needs.

2.  Cross-Platform Requirements🥱 :
Even though ncat is cross-platform, I wanted to build my own version to have complete control over the features, behavior, and customization possibilities for my specific use cases.

Traditional tools like `scp` require SSH setup, standard `netcat` lacks progress indication and error handling, and most GUI solutions are platform-specific. Gwencat-LAN provides a single, lightweight binary that works identically across platforms with intelligent auto-detection, making file transfers as simple as running one command on each machine.

### ⌚ Total time Spent in development , testing , README . 
Current implementation (development, testing, README): Approximately 1 hour
Time spent troubleshooting Netcat implementation: Around 2 hours 30 minutes, which unfortunately did not yield a stable result
Previous two variations: Successfully developed and tested in 2 hours, but both occasionally encountered errors

### 📝 Roadmap / To-do

-   [x] Release Cross-Platform Executables 

### 📞 Contact

**📧 Email:** pookielinuxuser@tutamail.com

### 📄 License

Licensed under **GNU Affero General Public License 3.0**

---

**🕒 Last Updated:** June 14, 2025
**🕒 First Published:** June 14, 2025
<div align="center">

**Happy file transferring! 🚀**

*Made with ❤️ for seamless cross-platform file transfers*

</div>
