# Backend Recipe Project

This is a simple recipe project for calculate COGS based on the recipe ingredients.

# Table of Contents

- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Project Setup](#project-setup)
- [Database Migrations & Seeding](#database-setup-migrations-and-seeding)
- [Running the Project](#running-the-project)
- [Running the Tests](#running-the-tests)
- [API Documentation](#api-documentation)

## Prerequisites

- [Go](https://golang.org/dl/) (v1.23 or later recommended)
- Postgresql (v13 or later recommended)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/) (optional)

## Project Structure

```
project-root/
│── src/
│   ├── config/        # Configuration files (e.g., database, environment variables)
│   ├── controller/    # Handles HTTP requests and responses
│   ├── db/            # Database connection and queries
│   ├── middleware/    # Middleware functions (e.g., authentication, logging)
│   ├── model/         # Data models for database entities
│   ├── router/        # API route definitions
│   ├── seeder/        # Database seeding scripts
│   ├── utils/         # Utility functions and helpers
│── .env               # Environment variables
│── .env.example       # Example environment variable file
│── .gitignore         # Git ignore rules
│── docker-compose.yml # Docker Compose configuration
│── Dockerfile         # Docker image configuration
│── go.mod             # Go module file
│── main.go            # Application entry point
│── README.md          # Project documentation
```

## Project Setup

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/legenhand/recipe-test.git
   cd recipe-test
   
2. **Install Dependencies:**

   ```bash
   go mod tidy
   ```
   
3. **Create a `.env` file:**

   ```bash
   cp .env.example .env
   ```
   
    Update the `.env` file with your database credentials.
   ```dotenv
   DB_HOST=your_database_host
   DB_PORT=your_database_port
   DB_USER=your_database_user
   DB_PASSWORD=your_database_password
   DB_NAME=your_database_name
   JWT_SECRET=your_jwt_secret
   BASE_URL=http://localhost:8080
```

## Database Setup, Migrations and Seeding
you can run this command to run postgresql in docker
```bash
docker compose up -d
```

Run Database Seeding

```bash
go run main.go seed
```

## Running the Project
You can run the project with:
```bash
go run main.go
```

## Running the Tests
You can run the tests with:
```bash
go test ./...
```

## Build and Run in Production with Docker

1. Build the Docker image:

   ```bash
   docker build -t recipe-project:latest .
   ```
2. Run the container

   ```bash
   docker run -d -p 8080:8080 --env-file .env recipe-project:latest
   ```



## API Documentation

Here is the [API Documentation](https://www.postman.com/hoshikuzu/recipe/collection/d4an0em/recipe-test?action=share&creator=16096015
) for the project.