# APAC Backend

APAC Backend is a RESTful API built with Go, providing authentication, user management, AI-powered trip planning (Gemini), and integration with PostgreSQL, Redis, Supabase, and Google OAuth.

## Features

- **Authentication**: Register, login, OTP verification, Google OAuth, refresh token, logout.
- **User Management**: User profile, edit profile (including photo upload to Supabase), user preferences.
- **Trip Planner**: Generate travel itineraries using Gemini AI, save and manage trip history.
- **Integrations**: PostgreSQL (main database), Redis (OTP & state), Supabase (media storage), Google OAuth, Email (Gmail SMTP).
- **Middleware**: JWT authentication.
- **Monitoring**: `/metrics` endpoint for application monitoring.

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   ├── bootstrap/
│   ├── domain/
│   ├── infra/
│   └── middleware/
├── resource/
│   └── schema.json
├── Dockerfile
├── compose.yaml
├── Makefile
├── go.mod
├── .env
└── init.sql
```

## Getting Started

### 1. Environment Setup

- Copy the `.env` file and fill in the required variables (see `.env` for reference).
- Make sure to provide credentials for the database, Redis, Supabase, Google OAuth, and Email.

### 2. Run with Docker

```sh
make run
```

or manually:

```sh
docker compose up --build
```

### 3. Main Endpoints

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/verify-otp` - Verify OTP
- `POST /api/v1/auth/login` - User login
- `GET  /api/v1/auth/google` - Google OAuth login
- `GET  /api/v1/user/profile` - Get user profile (JWT required)
- `PATCH /api/v1/user/profile` - Edit user profile (JWT required)
- `POST /api/v1/user/preferences` - Update user preferences (JWT required)
- `POST /api/v1/gemini` - Generate trip itinerary with AI (JWT required)
- `GET  /api/v1/trips/` - List user trips (JWT required)
- `GET  /api/v1/trips/:id` - Trip details (JWT required)
- `DELETE /api/v1/trips/:id` - Delete trip (JWT required)
- `GET  /metrics` - Application monitoring

## Local Development

1. Install Go (see `go.mod` for version).
2. Start PostgreSQL and Redis (you can use Docker).
3. Run the application:
   ```sh
   go run ./cmd/api
   ```

---

**Notes:**

- The AI itinerary schema is defined in [resource/schema.json](resource/schema.json).
