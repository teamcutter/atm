![ATM Machine](https://cdn-icons-png.flaticon.com/512/6059/6059866.png)

# ATM: A Lightweight In-Memory Database

**ATM** is a lightweight, high-performance in-memory database written in Go, designed to provide key-value storage similar to Redis. It supports efficient storage and retrieval of data, enabling developers to build fast and scalable applications.

---

## Features

- **Key-Value Storage**: Store and retrieve data with ease using a simple key-value mechanism.
- **High Performance**: Built for speed with an emphasis on low-latency data operations.
- **Concurrency**: Leverages Go’s concurrency model for handling multiple requests simultaneously.
- **Persistence (Optional)**: Support for saving in-memory data to disk for durability.
- **Minimal Dependencies**: Built using Go's standard library for simplicity and reliability.
- **Extensible Protocol**: Easily add new features or integrate with existing systems.

---

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/teamcutter/atm.git
   cd atm
    ```

2. Build the application:
   ```bash
   make build
    ```
3. Run the server:
   ```bash
   ./atm -p 8001 # Replace 8000 with your desired port.
    ``` 

### Usage

1. **Start the server**:  
   Run the ATM server as described above.

2. **Interact with the database**:  
   Use any TCP client (e.g., `netcat` or custom tools) to interact with the server. Example commands:

   - **Set a key-value pair**:
     ```bash
     SET mykey myvalue
     ```

   - **Retrieve a value by key**:
     ```bash
     GET mykey
     ```

   - **Delete a key**:
     ```bash
     DEL mykey
     ```

3. **Custom Protocol**:  
   The server communicates using a lightweight, custom protocol for simplicity and efficiency.
