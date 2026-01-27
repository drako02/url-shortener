<div align="center">

# 🔗 URL Shortener

**A production-ready URL shortening service with Firebase authentication and event-driven analytics**

[![Go Version](https://img.shields.io/badge/Go-1.23.3-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Firebase](https://img.shields.io/badge/Firebase-FFCA28?style=flat&logo=firebase&logoColor=black)](https://firebase.google.com)

[Features](#-features) • [Tech Stack](#-tech-stack) • [Quick Start](#-quick-start) • [API](#-api-overview)

</div>

---

## Motivation

Built to explore event-driven architecture and clean design patterns in Go. This project goes beyond basic URL shortening—it demonstrates asynchronous event processing with Kafka, secure authentication flows, and testable repository patterns. Perfect for understanding how to structure scalable backend services.

## Features

- **URL Shortening** – Generate unique short codes with automatic collision handling
- **Smart Redirection** – Fast redirects with click event tracking
- **Firebase Auth** – Token-based authentication with Google's infrastructure
- **User Management** – Per-user URL tracking with pagination support
- **Event-Driven Analytics** – Non-blocking Kafka integration for click metrics
- **Clean Architecture** – Testable layers (handlers → services → repositories)

## Tech Stack

**Backend** – Go, Gin Web Framework  
**Database** – PostgreSQL with GORM ORM  
**Authentication** – Firebase Admin SDK  
**Messaging** – Apache Kafka (Confluent)  
**Testing** – Testify, SQL Mock

## Learnt

- Implementing **repository pattern** for clean separation of concerns
- Using **Kafka for async processing** without blocking HTTP requests
- Writing **testable Go code** with dependency injection and mocks
- Integrating **Firebase Admin SDK** for production-grade auth
- Handling **database migrations** and relationships with GORM

## 🚀 Quick Start

**Prerequisites:** Go 1.23+, PostgreSQL, Firebase project

```bash
# Clone and install
git clone https://github.com/drako02/url-shortener.git
cd url-shortener
go mod download

# Configure environment (.env)
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=url_shortener

# Run
go run main.go
```

Server starts at `http://localhost:8080`

## API Overview

```http
POST   /create              # Create shortened URL
GET    /:shortCode          # Redirect to original URL
POST   /user-urls           # Get user's URLs (paginated)
POST   /users               # Create user
POST   /users/exists        # Check user existence
GET    /users/:uid          # Get user details
```

**Example:**
```bash
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "uid": "user123"}'
```

**Response:**
```json
{
  "short_code": "abc123",
  "original_url": "https://example.com",
  "created_at": "2026-01-26T10:00:00Z"
}
```

## Testing

```bash
go test ./... -v
```

Includes unit tests for services, repositories, and handlers with mocked dependencies.

## Future Plans

- [ ] Analytics dashboard with real-time metrics
- [ ] Custom domain support for branded links
- [ ] QR code generation
- [ ] Rate limiting and abuse prevention
- [ ] Link expiration and password protection

## Project Structure

```
├── handlers/       # HTTP request handlers
├── services/       # Business logic layer
├── repositories/   # Data access layer
├── models/         # Domain models
├── routes/         # Route definitions
├── config/         # Database, Firebase, Kafka setup
└── middlewares/    # Authentication middleware
```

---

<div align="center">

**Built by [drako02](https://github.com/drako02)**

[Report Bug](https://github.com/drako02/url-shortener/issues) • [View More Projects](https://github.com/drako02)

</div>