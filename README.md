# Storage API Server

A high-performance CRUD RESTful API server for local machine file management built in Go using `net/http` and `charlievieth/fastwalk`.

> [!CAUTION]
> This project was written almost by AI

## Features

- **File/Folder Operations**: List, search recursively, create, rename, delete, download, and upload files.
- **SSE Upload Progress**: Real-time server-sent events indicating percent complete and bytes written during file uploads.
- **Race Condition Safety**: Path-keyed locking (`sync.Map` of mutexes) ensures operations on the same path are fully serialized while unrelated paths remain parallel.
- **Security Checkers**:
  - Rejects symlinks and hardlinks.
  - Validates paths against escaping via path-traversal checking.
  - Checks for duplicates and validates form data inputs.
  - Implements request body limits (`MaxBytesReader`) on form inputs to prevent memory exhaustion.
- **Conditional Logging Middleware**: Build-tag gated logging (`//go:build !production`) logs format `METHOD /path {time}ms`.

---

## Getting Started

### Prerequisites

- Go (1.22+ recommended)
- `golangci-lint` (v2.x) for code analysis

### Building

Build the server binary:

```bash
make
```

Run tests:

```bash
make test
```

Clean build artifacts:

```bash
make clean
```

---

## CLI Usage

Start the server using the compiled binary by passing the storage path:

```bash
./bin/storage [flags] /path/to/storage
```

### Flags

- `-p`, `--port`: Port to listen on (overrides the `PORT` environment variable; defaults to `8080`).

---

## API Documentation

### 1. List directory

- **Route**: `GET /src/{path...}`
- **Response**: List of `Item` JSON objects.

### 2. Search files/folders recursively

- **Route**: `GET /src/{path...}?q=search_term`
- **Response**: List of matching `Item` JSON objects.

### 3. Create folder

- **Route**: `POST /src/{path...}`
- **Form Data**: `newName` (name of directory)
- **Response**: `{"path": "relative/path/to/new_folder"}`

### 4. Rename item

- **Route**: `PUT /src/{path...}`
- **Form Data**: `newName` (new base name)
- **Response**: `{"path": "new_relative_path"}` (for folder) or `{"status": "oke"}` (for file)

### 5. Delete item

- **Route**: `DELETE /src/{path...}`
- **Response**: `{"status": "oke"}`

### 6. Upload file (with SSE progress)

- **Route**: `POST /upload/{path...}`
- **Form Data**:
  - `name`: target name of the file
  - `file`: binary contents of the file
- **Response**: Server-Sent Events stream:
  - Event payload: `data: {"percent": 45, "bytesWritten": 4500000, "totalBytes": 10000000}`
  - Final success event payload: `data: {"status": "oke", "file": "filename.ext"}`

### 7. Download file

- **Route**: `GET /download/{path...}`
- **Response**: Serves the raw file content.
