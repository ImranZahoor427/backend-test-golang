# Banking Ledger

A lightweight banking ledger system built with Go, using PostgreSQL for account storage, MongoDB for transaction logs, and Kafka for async transaction processing.

## Features

- Create and fetch accounts with balance tracking
- Atomic balance updates with overdraft protection
- Record transactions asynchronously via Kafka
- REST API with Gin
- GORM for PostgreSQL, official Mongo driver for MongoDB

## Tech Stack

- **Go** + **Gin** – RESTful API
- **PostgreSQL** – Account and balance storage
- **MongoDB** – Transaction log storage
- **Kafka** – Async transaction processing
- **GORM** – ORM for PostgreSQL
- **Ginkgo** + **Gomega** – Testing framework

## Setup

### Prerequisites

- Go ≥ 1.20
- PostgreSQL
- MongoDB
- Kafka (locally or via Docker)

### Environment Variables

Before running the application, set up your environment variables by creating a `.env` file.

You can do this by copying the provided `.env.example`:

### Run with Docker Compose

```bash
docker compose up --build
```

### Run application locally

```bash
go run cmd/main.go
```

## Sample `curl` Requests

### Create an Account

```bash
curl --location 'http://localhost:8080/api/v1/accounts' \
--header 'Content-Type: application/json' \
--data '{
    "owner_name": "Alice",
    "initial_balance": 1000
}'
```

### Get Account by ID

```bash
curl --location 'http://localhost:8080/api/v1/accounts/18902ef3-1d70-48f9-b497-a1c10f2fe38f'
```

### Create a Transaction (Deposit / Withdrawal)

```bash
curl --location 'http://localhost:8080/api/v1/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "account_id": "18902ef3-1d70-48f9-b497-a1c10f2fe38f",
    "amount": 571000,
    "type": "deposit"
}'
```

### Get transactions

```bash
curl --location 'http://localhost:8080/api/v1/transactions/account/18902ef3-1d70-48f9-b497-a1c10f2fe38f?limit=20&offset=0' \
--header 'Content-Type: application/json'
```