# Swap Estimation

## Table of Contents
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Setup](#setup)
  - [Docker Setup](#docker-setup)
- [API Documentation](#api-documentation)
- [All Environment Variables](#all-environment-variables)
  - [Server Configuration](#server-configuration)
  - [Ethereum Connection](#ethereum-connection)
  - [Ethereum Client Configuration](#ethereum-client-configuration)
  - [Ethereum WebSocket Client Configuration](#ethereum-websocket-client-configuration)
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

### Ethereum Connection
| Name | Description | Default |
|------|-------------|---------|
| GETH_CLIENT_URL | Ethereum HTTP client URL | Required |
| GETH_WSS_CLIENT_URL | Ethereum WebSocket client URL | Required |

### Ethereum Client Configuration
| Name | Description | Default |
|------|-------------|---------|
| ETH_CLIENT_BLOCK_RANGE_SIZE | Maximum size of block range for querying | `9900` |

### Ethereum WebSocket Client Configuration
| Name | Description | Default |
|------|-------------|---------|
| ETH_WSS_CLIENT_LISTEN_PAIR_PERIOD | Period to refresh pair data | `2m` |

### Usage Examples

#### Docker Environment
When using Docker Compose, configure these variables in the `deployments/docker-compose.yml` file:

```yaml
version: '3'
services:
  backend:
    image: swap-estimation-app:latest
    environment:
      PORT: 8080
      HOST: 0.0.0.0
      GETH_CLIENT_URL: https://mainnet.infura.io/v3/your-project-id
      GETH_WSS_CLIENT_URL: wss://mainnet.infura.io/ws/v3/your-project-id
      ETH_CLIENT_BLOCK_RANGE_SIZE: 9900
      ETH_WSS_CLIENT_LISTEN_PAIR_PERIOD: 2m
```

## Development Resources

- [Go Modules Documentation](https://go.dev/wiki/Modules#quick-start)
- https://github.com/smartystreets/goconvey
- https://github.com/uber-go/mock
- https://github.com/rs/zerolog?tab=readme-ov-file#benchmarks
- https://github.com/uber-go/zap
- https://github.com/Uniswap/v2-sdk/blob/main/src/entities/pair.ts#L184
