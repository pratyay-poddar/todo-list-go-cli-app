# Simple To-Do List App

A minimalist, high-performance Command Line Interface (CLI) application built from scratch using Go (Golang) and SQLite. 

## 🎯 Project Purpose
This repository was created purely for **syntax practice and hands-on learning**. The core goal of this project was to strengthen my understanding of foundational Go concepts, including:
* Structs and custom slice methods
* Memory management using the `defer` keyword
* Handling pointer semantics safely
* Integrating raw SQL queries with a database driver (`modernc.org/sqlite`)

---

## 📁 File Architecture

The project follows a clean separation of concerns, dividing the core business logic from the database orchestration layer:

```text
simple-to-do/
├── cmd/
│   └── main.go         # Application Orchestrator (CLI Parsing & SQL Queries)
├── todo/
│   └── todo.go         # Core Business Logic (Item/Items Structs & Domain Methods)
├── database/
│   └── database.db     # SQLite Database File (Auto-generated on initialization)
├── go.mod              # Package Dependency Manifest
└── README.md           # Project Documentation Manual
