# 🗂️ Distributed File Storage System

A **distributed file storage system** built in **Golang**, designed for **efficient, secure, and scalable peer-to-peer (P2P)** file storage and retrieval.  
This system ensures **data integrity, encryption, and optimized storage** using a **content-addressable storage (CAS)** model.

---

## 🚀 Features

### 🔗 Peer-to-Peer Architecture
- Fully decentralized file sharing using a **TCP-based transport layer**.  
- Each peer can act as both a **client and server**, supporting dynamic participation.

### 🧩 Content-Addressable Storage (CAS)
- Implements **SHA-1 hashing** for deterministic file identification.  
- Guarantees deduplication and efficient lookup of files by content hash.

### 🔒 Secure Encryption
- Files are **AES-encrypted** before distribution.  
- Supports **key-based decryption**, ensuring confidentiality and controlled access.

### ⚙️ Fault Tolerance & Scalability
- Allows **multiple nodes** to store and retrieve chunks of data.  
- Resilient to node failures and supports **horizontal scalability** as peers join or leave the network.

### 📁 Optimized Storage Management
- Structured **file path transformation** for organized storage.  
- Efficient **deletion and cleanup mechanisms** to maintain system health.

---


## 🧰 Tech Stack

- **Language:** Go (Golang)  
- **Networking:** TCP sockets  
- **Hashing:** SHA-1  
- **Encryption:** AES  
- **Storage Model:** Content-Addressable Storage (CAS)  

---

## ⚡ Getting Started

### Prerequisites
- Go 1.21 or higher
- Git

### Installation
```bash
git clone https://github.com/sarthak21-negi/distributed-file-storage.git
cd distributed-file-storage
go mod tidy

```
---

### Run a Peer Node
go run main.go
