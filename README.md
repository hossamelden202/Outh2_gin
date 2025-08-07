# OAuth2 Authentication Microservice

A secure OAuth2 authentication service built with Go and Gin framework. Supports multiple OAuth providers (Google, GitHub, LinkedIn) with JWT-based session management and device tracking.

## Features

- **Multi-Provider OAuth2**: Google, GitHub, and LinkedIn authentication
- **JWT Token Management**: Secure access and refresh tokens with version control
- **Session Tracking**: Device-aware session management with Redis
- **Token Rotation**: Automatic token refresh mechanism
- **Logout Management**: Single and bulk logout functionality
- **Email Verification**: Email verification on account creation
- **Device Information**: Track user login locations and device details

## Tech Stack

- **Framework**: Gin Web Framework
- **Authentication**: JWT (golang-jwt)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Language**: Go 1.23

## Prerequisites

- Go 1.23 or higher
- PostgreSQL
- Redis
- Environment variables configured (see `.env.example`)

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd outh2

# Install dependencies
go mod download

# Configure environment variables
cp .env.example .env
# Edit .env with your credentials

# Run the service
go run main.go
```

## Environment Variables

```
# Server
PORT=8080

# Database
HOST=localhost
PORT=5432
user_name=postgres
password=your_password
db_name=oauth2_db
ssl=disable

# Redis
REDIS_URL=localhost:6379

# JWT
jwt_secret=your_secret_key

# Google OAuth
CLIENT_ID=your_google_client_id
CLIENT_SECRET=your_google_client_secret
RED_URL=http://localhost:8080/oauth2/callback/google
GOOGLE_STATE=state_value

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_RED_URL=http://localhost:8080/oauth2/callback/github
GITHUB_STATE=state_value

# LinkedIn OAuth
LINKDIN_CLIENT_ID=your_linkedin_client_id
LINKDIN_CLIENT_SECRET=your_linkedin_client_secret
LINKDIN_RED_URL=http://localhost:8080/oauth2/callback/linkedin
LINKDIN_STATE=state_value
```

## Database Schema

### Users Table
Stores user account information and authentication details.

### Device Record Table
Tracks user login locations and device information.

### Sessions Table
Maintains active user sessions with device mapping.

## API Endpoints

### Authorization
- `GET /oauth2/authorize/:provider` - Initiate OAuth flow

### OAuth Callbacks
- `GET /oauth2/callback/google` - Google OAuth callback
- `GET /oauth2/callback/github` - GitHub OAuth callback
- `GET /oauth2/callback/linkedin` - LinkedIn OAuth callback

### Session Management
- `POST /oauth2/refresh` - Refresh access token
- `POST /oauth2/logout` - Logout from current device
- `POST /oauth2/logoutall` - Logout from all devices
- `GET /oauth2/get-sessions` - Retrieve active sessions (requires auth)

## Security Features

### Token Management
- JWT tokens expire after 15 minutes
- Refresh tokens valid for 30 days
- Token versioning prevents token hijacking
- JTI (JWT ID) prevents token reuse

### Token Validation
- Signature verification using HMAC
- Expiration time validation
- Token version checking against database
- Blocklist management for revoked tokens

### Session Security
- Device fingerprinting and tracking
- IP geolocation recording
- Browser identification
- Session invalidation on logout

### Authorization
- Bearer token authentication via middleware
- Email verification for new accounts
- Role-based access control
- Device-aware authentication

## Running the Service

```bash
# Development mode
go run main.go

# Build and run
go build
./outh2
```

The service runs on `http://localhost:8080`

## API Documentation

See `swagger.json` for OpenAPI 3.0 specification and `security.html` for detailed security documentation.
