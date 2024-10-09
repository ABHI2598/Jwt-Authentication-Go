# Go JWT Authentication Project

## Overview

This project demonstrates how to implement JWT (JSON Web Token) authentication in a Go (Golang) application. It provides basic APIs for user registration, login, and protected routes accessible only with a valid JWT.

## Features

- **User Registration**: Allows users to register with a username and password.
- **User Login**: Validates user credentials and returns a signed JWT.
- **Protected Routes**: Routes that can only be accessed with a valid JWT.
- **JWT Middleware**: Middleware to verify JWT tokens on protected routes.

## Tech Stack

- **Go**: Main programming language.
- **JWT (JSON Web Tokens)**: For stateless authentication.
- **GIN**: HTTP request router for handling routes.
- **github.com/dgrijalva/jwt-go**: JWT library for creating and verifying tokens.
- **MongoDB**: For persistent user storage.

## Setup

### Prerequisites

Ensure you have the following installed:

- Go (v1.16 or higher)
- A database (e.g., PostgreSQL/MySQL/SQLite,MongoDB) for storing user data.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ABHI2598/Jwt-Authentication-Go.git
   cd Jwt-Authentication-Go

2. Install dependencies:
   ```bash
   go mod init main.go
   go mod tidy
   
3. Set up environment variables: Create a .env file in the project root to store sensitive information like JWT secrets and MongoDB connection details.
   
4. Run the application:
   ```bash
   go run main.go
5. The API will now be accessible at respective port defined in .env file.
   ```bash
   http://localhost:8080


## License
   This project is licensed under the MIT License. See the LICENSE file for more details.
