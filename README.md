# CLI Contact Management With Go

A simple command-line Contact Manager written in pure Go (Golang) with no third-party dependencies. It supports basic CRUD operations and can import/export contacts to/from JSON and CSV formats.

## Features

- Add a single or multiple contacts
- Edit and delete existing contacts
- List all saved contacts
- Search contacts by ID, name, or email
- Export contacts to JSON or CSV
- Import contacts from JSON or CSV
- Interactive CLI interface using `bufio.Scanner`

## Getting Started

### Prerequisites

- [Go](https://go.dev/doc/install) version `1.24.1` or above recommended

### Installation

- clone repository

```bash
   git clone git@github.com:Dwipasca/contact-management.git
```

- navigate to project repository

```bash
    cd contact-management
```

- run the application

```bash
    go run cmd/.main.go
```

- build the binary (optional)

```bash
    go build -o contact-management
```
