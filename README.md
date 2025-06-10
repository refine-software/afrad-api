# Afrad API

Afrad is a RESTful API for a clothing e-commerce platform. It handles user authentication, product listings, cart management, order processing, and more.

---

## 🚀 Get Up and Running

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Goose](https://github.com/pressly/goose) (for migrations)

### Environment Variables

Create a `.env` file with the following required environment variables:

```env
# APP
PORT=8080
# use 'dev' for development and 'prod' for production
APP_ENV=dev
MAX_OTP_REQUESTS_PER_DAY=10
OTP_EXP_IN_MIN=5

# DB
DB_HOST="afrad_db"
DB_PORT=5432
DB_DATABASE="afrad_db"
DB_USERNAME="afrad_api"
DB_PASSWORD="afrad1020"
DB_SCHEMA=public

DATABASE_URL=postgres://afrad_api:afrad1020@localhost:5432/afrad_db?sslmode=disable

# S3
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=
S3_REGION=
S3_BUCKET=

# Oauth 2.0
SESSION_KEY=
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

# JWT
ACCESS_TOKEN_SECRET=
ACCESS_TOKEN_EXP_IN_MIN=
REFRESH_TOKEN_SECRET=
REFRESH_TOKEN_EXP_IN_DAYS=

# Hash
HASHING_SECRET=

# Email Service
EMAIL=
PASSWORD=
```

### Run the Project

1. **Run the API with docker**

```bash
make docker-run # for production
make docker-down # for production

make docker-dev # for development
make docker-dev-down # for development
```

2. **Run Migrations**

```bash
# You only need to run migrations once
make migrate-up
```

---

## 📦 API Features

- User Registration (Local + OAuth)
- Email Verification with OTP
- Login / Logout
- Password Reset
- Products and Categories
- Cart Management
- Orders and Checkout
- Admin Routes

---

## 🛠 Project Structure

```
afrad-api/
├── cmd
│   └── api
├── config
│   └── config.go
├── docker-compose.dev.yml
├── docker-compose.yml
├── Dockerfile
├── Dockerfile.dev
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── auth
│   ├── database
│   ├── middleware
│   ├── models
│   ├── s3
│   ├── server
│   └── utils
├── Makefile
├── project-docs
│   ├── endpoints.md
│   ├── reposConventions.md
│   ├── schema.sql
│   └── TODO.md
└── README.md
```

---

## 📄 License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

This project is licensed under the [MIT License](LICENSE) © 2025 Afrad Team.
