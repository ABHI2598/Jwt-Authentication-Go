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
- **gin**: HTTP request router for handling routes.
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
   git clone https://github.com/yourusername/go-jwt-auth.git
   cd go-jwt-auth
