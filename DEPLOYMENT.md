# Backend Deployment Guide

Place Pro Platform - Go Backend API

## Overview

This is a Go backend application featuring:
- **Language**: Go 1.24
- **Framework**: Gorilla Mux for routing
- **Database**: PostgreSQL 15
- **Architecture**: Clean Architecture (Entity, UseCase, Repository, Handler)
- **Containerization**: Docker & Docker Compose
- **Features**: User management, Company management, Event management

## Prerequisites

- **Go**: Version 1.24 or higher
- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher
- **PostgreSQL**: Version 15+ (if running without Docker)

## Architecture

```
backend/
├── companyd/           # Company domain
│   ├── entity/         # Company entities
│   ├── handler/        # HTTP handlers
│   ├── presenter/      # Request/Response models
│   ├── repository/     # Database layer
│   └── usecase/        # Business logic
├── userd/              # User domain
│   ├── entity/         # User entities
│   ├── handler/        # HTTP handlers
│   ├── presenter/      # Request/Response models
│   ├── repository/     # Database layer
│   └── usecase/        # Business logic
├── docker-compose.yml  # Docker configuration
├── Dockerfile.golang   # Go application container
├── init.sql           # Database initialization
└── main.go            # Application entry point
```

## Environment Configuration

### Environment Variables

The application supports the following environment variables:

```bash
# Database Configuration
DB_HOST=localhost          # Database host
DB_PORT=5432              # Database port
DB_USER=myuser            # Database username
DB_PASSWORD=mypassword    # Database password
DB_NAME=myapp             # Database name

# Server Configuration
PORT=8080                 # Server port
```

## Local Development

### Option 1: Docker Compose (Recommended)

1. **Start all services**:
```bash
docker-compose up -d
```

2. **View logs**:
```bash
docker-compose logs -f app
```

3. **Stop services**:
```bash
docker-compose down
```

4. **Rebuild after code changes**:
```bash
docker-compose build app
docker-compose up -d app
```

### Option 2: Local Go Development

1. **Install dependencies**:
```bash
go mod download
```

2. **Start PostgreSQL** (using Docker):
```bash
docker run -d \
  --name postgres_local \
  -e POSTGRES_DB=myapp \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -p 5432:5432 \
  postgres:15
```

3. **Initialize database**:
```bash
psql -h localhost -U myuser -d myapp -f init.sql
```

4. **Run the application**:
```bash
go run main.go
```

### Development Commands

```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build binary
go build -o main .

# Format code
go fmt ./...

# Run linter
golangci-lint run

# Generate documentation
godoc -http=:6060
```

## Production Deployment

### Option 1: Docker Deployment (Recommended)

#### Single Server Deployment

1. **Prepare production environment**:
```bash
# Clone repository
git clone <repository-url>
cd back-end

# Create production environment file
cp .env.example .env
```

2. **Configure environment variables**:
```bash
# Edit .env file
DB_HOST=postgres
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_secure_password
DB_NAME=your_db_name
PORT=8080
```

3. **Deploy with Docker Compose**:
```bash
# Build and start services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

#### Production Docker Compose

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: postgres_prod
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build:
      context: .
      dockerfile: Dockerfile.golang
    image: placement-portal:latest
    container_name: golang_prod
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - PORT=8080
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/user/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

### Option 2: Kubernetes Deployment

#### Create Kubernetes Manifests

**namespace.yaml**:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: place-pro
```

**configmap.yaml**:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: place-pro
data:
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_NAME: "myapp"
  PORT: "8080"
```

**secret.yaml**:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: place-pro
type: Opaque
data:
  DB_USER: <base64-encoded-username>
  DB_PASSWORD: <base64-encoded-password>
```

**postgres-deployment.yaml**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: place-pro
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: DB_NAME
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_USER
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_PASSWORD
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
```

**app-deployment.yaml**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: place-pro-app
  namespace: place-pro
spec:
  replicas: 3
  selector:
    matchLabels:
      app: place-pro-app
  template:
    metadata:
      labels:
        app: place-pro-app
    spec:
      containers:
      - name: app
        image: placement-portal:latest
        envFrom:
        - configMapRef:
            name: app-config
        - secretRef:
            name: app-secrets
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /user/health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /user/health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### Deploy to Kubernetes

```bash
# Apply manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n place-pro

# View logs
kubectl logs -f deployment/place-pro-app -n place-pro
```

### Option 3: Cloud Platform Deployment

#### AWS ECS Deployment

1. **Create ECR repository**:
```bash
aws ecr create-repository --repository-name placement-portal
```

2. **Build and push image**:
```bash
# Build image
docker build -t placement-portal .

# Tag for ECR
docker tag placement-portal:latest <account-id>.dkr.ecr.<region>.amazonaws.com/placement-portal:latest

# Push to ECR
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/placement-portal:latest
```

3. **Create ECS task definition** and service using AWS Console or CLI.

#### Google Cloud Run Deployment

1. **Build and push to Container Registry**:
```bash
# Build image
gcloud builds submit --tag gcr.io/PROJECT-ID/placement-portal

# Deploy to Cloud Run
gcloud run deploy placement-portal \
  --image gcr.io/PROJECT-ID/placement-portal \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

#### Azure Container Instances

```bash
# Create resource group
az group create --name myResourceGroup --location eastus

# Create container
az container create \
  --resource-group myResourceGroup \
  --name placement-portal \
  --image placement-portal:latest \
  --ports 8080 \
  --environment-variables \
    DB_HOST=your-db-host \
    DB_USER=your-db-user \
    DB_NAME=your-db-name
```

### Option 4: Traditional Server Deployment

#### Build and Deploy Binary

1. **Build for production**:
```bash
# Build for Linux (if building on different OS)
GOOS=linux GOARCH=amd64 go build -o placement-portal main.go

# Or build for current OS
go build -o placement-portal main.go
```

2. **Create systemd service** (`/etc/systemd/system/placement-portal.service`):
```ini
[Unit]
Description=Place Pro Platform API
After=network.target

[Service]
Type=simple
User=app
WorkingDirectory=/opt/placement-portal
ExecStart=/opt/placement-portal/placement-portal
Restart=on-failure
RestartSec=5
Environment=DB_HOST=localhost
Environment=DB_PORT=5432
Environment=DB_USER=myuser
Environment=DB_PASSWORD=mypassword
Environment=DB_NAME=myapp
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

3. **Enable and start service**:
```bash
sudo systemctl enable placement-portal
sudo systemctl start placement-portal
sudo systemctl status placement-portal
```

## Database Management

### Database Migrations

The application uses `init.sql` for initial schema. For production:

1. **Create migration system**:
```bash
# Install migrate tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq initial_schema
```

2. **Run migrations**:
```bash
migrate -path migrations -database "postgres://user:pass@host:port/dbname?sslmode=disable" up
```

### Database Backup

```bash
# Backup database
docker exec postgres_db pg_dump -U myuser myapp > backup.sql

# Restore database
docker exec -i postgres_db psql -U myuser myapp < backup.sql
```

## Monitoring and Logging

### Health Checks

The application provides health check endpoints:

- `GET /user/health` - User service health
- `GET /company/health` - Company service health
- `GET /user/dbtest` - Database connectivity test

### Logging

Add structured logging:

```go
import "log/slog"

// Configure structured logging
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(logger)

// Use in handlers
slog.Info("Request received", "method", r.Method, "path", r.URL.Path)
```

### Metrics

Add Prometheus metrics:

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Add metrics endpoint
http.Handle("/metrics", promhttp.Handler())
```

## Security Considerations

### 1. Environment Variables

- Never commit sensitive data to version control
- Use secrets management (AWS Secrets Manager, HashiCorp Vault)
- Rotate passwords regularly

### 2. Database Security

```sql
-- Create application user with limited permissions
CREATE USER app_user WITH ENCRYPTED PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE myapp TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
```

### 3. CORS Configuration

Update CORS settings for production:

```go
allowedOrigins := []string{
    "https://yourdomain.com",
    "https://www.yourdomain.com",
}
```

### 4. HTTPS Configuration

Use reverse proxy (nginx) for HTTPS:

```nginx
server {
    listen 443 ssl;
    server_name api.yourdomain.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Performance Optimization

### 1. Database Optimization

```sql
-- Add indexes for better performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_companies_name ON companies(company_name);
CREATE INDEX idx_events_date ON events(date);
```

### 2. Connection Pooling

```go
// Configure database connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 3. Caching

Add Redis for caching:

```go
import "github.com/go-redis/redis/v8"

// Configure Redis client
rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
```

## Troubleshooting

### Common Issues

1. **Database Connection Issues**:
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Test connection
docker exec -it postgres_db psql -U myuser -d myapp
```

2. **Application Won't Start**:
```bash
# Check logs
docker-compose logs app

# Check port availability
lsof -i :8080

# Rebuild container
docker-compose build --no-cache app
```

3. **CORS Issues**:
   - Verify allowed origins in handler
   - Check preflight OPTIONS handling
   - Ensure proper headers are set

### Debug Commands

```bash
# Check container status
docker-compose ps

# View application logs
docker-compose logs -f app

# Access container shell
docker-compose exec app sh

# Check database
docker-compose exec postgres psql -U myuser -d myapp

# Monitor resource usage
docker stats
```

## Maintenance

### Regular Tasks

```bash
# Update dependencies
go get -u ./...
go mod tidy

# Security scan
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Build and test
go build ./...
go test ./...

# Database maintenance
docker-compose exec postgres psql -U myuser -d myapp -c "VACUUM ANALYZE;"
```

### Backup Strategy

```bash
# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec postgres_db pg_dump -U myuser myapp > "backup_${DATE}.sql"
gzip "backup_${DATE}.sql"

# Keep only last 7 days of backups
find . -name "backup_*.sql.gz" -mtime +7 -delete
```

---

## Quick Deployment Checklist

- [ ] Configure environment variables
- [ ] Set up database (PostgreSQL)
- [ ] Build application image
- [ ] Configure CORS for production domains
- [ ] Set up reverse proxy (nginx) for HTTPS
- [ ] Configure monitoring and logging
- [ ] Set up database backups
- [ ] Test all API endpoints
- [ ] Configure health checks
- [ ] Set up SSL certificates

## API Endpoints

### User Management
- `POST /user/login` - User authentication
- `GET /user/list` - List all users
- `POST /user/create` - Create new user
- `DELETE /user/delete/{id}` - Delete user
- `GET /user/health` - Health check

### Company Management
- `GET /company/list` - List all companies
- `POST /company/create` - Create new company
- `PUT /company/update/{id}` - Update company
- `DELETE /company/delete/{id}` - Delete company
- `GET /company/list/{username}` - List companies by officer

### Event Management
- `GET /event/list` - List all events
- `POST /event/create` - Create new event
- `PUT /event/update/{id}` - Update event
- `DELETE /event/delete/{id}` - Delete event

For support, check the main project documentation or contact the development team. 