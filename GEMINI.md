# Webhook Notification Service

## Project Overview
This project is a Go-based webhook notification service designed to receive webhooks and forward notifications to messaging platforms. Currently, it focuses on integration with **Feishu (Lark)**.

The project follows a layered architectural pattern (inspired by Clean Architecture or DDD), organizing code into `domain`, `infrastructure`, and `application` layers.

### Key Technologies
- **Language**: Go 1.25.0
- **SDK**: [Lark Open API SDK (Go)](https://github.com/larksuite/oapi-sdk-go/v3)
- **Environment**: Configured via `.env` file

## Project Structure
- `internal/domain/`: Contains core business logic and interfaces (e.g., `MessageSender`).
- `internal/infrastructure/`: Concrete implementations of domain interfaces and external integrations.
    - `feishu/`: LarkBot implementation for sending messages via Feishu.
    - `config/`: Configuration management logic.
- `internal/application/`: Application-level orchestration (currently empty).

## Building and Running

### Prerequisites
- Go 1.25.0 or later.
- A `.env` file with the following keys:
  - `FEISHU_APP_ID`
  - `FEISHU_APP_SECRET`
  - `FEISHU_OPEN_ID`

### Commands
- **Install dependencies**: `go mod download`
- **Build**: `go build ./...`
- **Test**: `go test ./...`
- **Run**: (Note: A `main.go` entry point is not yet implemented)

## Development Conventions
- **Layered Architecture**: Keep business logic in the `domain` layer and external details (like API clients) in the `infrastructure` layer.
- **Interfaces**: Define interfaces in the `domain` layer to allow for multiple implementations (e.g., adding more notification providers).
- **Configuration**: Use `internal/infrastructure/config` for centralized configuration access.
