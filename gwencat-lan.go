package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	mode         = flag.String("mode", "", "Mode: 'send' or 'receive'")
	remoteIP     = flag.String("r", "", "Remote IP address (target IP for sender, source IP filter for receiver)")
	port         = flag.Int("p", 0, "Port to connect to or listen on")
	timeout      = flag.Int("t", 30, "Connection timeout in seconds")
	verify       = flag.Bool("v", false, "Verify file integrity with SHA256 checksum")
	showProgress = flag.Bool("progress", false, "Show transfer progress")
	listenAll    = flag.Bool("a", false, "Listen on all interfaces (0.0.0.0) - receiver mode only")
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		printUsage()
		os.Exit(1)
	}

	filePath := flag.Arg(0)

	// Auto-detect mode if not specified
	if *mode == "" {
		if fileExists(filePath) {
			*mode = "send"
		} else {
			*mode = "receive"
		}
	}

	// Validate parameters based on mode
	if *mode == "send" {
		if *remoteIP == "" || *port == 0 {
			fmt.Println("Error: Send mode requires -r <remote_ip> -p <port>")
			printUsage()
			os.Exit(1)
		}
		if !fileExists(filePath) {
			log.Fatalf("File does not exist: %s", filePath)
		}
		sendFile(*remoteIP, *port, filePath)
	} else if *mode == "receive" {
		if *port == 0 {
			fmt.Println("Error: Receive mode requires -p <port>")
			printUsage()
			os.Exit(1)
		}
		receiveFile(*remoteIP, *port, filePath)
	} else {
		fmt.Println("Error: Invalid mode. Use 'send' or 'receive'")
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Enhanced File Transfer Utility")
	fmt.Println("\nUsage:")
	fmt.Println("  Auto-detect mode:")
	fmt.Println("    goncat -r <remote_ip> -p <port> [options] <file_path>")
	fmt.Println("    (if file exists = send mode, if not = receive mode)")
	fmt.Println("")
	fmt.Println("  Explicit mode:")
	fmt.Println("    Send:    goncat -mode send -r <remote_ip> -p <port> [options] <file_to_send>")
	fmt.Println("    Receive: goncat -mode receive -p <port> [options] <file_to_save>")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -mode <send|receive>  Explicitly set mode (optional, auto-detected)")
	fmt.Println("  -r <ip>              Remote IP (required for send, optional filter for receive)")
	fmt.Println("  -p <port>            Port number (required)")
	fmt.Println("  -t <seconds>         Connection timeout (default: 30)")
	fmt.Println("  -v                   Verify file integrity with SHA256 checksum")
	fmt.Println("  -progress            Show transfer progress")
	fmt.Println("  -a                   Listen on all interfaces (receive mode only)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Send file to 192.168.1.100 on port 4444")
	fmt.Println("  goncat -r 192.168.1.100 -p 4444 myfile.txt")
	fmt.Println("")
	fmt.Println("  # Receive file on port 4444, save as newfile.txt")
	fmt.Println("  goncat -p 4444 newfile.txt")
	fmt.Println("")
	fmt.Println("  # Receive with IP filtering and verification")
	fmt.Println("  goncat -mode receive -r 192.168.1.50 -p 4444 -v received.txt")
}

func sendFile(remoteIP string, port int, filePath string) {
	address := net.JoinHostPort(remoteIP, strconv.Itoa(port))

	log.Printf("Connecting to %s (timeout: %ds)...", address, *timeout)

	conn, err := net.DialTimeout("tcp", address, time.Duration(*timeout)*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", address, err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Duration(*timeout*10) * time.Second))

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}

	fileName := filepath.Base(filePath)
	fileSize := fileInfo.Size()

	log.Printf("Sending file: %s (%s)", fileName, formatBytes(fileSize))

	var originalChecksum string
	if *verify {
		log.Printf("Calculating checksum...")
		originalChecksum, err = calculateChecksum(filePath)
		if err != nil {
			log.Fatalf("Failed to calculate checksum: %v", err)
		}
		log.Printf("File checksum: %s", originalChecksum)
	}

	var bytesSent int64
	start := time.Now()

	if *showProgress {
		bytesSent, err = copyWithProgress(conn, file, fileSize, "Sending")
	} else {
		bytesSent, err = io.Copy(conn, file)
	}

	if err != nil {
		log.Fatalf("Failed to send file: %v", err)
	}

	duration := time.Since(start)
	speed := calculateSpeed(bytesSent, duration)

	log.Printf("✓ File sent successfully: %s (%s in %v, %s/s)",
		fileName, formatBytes(bytesSent), duration.Round(time.Millisecond), formatBytes(speed))

	if *verify {
		log.Printf("✓ Checksum: %s", originalChecksum)
	}
}

func receiveFile(sourceIPFilter string, port int, filePath string) {
	var listenAddr string

	if *listenAll {
		listenAddr = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
		log.Printf("Listening on all interfaces at port %d", port)
	} else {
		// Listen on all local interfaces, but we'll filter connections if sourceIPFilter is provided
		listenAddr = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
		if sourceIPFilter != "" {
			log.Printf("Listening on port %d, will accept connections only from %s", port, sourceIPFilter)
		} else {
			log.Printf("Listening on port %d, will accept connections from any IP", port)
		}
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}
	defer listener.Close()

	log.Printf("Waiting for connection (timeout: %ds)...", *timeout)

	if tcpListener, ok := listener.(*net.TCPListener); ok {
		tcpListener.SetDeadline(time.Now().Add(time.Duration(*timeout) * time.Second))
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Failed to accept connection: %v", err)
	}
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

	// Apply IP filtering if specified
	if sourceIPFilter != "" && remoteAddr.IP.String() != sourceIPFilter {
		log.Printf("✗ Rejected connection from %s (expected %s)", remoteAddr.IP.String(), sourceIPFilter)
		return
	}

	log.Printf("✓ Connection established from %s", remoteAddr.IP.String())

	conn.SetDeadline(time.Now().Add(time.Duration(*timeout*10) * time.Second))

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create directory: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	log.Printf("Receiving file...")

	var bytesReceived int64
	start := time.Now()

	if *showProgress {
		bytesReceived, err = copyWithProgressUnknownSize(file, conn, "Receiving")
	} else {
		bytesReceived, err = io.Copy(file, conn)
	}

	if err != nil {
		log.Fatalf("Failed to receive file: %v", err)
	}

	duration := time.Since(start)
	speed := calculateSpeed(bytesReceived, duration)
	fileName := filepath.Base(filePath)

	log.Printf("✓ File received successfully: %s (%s in %v, %s/s)",
		fileName, formatBytes(bytesReceived), duration.Round(time.Millisecond), formatBytes(speed))

	if *verify {
		log.Printf("Calculating checksum for verification...")
		receivedChecksum, err := calculateChecksum(filePath)
		if err != nil {
			log.Printf("⚠ Warning: Could not calculate checksum: %v", err)
		} else {
			log.Printf("✓ File checksum: %s", receivedChecksum)
			log.Printf("Note: Compare with sender's checksum to verify integrity")
		}
	}
}

func calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func copyWithProgress(dst io.Writer, src io.Reader, total int64, operation string) (int64, error) {
	var written int64
	buf := make([]byte, 32*1024)
	lastUpdate := time.Now()

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}

			if time.Since(lastUpdate) > 100*time.Millisecond {
				if total > 0 {
					percent := float64(written) / float64(total) * 100
					fmt.Printf("\r%s: %.1f%% (%s/%s)", operation, percent,
						formatBytes(written), formatBytes(total))
				} else {
					fmt.Printf("\r%s: %s", operation, formatBytes(written))
				}
				lastUpdate = time.Now()
			}
		}
		if er != nil {
			if er != io.EOF {
				return written, er
			}
			break
		}
	}
	fmt.Println()
	return written, nil
}

func copyWithProgressUnknownSize(dst io.Writer, src io.Reader, operation string) (int64, error) {
	return copyWithProgress(dst, src, 0, operation)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func calculateSpeed(bytes int64, duration time.Duration) int64 {
	if duration.Seconds() == 0 {
		return 0
	}
	return int64(float64(bytes) / duration.Seconds())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
