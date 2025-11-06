# Weather Service

A Go service that retrieves current temperature by city name using the [Open-Meteo API](https://open-meteo.com/).

## Overview

Weather Service is an HTTP service that periodically fetches weather data from Open-Meteo API, stores historical readings in PostgreSQL, and provides a REST API to query the latest temperature for any city.

## Features

- Retrieve current temperature by city name
- Automatic periodic updates (configurable interval)
- Historical data storage in PostgreSQL
- Support for any city via Open-Meteo Geocoding API

## Requirements

- Go 1.25.1 or higher
- Docker and Docker Compose
- Internet access for Open-Meteo API

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/idakhno/weather-service.git
cd weather-service
```

### 2. Start the database

Start PostgreSQL using Docker Compose. This will automatically create the `weather` database and apply migrations:

```bash
docker compose up -d
```

Verify the container is running:

```bash
docker ps
```

### 3. Run the service

```bash
go run cmd/server/main.go
```

The service starts on port `3000` and begins automatically updating temperature data for Moscow every 10 seconds.

### 4. Verify

In a new terminal, make a request:

```bash
curl http://localhost:3000/moscow
```

Example response:

```json
{
  "Time": "2024-01-15T14:30:00Z",
  "Temperature2m": -5.2
}
```

Note: If no data is available yet, wait a few seconds as the background task updates temperature every 10 seconds.

## API Usage

### Get Temperature

**Request:**
```bash
GET /{city}
```

**Examples:**
```bash
curl http://localhost:3000/moscow
curl http://localhost:3000/saint-petersburg
curl http://localhost:3000/novosibirsk
```

**Response (200 OK):**
```json
{
  "Time": "2024-01-15T14:30:00Z",
  "Temperature2m": -5.2
}
```

**Response (404 Not Found):**
```
not found
```

## Configuration

Current settings are hardcoded in `cmd/server/main.go`:

- **Server port:** `3000`
- **Database:** `postgresql://postgres:postgres@localhost:5432/weather`
- **Update interval:** `10 seconds`
- **Default city:** `moscow`

For production, these should be moved to environment variables or a configuration file.

## Database Schema

The `reading` table stores historical weather measurements:

```sql
CREATE TABLE reading (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP NOT NULL
);
```

## Development

### Install dependencies

```bash
go mod download
```

### Run tests

```bash
go test ./...
```

### Stop the database

```bash
docker compose down
```

To remove all data:

```bash
docker compose down -v
```

## License

MIT © 2025
