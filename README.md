# Go Zero-Knowledge Proof API

This repository provides an implementation of a zero-knowledge proof system in Go, allowing nodes to prove the possession of certain data without revealing the data itself. The system uses SHA-256 for hashing and offers an API for data broadcast, receipt, and acknowledgment between nodes.

This repository also makes use of The http/2 API Server located at https://github.com/AndrewDonelson/go-http2-api

## Project Structure

- `pkg/zeroknowledge.go`: Contains the core logic for generating and verifying hashes.
- `pkg/api.go`: Contains the API logic, including endpoints for broadcasting data and acknowledgments.
- `cmd/server/main.go`: Main entry point for starting the example server.

## Installation and Setup

1. Clone the repository:
```bash
git clone https://github.com/AndrewDonelson/go-zeroknowledge-proof.git
cd go-zeroknowledge-proof
```

2. Install the required Go dependencies:
```bash
go mod tidy
```

## Running the Project

Navigate to the `cmd/server` directory and run:

```bash
go run main.go
```

The server will start on port 8080, offering HTTP/2 support.

## API Endpoints Overview

- `/prove`: Computes and returns the hash of given data.
- `/verify`: Compares provided and expected hashes.
- `/broadcast`: Broadcasts data to all authorized nodes.
- `/acknowledge`: Receives acknowledgments from nodes regarding data acceptance.

## Contribution

Feel free to submit pull requests or raise issues.

## License

[MIT License](LICENSE)
