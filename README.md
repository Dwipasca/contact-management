# CLI Contact Management With Go

A simple command-line Contact Manager written in pure Go (Golang) with no third-party dependencies. It supports basic CRUD operations and can import/export contacts to/from JSON and CSV formats.

## Features

- Add, edit, and delete contacts
- Add multiple contacts at once
- List all contacts
- Search contacts by id, name, or email
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
    go build -o contact-management-app ./cmd/main.go
```

- run application

```bash
    ./contact-management-app
```

## Usage

The application provides a menu-driven interface with the following options:

1. Add Contact
2. Add Multiple Contacts
3. Edit Contact
4. Delete Contact
5. Show List Contact
6. Search Contact
7. Export Contacts
8. Import Contacts
9. Exit

Follow the on-screen prompts to use each feature.

## Project Structure

This project follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) guidelines:

- `/cmd` - Main applications
- `/internal` - Private application code
  - `/domain` - Domain models
  - `/repository` - Data access layer
  - `/usecase` - Business logic
  - `/handler` - UI handlers
- `/ui` - User interface utilities
