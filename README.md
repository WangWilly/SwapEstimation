# Swap Estimation

## Table of Contents
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Docker Setup](#docker-setup)
- [API Documentation](#api-documentation)
- [All Environment Variables](#all-environment-variables)
  - [Server Configuration](#server-configuration)
- [Development Resources](#development-resources)

## Installation

### Prerequisites

- Go 1.24 or higher
- [GVM](https://github.com/moovweb/gvm) (optional, for managing Go versions)

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/WangWilly/SwapEstimation.git
   cd SwapEstimation
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
   This will create a Docker image named `swap-estimation:latest`.

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

## API Documentation

### Endpoints

```bash
curl --location 'http://localhost:8080/estimate?pool=0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852&src=0xdAC17F958D2ee523a2206206994597C13D831ec7&dst=0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2&src_amount=10000000'
```

Response (200 OK):
```json
3902524309783809
```

Request Parameters:
- `pool` (string, required): The address of the liquidity pool
- `src` (string, required): The source token address
- `dst` (string, required): The destination token address
- `src_amount` (number, required): The amount of source token to swap

Error Responses:
- 400 Bad Request: Invalid request format or missing required fields
- 500 Internal Server Error: Server-side processing error

## All Environment Variables

### Server Configuration
| Name | Description | Default |
|------|-------------|---------|
| PORT | The port on which the service listens | `8080` |
| HOST | The host address for the service | `0.0.0.0` |

### Usage Examples

#### Docker Environment
When using Docker Compose, configure these variables in the `deployments/docker-compose.yml` file.

## Development Resources

- [Go Modules Documentation](https://go.dev/wiki/Modules#quick-start)
- https://github.com/smartystreets/goconvey
- https://github.com/uber-go/mock
- https://github.com/rs/zerolog?tab=readme-ov-file#benchmarks
- https://github.com/uber-go/zap
