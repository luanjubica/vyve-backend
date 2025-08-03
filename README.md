# Vyve Backend

A comprehensive Go backend API for the Vyve relationship management application, built with Gin, GORM, Redis, and Supabase authentication.

## 🚀 Features

- **🔐 Authentication**: Supabase JWT validation
- **💾 Database**: PostgreSQL with GORM ORM and normalized schema
- **⚡ Caching**: Redis for performance optimization  
- **🤖 Background Jobs**: Cron-based pattern detection and insights
- **📊 Analytics**: Relationship health scoring and energy analysis
- **🔔 Smart Nudges**: AI-powered recommendations and reminders
- **🏗️ Clean Architecture**: Well-structured, maintainable codebase

## 📋 Prerequisites

- Go 1.21 or later
- Docker & Docker Compose (recommended)
- Supabase project with database

## 🛠️ Quick Setup

### 1. Create the project
```bash
# Run this script to create the complete project structure
bash create_vyve_backend.sh
cd vyve-backend
```

### 2. Configure environment
```bash
# Copy and update environment variables
cp .env.example .env
# Edit .env with your Supabase credentials
```

### 3. Set up the project
```bash
make setup
```

### 4. Run the application

**Option A: With Docker (Recommended)**
```bash
make docker-run
```

**Option B: Local Development**
```bash
make dev  # With hot reload
# or
make run  # Single run
```

### 5. Verify it's working
```bash
curl http://localhost:8080/health
```

## 📚 API Documentation

### Core Endpoints

#### Authentication
```
POST /auth/verify     # Verify Supabase JWT token
GET  /auth/profile    # Get user profile
PUT  /auth/profile    # Update user profile
```

#### Lookups (Reference Data)
```
GET /api/v1/lookups   # Get all lookup data for frontend
```

#### People Management
```
GET    /api/v1/people           # Get all user's people
POST   /api/v1/people           # Add new person
GET    /api/v1/people/:id       # Get person details
PUT    /api/v1/people/:id       # Update person
DELETE /api/v1/people/:id       # Delete person
GET    /api/v1/people/:id/stats # Get person statistics
```

#### Interactions
```
GET    /api/v1/interactions     # Get interactions (with filters)
POST   /api/v1/interactions     # Log new interaction
GET    /api/v1/interactions/:id # Get interaction details
PUT    /api/v1/interactions/:id # Update interaction
DELETE /api/v1/interactions/:id # Delete interaction
```

#### Analytics & Insights
```
GET /api/v1/analytics/dashboard     # Dashboard metrics
GET /api/v1/analytics/patterns      # Relationship patterns
GET /api/v1/analytics/health-scores # Health score trends
GET /api/v1/analytics/energy        # Energy analysis
```

#### Nudges & Notifications
```
GET    /api/v1/nudges              # Get user nudges
POST   /api/v1/nudges/:id/read     # Mark nudge as read
DELETE /api/v1/nudges/:id          # Dismiss nudge
GET    /api/v1/nudges/generate     # Trigger nudge generation
```

#### Daily Reflections
```
GET /api/v1/reflections           # Get user reflections
POST /api/v1/reflections          # Save new reflection
GET /api/v1/reflections/:id       # Get reflection details
PUT /api/v1/reflections/:id       # Update reflection
```

## 🏗️ Project Structure

```
vyve-backend/
├── cmd/server/         # Application entry point
├── internal/
│   ├── handlers/       # HTTP handlers (7 files)
│   ├── services/       # Business logic (5 files)
│   ├── models/         # Database models (8 files)
│   ├── middleware/     # HTTP middleware (4 files)
│   ├── auth/          # Authentication logic
│   ├── database/      # Database connection & migrations
│   └── jobs/          # Background jobs (3 files)
├── pkg/
│   ├── utils/         # Utility functions (2 files)
│   └── validator/     # Validation helpers
├── configs/           # Configuration (2 files)
├── scripts/           # Setup and utility scripts
├── deployments/       # Docker and Kubernetes configs
├── docs/             # Documentation (3 files)
├── tests/            # Integration and unit tests (5 files)
└── [config files]    # Docker, Make, Git, etc.
```

## 🛠️ Development Commands

```bash
# Development
make dev          # Run with hot reload
make run          # Run once
make build        # Build binary
make test         # Run tests
make format       # Format code
make lint         # Lint code

# Dependencies
make deps         # Download dependencies
make setup        # Initial setup

# Docker
make docker-build # Build Docker image
make docker-run   # Run with Docker Compose
make docker-stop  # Stop Docker containers

# Utilities
make clean        # Clean build files
make check        # Run all checks (format, lint, test)
```

## 🔧 Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `ENVIRONMENT` | Environment (development/production) | No |
| `PORT` | Server port | No |
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `SUPABASE_URL` | Supabase project URL | Yes |
| `SUPABASE_ANON_KEY` | Supabase anonymous key | Yes |
| `SUPABASE_JWT_SECRET` | Supabase JWT secret | Yes |
| `REDIS_URL` | Redis connection string | No |
| `LOG_LEVEL` | Logging level | No |

## 🤖 Background Jobs

The application runs several automated background jobs:

- **Pattern Detection** (hourly): Identifies concerning relationship patterns
- **Reconnection Reminders** (daily 9 AM): Suggests people to reconnect with  
- **Energy Insights** (daily 6 PM): Finds energy-giving relationships
- **Health Score Updates** (every 6 hours): Recalculates relationship health
- **Cleanup** (daily midnight): Removes old nudges and expired data

## 🔐 Security Features

- JWT token validation for all API endpoints
- Row Level Security (RLS) in database
- Rate limiting middleware
- Input validation on all endpoints
- SQL injection protection via GORM
- CORS configuration for frontend integration

## 📊 Database Schema

The application uses a normalized PostgreSQL schema with:

- **9 lookup tables** for reference data (categories, methods, statuses, etc.)
- **8 main tables** for user data (people, interactions, reflections, etc.)
- **Foreign key relationships** ensuring data integrity
- **Indexes** for optimal query performance
- **Triggers** for automatic health score calculation

## 🚀 Deployment

### Docker Deployment
```bash
# Production deployment
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Cloud Deployment
1. Set environment variables in your cloud platform
2. Deploy from Git repository  
3. Ensure Supabase database is accessible
4. Configure Redis if using external Redis

## 🧪 Testing

```bash
make test           # Run all tests
make test-coverage  # Run with coverage report
```

## 📄 License

MIT License

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `make check` to verify code quality
5. Submit a pull request

---

Built with ❤️ for intentional relationships
