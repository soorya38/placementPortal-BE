# CORS Configuration

This document explains how to configure CORS (Cross-Origin Resource Sharing) allowed origins for the backend API.

## Environment Variable

The backend now supports configuring CORS allowed origins through the `CORS_ALLOWED_ORIGINS` environment variable.

### Setting up CORS_ALLOWED_ORIGINS

**Format:** Comma-separated list of allowed origins
**Environment Variable:** `CORS_ALLOWED_ORIGINS`

### Examples

#### Using Environment Variables

```bash
# Single origin
export CORS_ALLOWED_ORIGINS="https://yourdomain.com"

# Multiple origins
export CORS_ALLOWED_ORIGINS="https://yourdomain.com,https://www.yourdomain.com,http://localhost:3000"
```

#### Using Docker Compose

Update your `docker-compose.yml`:

```yaml
services:
  app:
    environment:
      - CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com,http://localhost:3000
```

#### Using .env file

Create a `.env` file in your project root:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=myapp

# Server Configuration
PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com,http://localhost:3000
```

### Default Values

If `CORS_ALLOWED_ORIGINS` is not set, the following default origins are used:

- `https://0f22-2402-3a80-1325-cd70-dd05-94a2-213-dd84.ngrok-free.app`
- `https://place-pro-platform-88.vercel.app`
- `http://localhost:8081`

### Production Deployment

For production environments, make sure to:

1. Set the `CORS_ALLOWED_ORIGINS` environment variable to only include your production domains
2. Remove any localhost or development URLs from the production configuration
3. Use HTTPS origins whenever possible

### Security Notes

- Only add trusted domains to the allowed origins list
- Avoid using wildcards (*) in production
- Keep the list as minimal as possible to reduce attack surface
- Regularly review and update the allowed origins list

### Troubleshooting

If you're experiencing CORS issues:

1. Check that your frontend domain is included in `CORS_ALLOWED_ORIGINS`
2. Ensure there are no typos in the origin URLs
3. Verify that the protocol (http/https) matches exactly
4. Check the server logs for CORS-related messages

### Testing CORS Configuration

You can test CORS configuration by checking the response headers:

```bash
curl -H "Origin: https://yourdomain.com" \
     -H "Access-Control-Request-Method: GET" \
     -H "Access-Control-Request-Headers: X-Requested-With" \
     -X OPTIONS \
     http://localhost:8080/user/health
```

The response should include:
- `Access-Control-Allow-Origin: https://yourdomain.com`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With, ngrok-skip-browser-warning` 