# Labs-HR-Go

## Table of Contents
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Docker Setup](#docker-setup)
- [Migrations](#migrations)
  - [Prerequisites](#prerequisites-1)
  - [Local Database Setup](#local-database-setup)
  - [Creating Migrations](#creating-migrations)
  - [Running Migrations](#running-migrations)
  - [Migration Best Practices](#migration-best-practices)
  - [Troubleshooting](#troubleshooting)
  - [References](#references)
- [API Documentation](#api-documentation)
  - [Employee Endpoints](#employee-endpoints)
  - [Attendance Endpoints](#attendance-endpoints)
- [All Environment Variables](#all-environment-variables)
  - [Server Configuration](#server-configuration)
  - [Database Configuration](#database-configuration)
  - [Application Features](#application-features)
  - [Usage Examples](#usage-examples)
- [Development Resources](#development-resources)

## Installation

### Prerequisites

- Go 1.24 or higher
- [GVM](https://github.com/moovweb/gvm) (optional, for managing Go versions)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/WangWilly/labs-hr-go.git
   cd labs-hr-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the development server:
   ```bash
   ./scripts/dev.sh
   ```

### Docker Setup

1. Build the Docker image using the provided build script:
   ```bash
   ./scripts/build.sh
   ```
   This will create a Docker image named `labs-hr-go:latest`.

2. Run the service using Docker Compose:
   ```bash
   cd deployments
   docker compose up -d
   ```
   This will start the service in detached mode, listening on port 8080.

3. To stop the service:
   ```bash
   docker compose down
   ```

4. Monitor logs:
   ```bash
   docker compose logs -f
   ```

## Migrations

### Prerequisites

- Docker and Docker Compose for local database setup

### Local Database Setup

Set up the local database environment using the provided script:

```bash
./scripts/local_setup.sh
```

This script will:
1. Navigate to the deployment directory
2. Verify Docker is running
3. Check if Docker Compose is installed
4. Start the required database containers

### Creating Migrations

1. Create a new migration file in the `database/migrations` directory following the naming convention `NNNNN_description.go`, where `NNNNN` is a sequential number.

   Example structure ([see example](./database/migrations/00001_init.go)):
   ```go
   package migrations

   import (
       "context"

       "github.com/WangWilly/labs-hr-go/pkgs/models"
       "github.com/WangWilly/labs-hr-go/pkgs/utils"
   )

   // UpNNNNNDescription runs the "up" migration
   func UpNNNNNDescription(ctx context.Context) error {
       // Migration logic to apply changes
       return nil
   }

   // DownNNNNNDescription runs the "down" migration
   func DownNNNNNDescription(ctx context.Context) error {
       // Migration logic to revert changes
       return nil
   }
   ```

2. In the `Up` function, implement the changes you want to apply to the database (create tables, add columns, etc.).

3. In the `Down` function, implement the logic to revert those changes (drop tables, remove columns, etc.).

### Running Migrations

Execute migrations using the migration command:

```bash
go run database/cmd/main.go
```

### Migration Best Practices

1. Always create both `Up` and `Down` functions for each migration
2. Test migrations in development before applying to production
3. Keep migrations idempotent when possible
4. Use transactions for complex migrations to ensure atomicity

### Troubleshooting

If you encounter issues with migrations:

1. Check the database connection settings
2. Verify that the migration files are properly formatted
3. Look for error messages in the application logs
4. Reset the database if needed during development

### References

- [Gormigrate Documentation](https://github.com/go-gormigrate/gormigrate)
- [Goose Documentation](https://github.com/pressly/goose)
- [GORM Migrations](https://gorm.io/docs/migration.html)
- [Stack Overflow: GORM Migration with golang-migrate](https://stackoverflow.com/questions/64510093/gorm-migration-using-golang-migrate-migrate)

## API Documentation

The HR management system provides RESTful APIs for managing employees and attendance records. All requests and responses use JSON format.

### Employee Endpoints

#### Create Employee

Creates a new employee record with their position information.

```bash
curl --location 'http://localhost:8080/employee' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Will",
    "age": 39,
    "address": "united states",
    "phone": "654321232",
    "email": "test@goooo.co",
    "position": "tester",
    "department": "tech",
    "salary": 4000,
    "start_date": 1746365072
}'
```

Response (200 OK):
```json
{
   "employee_id": 1,
   "position_id": 1
}
```

Request Parameters:
- `name` (string, required): Employee's full name
- `age` (integer, required): Employee's age
- `address` (string, required): Employee's address
- `phone` (string, required): Contact phone number
- `email` (string, required): Contact email address
- `position` (string, required): Job position title
- `department` (string, required): Department name
- `salary` (number, required): Monthly salary amount
- `start_date` (unix timestamp, required): Employment start date

Error Responses:
- 400 Bad Request: Invalid request format or missing required fields
- 500 Internal Server Error: Server-side processing error

#### Get Employee

Retrieves detailed information about an employee by ID.

```bash
curl --location 'http://localhost:8080/employee/1'
```

Response (200 OK):
```json
{
   "employee_id": 1,
   "name": "Will",
   "age": 39,
   "phone": "654321232",
   "email": "test@goooo.co",
   "address": "united states",
   "created_at": "2025-05-04 13:26:51",
   "updated_at": "2025-05-04 13:26:51",
   "position_id": 1,
   "position": "tester",
   "department": "tech",
   "salary": 4000,
   "start_date": "2025-05-04 00:00:00"
}
```

Error Responses:
- 400 Bad Request: Invalid ID format
- 404 Not Found: Employee not found

#### Update Employee

Updates an existing employee's information. All fields are optional - only include fields you want to update.

```bash
curl --location --request PUT 'http://localhost:8080/employee/1' \
--header 'Content-Type: application/json' \
--data '{
    "address": "taiwan"
}'
```

Response (200 OK):
```json
{
   "id": 1,
   "name": "Will",
   "age": 39,
   "address": "taiwan",
   "phone": "654321232",
   "email": "test@goooo.co"
}
```

Request Parameters:
- `name` (string, optional): Updated employee name
- `age` (integer, optional): Updated employee age
- `address` (string, optional): Updated address
- `phone` (string, optional): Updated phone number
- `email` (string, optional): Updated email address

Error Responses:
- 400 Bad Request: Invalid ID format or request body
- 404 Not Found: Employee not found
- 500 Internal Server Error: Update operation failed

#### Promote Employee

Updates an employee's position, department, or salary information.

```bash
curl --location 'http://localhost:8080/promote/1' \
--header 'Content-Type: application/json' \
--data '{
    "position": "tester2",
    "department": "tech",
    "salary": 5000,
    "start_date": 1747365072
}'
```

Response (200 OK):
```json
{
   "position_id": 5,
   "start_date": "2025-05-16 11:11:12"
}
```

Request Parameters:
- `position` (string, required): New position title
- `department` (string, required): New department name
- `salary` (number, required): New salary amount
- `start_date` (unix timestamp, required): When the promotion takes effect

Error Responses:
- 400 Bad Request: Invalid request format or missing required fields
- 404 Not Found: Employee not found
- 500 Internal Server Error: Promotion operation failed

### Attendance Endpoints

#### Clock In

Records an attendance entry when an employee starts work.

```bash
curl --location 'http://localhost:8080/attendance' \
--header 'Content-Type: application/json' \
--data '{
    "employee_id": 1
}'
```

Response (200 OK):
```json
{
   "attendance_id": 1,
   "position_id": 3,
   "clock_in_time": "2025-05-04 13:41:15",
   "clock_out_time": ""
}
```

Request Parameters:
- `employee_id` (integer, required): ID of the employee clocking in

Error Responses:
- 400 Bad Request: Invalid employee ID or already clocked in
- 404 Not Found: Employee not found
- 500 Internal Server Error: Clock-in operation failed

#### Clock Out

Records when an employee ends their workday.

```bash
curl --location --request PUT 'http://localhost:8080/attendance/1' \
--header 'Content-Type: application/json'
```

Response (200 OK):
```json
{
    "attendance_id": 1,
    "position_id": 3,
    "clock_in_time": "2025-05-04 13:41:15",
    "clock_out_time": "2025-05-04 17:30:22"
}
```

Error Responses:
- 400 Bad Request: Invalid ID format or already clocked out
- 404 Not Found: Attendance record not found
- 500 Internal Server Error: Clock-out operation failed

#### Get Attendance Record

Retrieves an attendance record by ID.

```bash
curl --location 'http://localhost:8080/attendance/1'
```

Response (200 OK):
```json
{
    "attendance_id": 1,
    "position_id": 3,
    "clock_in_time": "2025-05-04 13:41:15",
    "clock_out_time": "2025-05-04 17:30:22"
}
```

Error Responses:
- 400 Bad Request: Invalid ID format
- 404 Not Found: Attendance record not found

## All Environment Variables

### Server Configuration
| Name | Description | Default |
|------|-------------|---------|
| PORT | The port on which the service listens | `8080` |
| HOST | The host address for the service | `0.0.0.0` |

### Database Configuration
| Name | Description | Default |
|------|-------------|---------|
| DB_HOST | Database server hostname | - |
| DB_PORT | Database server port | `3306` |
| DB_USER | Database username | - |
| DB_PASSWORD | Database password | - |
| DB_DATABASE | Database name | - |

### Application Features
| Name | Description | Default |
|------|-------------|---------|
| DB_SEED | Whether to seed the database with sample data | `false` |

### Usage Examples

#### Local Development
```bash
# Run with database seeding enabled
DB_SEED=true go run cmd/main.go

# Custom database connection
DB_HOST=localhost DB_PORT=3306 DB_USER=myuser DB_PASSWORD=mypass DB_DATABASE=mydb go run cmd/main.go
```

#### Docker Environment
When using Docker Compose, configure these variables in the `deployments/docker-compose.yml` file.

## Development Resources

- [Go Modules Documentation](https://go.dev/wiki/Modules#quick-start)
- https://github.com/smartystreets/goconvey
- https://github.com/uber-go/mock
- https://github.com/rs/zerolog?tab=readme-ov-file#benchmarks
- https://github.com/uber-go/zap
