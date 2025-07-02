# Environment Variables Configuration

This file documents all environment variables used in the application. Use this as a template for setting up your environment.

## Required Environment Variables

### Database Configuration

```bash
DB_HOST=localhost          # Database host (default: localhost)
DB_PORT=5432              # Database port (default: 5432)
DB_USER=myuser            # Database username (default: myuser)
DB_PASSWORD=mypassword    # Database password (default: mypassword)
DB_NAME=myapp             # Database name (default: myapp)
```

### Server Configuration

```bash
PORT=8080                 # Server port (default: 8080)
```

### CORS Configuration

```bash
# Comma-separated list of allowed origins for CORS requests
CORS_ALLOWED_ORIGINS=https://0f22-2402-3a80-1325-cd70-dd05-94a2-213-dd84.ngrok-free.app,https://place-pro-platform-88.vercel.app,https://localhost:8081,http://localhost:8081
```

## Environment-Specific Examples

### Local Development

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=myapp
PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://localhost:8081,http://localhost:8081,http://127.0.0.1:3000
```

### Production

```bash
DB_HOST=your-production-db-host
DB_PORT=5432
DB_USER=your-production-user
DB_PASSWORD=your-secure-production-password
DB_NAME=your-production-db
PORT=8080
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### Staging

```bash
DB_HOST=your-staging-db-host
DB_PORT=5432
DB_USER=your-staging-user
DB_PASSWORD=your-staging-password
DB_NAME=your-staging-db
PORT=8080
CORS_ALLOWED_ORIGINS=https://staging.yourdomain.com,https://dev.yourdomain.com
```

## How to Set Environment Variables

### Option 1: Using .env file (Recommended for local development)

1. Create a `.env` file in the project root
2. Add the variables in `KEY=VALUE` format:

```bash
# .env file
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=myapp
PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://localhost:8081,http://localhost:8081
```

### Option 2: Export environment variables

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=myuser
export DB_PASSWORD=mypassword
export DB_NAME=myapp
export PORT=8080
export CORS_ALLOWED_ORIGINS=http://localhost:3000,https://localhost:8081,http://localhost:8081
```

### Option 3: Docker Compose (Already configured)

The `docker-compose.yml` file already includes these environment variables. You can modify them directly in the file or override them with a `.env` file.

## Security Notes

- **Never commit sensitive values** like passwords to version control
- Use strong, unique passwords for production databases
- Limit CORS origins to trusted domains only
- Consider using secrets management systems for production deployments
- Regularly rotate database passwords

## Default Values

If environment variables are not set, the application uses these defaults:

- `DB_HOST`: localhost
- `DB_PORT`: 5432
- `DB_USER`: myuser
- `DB_PASSWORD`: mypassword
- `DB_NAME`: myapp
- `PORT`: 8080
- `CORS_ALLOWED_ORIGINS`: Uses hardcoded defaults (ngrok, vercel, https://localhost:8081, http://localhost:8081) 