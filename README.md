
# Project Name: JETHUR

## Description
JETHUR is a Go language project structured to manage tasks and database migrations effectively. This backend service aims to provide a robust framework for task scheduling and execution with an emphasis on modularity and separation of concerns.

## Installation

### Prerequisites
- Go (Version specified in `go.mod`)
- Database (As configured in `database.go`)

### Setup# 🚀 JETHUR API – Setup Guide

This guide walks you through setting up and running the JETHUR API locally.

## 1. Clone the Repository
git clone [repository_url]  
cd JETHUR

## 2. Install Go Dependencies
go mod tidy

## 3. Configure Environment Variables
Copy the example environment file and update it:  
cp .env.example .env  

Edit `.env` with your environment-specific settings (database, server port, etc).  
Set `ROOT_PATH` to your project folder.

## 4. Database Migrations with Goose
We use Goose for managing migrations.

### Install Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

### Configure Migration Script
cp script.example.bat script.bat  
Edit `script.bat` with your database connection details.

## 5. Run the Application

### Option A: Normal Run
go run cmd/api/main.go

### Option B: Run with Air (hot reload for development)
1. Install Air:
   go install github.com/cosmtrek/air@latest

2. Create a `.air.toml` config file in the project root (if not already present).

3. Start the server with auto-reload:
   air

---

## ✅ Quick Setup Checklist
- Repository cloned  
- Dependencies installed (`go mod tidy`)  
- `.env` configured (with `ROOT_PATH`)  
- Goose installed & migration script updated  
- Server started successfully (via `go run` or `air`)  

```

## Project Structure

```
JETHUR/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── task/
│       └── task.go
├── config/
│   └── database.go
├── grc/
│   ├── common/
│   └── organization/
├── middleware/
│   └── recovery.go
├── migrations/
│   └── 20240213175253_create_error_logs_table.sql
├── tmp/
├── utility/
├── .air.toml
├── .env
├── go.mod
├── go.sum
└── script.bat
```

## Usage
Explain how to use the application, any available commands, and examples of common use cases.

## Contributing
Guidelines for how to make contributions to the project. It might include:

- Code style guides
- Git commit message guidelines
- Process for submitting a pull request

## License
Specify the license under which the project is made available.

## Acknowledgments
Credit any authors or inspiration.

## Contact Information
Provide information on how to reach out to the maintainer or the team for queries.

