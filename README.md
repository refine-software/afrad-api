# afrad-api

## Get Up and Running

## ENV File

**NOTE:** These environment variables are for development, so it's safe to upload them.

```txt
# APP
PORT=8080
APP_ENV=dev # use 'dev' for development and 'prod' for production

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
```
