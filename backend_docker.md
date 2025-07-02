# How to Run the Golang + PostgreSQL Docker Setup

## Prerequisites

Make sure you have the following installed on your system:
- **Docker** (version 20.10 or later)
- **Docker Compose** (version 2.0 or later)

Check your versions:
```bash
docker --version
docker-compose --version
```

## Step-by-Step Setup

### Step 1: Create Project Directory
```bash
mkdir golang-postgres-app
cd golang-postgres-app
```

### Step 2: Create Required Files

Create the following directory structure:
```
golang-postgres-app/
├── Dockerfile.golang
├── Dockerfile.postgres
├── docker-compose.yml
├── main.go
├── go.mod
├── init.sql
└── init-scripts/
    └── 01-init.sql
```

**Create the init-scripts directory:**
```bash
mkdir init-scripts
```

### Step 3: Copy All the Code Files

Copy each file content from the artifacts above:

1. **Dockerfile.golang** - Copy the Golang Dockerfile content
2. **Dockerfile.postgres** - Copy the PostgreSQL Dockerfile content  
3. **docker-compose.yml** - Copy the Docker Compose configuration
4. **main.go** - Copy the sample Golang application
5. **go.mod** - Copy the Go module file
6. **init.sql** - Copy the application init.sql file
7. **init-scripts/01-init.sql** - Copy the database initialization script

### Step 4: Build and Run

**Option 1: Build and run in one command (Recommended)**
```bash
docker-compose up --build
```

**Option 2: Build first, then run**
```bash
# Build the images
docker-compose build

# Run the containers
docker-compose up
```

**Option 3: Run in background (detached mode)**
```bash
docker-compose up --build -d
```

## What You Should See

When running successfully, you'll see output similar to:
```
Creating network "golang-postgres-app_app-network" with driver "bridge"
Creating volume "golang-postgres-app_postgres_data" with default driver
Building postgres
Building app
Creating postgres_db ... done
Creating golang_app ... done
Attaching to postgres_db, golang_app
postgres_db | PostgreSQL init process complete; ready for start up.
golang_app  | Successfully connected to PostgreSQL!
golang_app  | Found and reading SQL file: ./init.sql
golang_app  | Successfully executed all SQL statements from ./init.sql
golang_app  | Server starting on :8080
```

## Step 5: Test the Application

Open a new terminal and test the API endpoints:

**1. Check if the app is running:**
```bash
curl http://localhost:8080/
```
Expected response:
```json
{"message":"Welcome to the Golang + PostgreSQL API!","status":"running"}
```

**2. Check health status:**
```bash
curl http://localhost:8080/health
```
Expected response:
```json
{"status":"healthy","database":"connected"}
```

**3. Get all users:**
```bash
curl http://localhost:8080/users
```
Expected response:
```json
[
  {"id":1,"name":"John Doe","email":"john@example.com"},
  {"id":2,"name":"Jane Smith","email":"jane@example.com"},
  {"id":3,"name":"Alice Johnson","email":"alice@example.com"}
]
```

**4. Create a new user:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"New User","email":"newuser@example.com"}'
```

## Step 6: View Logs

**View logs from both services:**
```bash
docker-compose logs
```

**View logs from specific service:**
```bash
docker-compose logs app
docker-compose logs postgres
```

**Follow logs in real-time:**
```bash
docker-compose logs -f
```

## Management Commands

**Stop the services:**
```bash
docker-compose down
```

**Stop and remove volumes (deletes database data):**
```bash
docker-compose down -v
```

**Restart services:**
```bash
docker-compose restart
```

**Rebuild and restart:**
```bash
docker-compose up --build --force-recreate
```

## Accessing the Database Directly

If you want to connect to PostgreSQL directly:
```bash
docker-compose exec postgres psql -U myuser -d myapp
```

Once connected, you can run SQL commands:
```sql
\dt                    -- List tables
SELECT * FROM users;   -- View users
SELECT * FROM logs;    -- View logs
\q                     -- Quit
```

## Troubleshooting

**If you get connection errors:**
1. Make sure both containers are running: `docker-compose ps`
2. Check logs: `docker-compose logs`
3. Restart: `docker-compose down && docker-compose up --build`

**If port 8080 is already in use:**
Change the port mapping in docker-compose.yml:
```yaml
ports:
  - "8081:8080"  # Use port 8081 instead
```

**If you get permission errors:**
Make sure Docker daemon is running and you have proper permissions:
```bash
sudo docker-compose up --build  # On Linux, if needed
```

## File Permissions

Make sure all files have proper read permissions:
```bash
chmod +r *.go *.sql *.yml Dockerfile.*
chmod +r init-scripts/*
```

That's it! Your Golang application with PostgreSQL should now be running successfully.