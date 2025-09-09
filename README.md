# P2P File Sharing

A peer-to-peer file sharing application that enables direct file transfer between devices over TCP connections with streaming.


https://github.com/user-attachments/assets/faa4a279-79c2-4e1a-aa0e-f1e8f3445893


## Features

- **Direct P2P Communication**: No centralized server required - files are transferred directly between peers
- **TCP Streaming**: Efficient file transfer using TCP protocol with chunked streaming
- **CRC Validation**: Built-in Cyclic Redundancy Check (CRC) validation ensures data integrity during transfer
- **Simple CLI Interface**: Easy-to-use command-line interface for sending and receiving files
- **Port Flexibility**: Configure custom ports for file transfer operations

## How It Works

The application operates on a simple client-server model where:

- One peer acts as the **sender** (server) and shares a file on a specified port
- Another peer acts as the **receiver** (client) and connects to download the file
- Files are transmitted in chunks over TCP with CRC validation for each chunk
- Both peers must be on the same network or have network connectivity

## Installation

```bash
# Clone the repository
git clone https://github.com/rushikeshg25/p2p-file-sharing.git
cd p2p-file-sharing

# Install dependencies (if any)
# Add installation instructions based on your project setup
```

## Usage

### Sending a File

To share a file, use the `send` command:

```bash
p2p-share send <filename> <port>
```

**Example:**

```bash
p2p-share send 50mb.mov 3001
```

This command:

- Starts a server on port 3001
- Makes the file `50mb.mov` available for download
- Waits for incoming connections from receivers

### Receiving a File

To download a file, use the `receive` command:

```bash
p2p-share receive <filename> <port>
```

**Example:**

```bash
p2p-share receive filename 3001
```

This command:

- Connects to a sender on port 3001
- Downloads the file
- Saves the file to the current directory

## Technical Details

### Protocol

- **Transport Layer**: TCP for reliable data transmission
- **Data Integrity**: CRC (Cyclic Redundancy Check) validation for each file chunk
- **Streaming**: Files are transmitted in chunks to handle large files efficiently

### Network Requirements

- Both sender and receiver must have network connectivity
- Firewall rules may need to be configured to allow connections on specified ports
- For local network transfers, ensure both devices are on the same subnet

## Example Workflow

1. **Sender Side**:

   ```bash
   p2p-share send document.pdf 3001
   # Server starts and waits for connections on port 3001
   ```

2. **Receiver Side**:
   ```bash
   p2p-share receive document.pdf 3001
   # Connects to sender and downloads document.pdf
   ```

## Error Handling

- **CRC Validation**: Each chunk is validated using CRC checksums
- **Connection Errors**: Automatic retry mechanisms for network interruptions
- **File Integrity**: Complete file validation after transfer completion

## Port Configuration

- Default ports can be customized based on your network setup
- Ensure the specified port is not blocked by firewall rules
- Use ports above 1024 to avoid requiring administrator privileges
