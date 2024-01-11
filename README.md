# Go-Web-CRM

Go-Web-CRM is a simple Customer Relationship Management (CRM) system written in Go, leveraging an SQLite database. The project enables users to perform basic CRM functions through an intuitive web interface.

## Features

- Add new clients to the database.
- Browse a list of clients with a user-friendly interface.
- Basic client list filtering functionality.
- Remove clients from the database.

## Usage Instructions

1. Compile the program using the command `go build -o app.exe`.
2. Run the program with `./app.exe`.
3. Open a web browser and visit [http://localhost:8080](http://localhost:8080).

## Remember about database

To ensure proper functionality, please remember that the database file (database.db) must be located in the same directory as the executable (app.exe/main.go). Make sure to keep both files in the same folder to enable seamless interaction with the database.
