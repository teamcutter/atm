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
   go build
    ```
3. Run the server:
   ```bash
   ./atm -pass 12345 -login user -p :8001 
    ``` 

### Usage

1. **Start the server**:  
   Run the ATM server as described above.

2. **Use independent client with desired protocol**

### TCP-Based Protocol Specification

The server already uses a simple TCP-based protocol. Let’s formalize it for distribution:

#### 1. Authentication

- **Format**:  
  `login:password\n`  
  (Plaintext string, terminated by a newline `\n`.)

- **Example**:  
  `user:12345\n`

- **Response**:
  - **Success**:  
    `OK\n`
  - **Failure**:  
    `ERROR: <message>\n`  
    _(e.g., `ERROR: invalid login or password\n`)_

---

#### 2. Commands

- **Format**:  
  `<header><keyLen><key><valueLen><value>`  
  (Binary, no separators except for lengths.)

  - `<header>`: 3 bytes (e.g., `SET`, `GET`, `DEL`).
  - `<keyLen>`: 4 bytes (uint32, big-endian), length of the key in bytes.
  - `<key>`: Variable-length string (UTF-8 encoded).
  - `<valueLen>`: 4 bytes (uint32, big-endian), length of the value in bytes (**only for SET**).
  - `<value>`: Variable-length string (UTF-8 encoded, **only for SET**).
  - Terminated by `\n`.

#### 3. Examples

- **SET user 1**:
  - **Hex**:  
    `53455400000004757365720000000131\n`
  - **Breakdown**:  
    `SET (3)` + `00000004 (4)` + `user (4)` + `00000001 (4)` + `1 (1)` + `\n (1)`

- **GET user**:
  - **Hex**:  
    `4745540000000475736572\n`
  - **Breakdown**:  
    `GET (3)` + `00000004 (4)` + `user (4)` + `\n (1)`

- **DEL user**:
  - **Hex**:  
    `44454c0000000475736572\n`
  - **Breakdown**:  
    `DEL (3)` + `00000004 (4)` + `user (4)` + `\n (1)`

---

#### 4. Response

- **Success**:  
  `<command> <key> = <value>\n`  
  _(e.g., `SET user = 1\n`, `GET user = 1\n`)_

- **Failure**:  
  `ERROR: <message>\n`  
  _(e.g., `ERROR: no record with such key\n`)_

---

#### 5. Notes

- All lengths are in **bytes**, not characters.
- The server expects **binary data** for commands, **not text**, but **responses are human-readable text**.