# Place Pro Platform - Backend API

A comprehensive Go-based backend API for placement management featuring user authentication, company management, and event scheduling.

## 🚀 Features

- **User Management**: Authentication, role-based access, user CRUD operations
- **Company Management**: Company profiles, application tracking, officer assignments
- **Event Management**: Calendar events, notifications, scheduling
- **Clean Architecture**: Domain-driven design with separation of concerns
- **RESTful API**: Well-structured endpoints with proper HTTP methods
- **Database Integration**: PostgreSQL with connection pooling
- **Docker Support**: Containerized deployment with Docker Compose
- **CORS Configuration**: Flexible cross-origin resource sharing setup
- **Environment-based Configuration**: Easy deployment across environments

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Production Deployment](#production-deployment)
- [Environment Configuration](#environment-configuration)
- [Database Setup](#database-setup)
- [Security Configuration](#security-configuration)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [API Documentation](#api-documentation)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## 🛠 Prerequisites

### For Development
- **Go**: Version 1.24 or higher
- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher
- **PostgreSQL**: Version 15+ (if running without Docker)

### For Production
- **Linux Server**: Ubuntu 20.04+ or CentOS 8+ recommended
- **Docker & Docker Compose**: Latest stable versions
- **Domain Name**: For HTTPS setup
- **SSL Certificate**: Let's Encrypt or commercial certificate
- **Reverse Proxy**: Nginx or similar (recommended)
- **Minimum Resources**: 2GB RAM, 20GB disk space

## 🚀 Quick Start

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd back-end
   ```

2. **Start with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Verify installation**
   ```bash
   curl http://localhost:8080/user/health
   curl http://localhost:8080/company/health
   ```

The API will be available at `http://localhost:8080`

## 🏭 Production Deployment

### Step 1: Server Preparation

1. **Update system packages**
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

2. **Install Docker**
   ```bash
   # Remove old versions
   sudo apt remove docker docker-engine docker.io containerd runc

   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   sudo usermod -aG docker $USER
   ```

3. **Install Docker Compose**
   ```bash
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   ```

4. **Configure firewall**
   ```bash
   sudo ufw allow ssh
   sudo ufw allow 80
   sudo ufw allow 443
   sudo ufw enable
   ```

### Step 2: Application Setup

1. **Create application directory**
   ```bash
   sudo mkdir -p /opt/place-pro
   sudo chown $USER:$USER /opt/place-pro
   cd /opt/place-pro
   ```

2. **Clone repository**
   ```bash
   git clone <repository-url> .
   ```

3. **Create production environment file**
   ```bash
   cp ENVIRONMENT_VARIABLES.md .env
   # Edit .env with your production values
   nano .env
   ```

4. **Set production environment variables**
   ```bash
   # .env file
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=your_production_user
   DB_PASSWORD=your_secure_password
   DB_NAME=your_production_db
   PORT=8080
   CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
   ```

### Step 3: Database Configuration

1. **Create production Docker Compose override**
   ```bash
   cat > docker-compose.prod.yml << 'EOF'
   version: '3.8'

   services:
     postgres:
       environment:
         POSTGRES_DB: ${DB_NAME}
         POSTGRES_USER: ${DB_USER}
         POSTGRES_PASSWORD: ${DB_PASSWORD}
       volumes:
         - postgres_data:/var/lib/postgresql/data
         - ./backups:/backups
       restart: unless-stopped
       healthcheck:
         test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
         interval: 30s
         timeout: 10s
         retries: 5

     app:
       environment:
         - DB_HOST=postgres
         - DB_PORT=5432
         - DB_USER=${DB_USER}
         - DB_PASSWORD=${DB_PASSWORD}
         - DB_NAME=${DB_NAME}
         - PORT=8080
         - CORS_ALLOWED_ORIGINS=${CORS_ALLOWED_ORIGINS}
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
   EOF
   ```

2. **Deploy the application**
   ```bash
   # Load environment variables
   source .env

   # Start services
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

   # Check status
   docker-compose ps
   ```

### Step 4: Reverse Proxy Setup (Nginx)

1. **Install Nginx**
   ```bash
   sudo apt install nginx -y
   ```

2. **Create Nginx configuration**
   ```bash
   sudo tee /etc/nginx/sites-available/place-pro << 'EOF'
   upstream place_pro_backend {
       server localhost:8080;
   }

   server {
       listen 80;
       server_name yourdomain.com www.yourdomain.com;
       
       # Redirect HTTP to HTTPS
       return 301 https://$server_name$request_uri;
   }

   server {
       listen 443 ssl http2;
       server_name yourdomain.com www.yourdomain.com;

       # SSL Configuration
       ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
       ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
       ssl_protocols TLSv1.2 TLSv1.3;
       ssl_ciphers HIGH:!aNULL:!MD5;

       # Security Headers
       add_header X-Content-Type-Options nosniff;
       add_header X-Frame-Options DENY;
       add_header X-XSS-Protection "1; mode=block";
       add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

       # API Proxy
       location / {
           proxy_pass http://place_pro_backend;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
           
           # CORS headers (if needed)
           proxy_hide_header Access-Control-Allow-Origin;
           add_header Access-Control-Allow-Origin $http_origin always;
           
           # Timeouts
           proxy_connect_timeout 30s;
           proxy_send_timeout 30s;
           proxy_read_timeout 30s;
       }

       # Health check endpoint
       location /health {
           access_log off;
           proxy_pass http://place_pro_backend/user/health;
       }
   }
   EOF
   ```

3. **Enable site and test configuration**
   ```bash
   sudo ln -s /etc/nginx/sites-available/place-pro /etc/nginx/sites-enabled/
   sudo nginx -t
   ```

### Step 5: SSL Certificate Setup

1. **Install Certbot**
   ```bash
   sudo apt install certbot python3-certbot-nginx -y
   ```

2. **Obtain SSL certificate**
   ```bash
   sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com
   ```

3. **Test auto-renewal**
   ```bash
   sudo certbot renew --dry-run
   ```

4. **Start Nginx**
   ```bash
   sudo systemctl enable nginx
   sudo systemctl start nginx
   ```

### Step 6: Monitoring Setup

1. **Create monitoring script**
   ```bash
   cat > /opt/place-pro/monitor.sh << 'EOF'
   #!/bin/bash
   
   # Health check script
   LOG_FILE="/var/log/place-pro-monitor.log"
   
   check_service() {
       local service_name=$1
       local health_url=$2
       
       if curl -f -s $health_url > /dev/null; then
           echo "$(date): $service_name - OK" >> $LOG_FILE
           return 0
       else
           echo "$(date): $service_name - FAILED" >> $LOG_FILE
           return 1
       fi
   }
   
   # Check services
   check_service "User Service" "http://localhost:8080/user/health"
   check_service "Company Service" "http://localhost:8080/company/health"
   
   # Check database
   if docker exec postgres_db pg_isready -U ${DB_USER} -d ${DB_NAME} > /dev/null; then
       echo "$(date): Database - OK" >> $LOG_FILE
   else
       echo "$(date): Database - FAILED" >> $LOG_FILE
   fi
   EOF
   
   chmod +x /opt/place-pro/monitor.sh
   ```

2. **Setup cron job for monitoring**
   ```bash
   (crontab -l 2>/dev/null; echo "*/5 * * * * /opt/place-pro/monitor.sh") | crontab -
   ```

### Step 7: Backup Configuration

1. **Create backup script**
   ```bash
   cat > /opt/place-pro/backup.sh << 'EOF'
   #!/bin/bash
   
   BACKUP_DIR="/opt/place-pro/backups"
   DATE=$(date +%Y%m%d_%H%M%S)
   
   # Create backup directory
   mkdir -p $BACKUP_DIR
   
   # Database backup
   docker exec postgres_db pg_dump -U ${DB_USER} ${DB_NAME} > $BACKUP_DIR/db_backup_$DATE.sql
   
   # Compress backup
   gzip $BACKUP_DIR/db_backup_$DATE.sql
   
   # Remove backups older than 7 days
   find $BACKUP_DIR -name "db_backup_*.sql.gz" -mtime +7 -delete
   
   echo "Backup completed: db_backup_$DATE.sql.gz"
   EOF
   
   chmod +x /opt/place-pro/backup.sh
   ```

2. **Schedule daily backups**
   ```bash
   (crontab -l 2>/dev/null; echo "0 2 * * * /opt/place-pro/backup.sh") | crontab -
   ```

## ⚙️ Environment Configuration

### Required Environment Variables

| Variable | Description | Default | Production Example |
|----------|-------------|---------|-------------------|
| `DB_HOST` | Database host | localhost | postgres |
| `DB_PORT` | Database port | 5432 | 5432 |
| `DB_USER` | Database username | myuser | prod_user |
| `DB_PASSWORD` | Database password | mypassword | secure_password_123 |
| `DB_NAME` | Database name | myapp | place_pro_db |
| `PORT` | Server port | 8080 | 8080 |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | localhost origins | https://yourdomain.com |

### Environment Files

Create different environment files for different stages:

**Production (.env.prod)**
```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=super_secure_password_123
DB_NAME=place_pro_production
PORT=8080
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

**Staging (.env.staging)**
```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=staging_user
DB_PASSWORD=staging_password_123
DB_NAME=place_pro_staging
PORT=8080
CORS_ALLOWED_ORIGINS=https://staging.yourdomain.com
```

## 🔒 Security Configuration

### 1. Database Security

```sql
-- Create production user with limited permissions
CREATE USER prod_user WITH ENCRYPTED PASSWORD 'super_secure_password_123';
CREATE DATABASE place_pro_production OWNER prod_user;
GRANT CONNECT ON DATABASE place_pro_production TO prod_user;
GRANT USAGE ON SCHEMA public TO prod_user;
GRANT CREATE ON SCHEMA public TO prod_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO prod_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO prod_user;
```

### 2. Application Security

- **Environment Variables**: Never commit sensitive data to version control
- **Password Policies**: Implement strong password requirements
- **Rate Limiting**: Consider adding rate limiting middleware
- **Input Validation**: Validate all user inputs
- **HTTPS Only**: Force HTTPS in production
- **Security Headers**: Implemented in Nginx configuration

### 3. Network Security

```bash
# Configure firewall
sudo ufw deny incoming
sudo ufw allow outgoing
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable

# Disable root login (optional but recommended)
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart ssh
```

## 📊 Monitoring & Maintenance

### Health Check Endpoints

- `GET /user/health` - User service health
- `GET /company/health` - Company service health
- `GET /user/dbtest` - Database connectivity test

### Log Management

1. **Application Logs**
   ```bash
   # View application logs
   docker-compose logs -f app
   
   # View database logs
   docker-compose logs -f postgres
   ```

2. **Nginx Logs**
   ```bash
   # Access logs
   sudo tail -f /var/log/nginx/access.log
   
   # Error logs
   sudo tail -f /var/log/nginx/error.log
   ```

3. **System Monitoring**
   ```bash
   # Monitor resource usage
   docker stats
   
   # Monitor disk usage
   df -h
   
   # Monitor memory usage
   free -m
   ```

### Performance Optimization

1. **Database Optimization**
   ```sql
   -- Add indexes for better performance
   CREATE INDEX idx_users_username ON users(username);
   CREATE INDEX idx_companies_name ON companies(company_name);
   CREATE INDEX idx_events_date ON events(date);
   
   -- Configure PostgreSQL
   ALTER SYSTEM SET max_connections = 200;
   ALTER SYSTEM SET shared_buffers = '256MB';
   ALTER SYSTEM SET effective_cache_size = '1GB';
   SELECT pg_reload_conf();
   ```

2. **Connection Pooling** (already configured in the application)

3. **Nginx Optimization**
   ```nginx
   # Add to nginx.conf
   worker_processes auto;
   worker_connections 1024;
   
   # Enable gzip compression
   gzip on;
   gzip_types text/plain application/json application/javascript text/css;
   ```

### Maintenance Tasks

1. **Weekly Tasks**
   ```bash
   # Update system packages
   sudo apt update && sudo apt upgrade -y
   
   # Clean up Docker
   docker system prune -f
   
   # Rotate logs
   sudo logrotate -f /etc/logrotate.conf
   ```

2. **Monthly Tasks**
   ```bash
   # Update SSL certificates
   sudo certbot renew
   
   # Database maintenance
   docker exec postgres_db psql -U ${DB_USER} -d ${DB_NAME} -c "VACUUM ANALYZE;"
   
   # Security updates
   sudo unattended-upgrades
   ```

## 📚 API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/user/login` | User authentication |
| GET | `/user/health` | Health check |

### User Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/user/list` | List all users |
| POST | `/user/create` | Create new user |
| DELETE | `/user/delete/{id}` | Delete user |
| GET | `/user/dbtest` | Database test |

### Company Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/company/list` | List all companies |
| POST | `/company/create` | Create new company |
| PUT | `/company/update/{id}` | Update company |
| DELETE | `/company/delete/{id}` | Delete company |
| GET | `/company/list/{username}` | List companies by officer |
| GET | `/company/health` | Health check |

### Company Temporary Updates

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/company/temp/update` | Create temporary update |
| GET | `/company/temp/list` | List pending updates |
| PUT | `/company/temp/status/{id}` | Update status |
| PUT | `/company/temp/approve/{id}` | Approve update |

### Event Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/event/list` | List all events |
| POST | `/event/create` | Create new event |

## 🔧 Troubleshooting

### Common Issues

1. **Application won't start**
   ```bash
   # Check logs
   docker-compose logs app
   
   # Check environment variables
   docker-compose config
   
   # Restart services
   docker-compose restart
   ```

2. **Database connection issues**
   ```bash
   # Check database status
   docker-compose ps postgres
   
   # Test connection
   docker exec postgres_db pg_isready -U ${DB_USER} -d ${DB_NAME}
   
   # Check logs
   docker-compose logs postgres
   ```

3. **CORS issues**
   - Verify `CORS_ALLOWED_ORIGINS` includes your frontend domain
   - Check browser console for specific CORS errors
   - Ensure protocol (http/https) matches exactly

4. **SSL certificate issues**
   ```bash
   # Check certificate status
   sudo certbot certificates
   
   # Renew certificates
   sudo certbot renew
   
   # Test Nginx configuration
   sudo nginx -t
   ```

### Emergency Procedures

1. **Service Recovery**
   ```bash
   # Quick restart
   docker-compose restart
   
   # Full rebuild
   docker-compose down
   docker-compose up -d --build
   ```

2. **Database Recovery**
   ```bash
   # Restore from backup
   gunzip /opt/place-pro/backups/db_backup_YYYYMMDD_HHMMSS.sql.gz
   docker exec -i postgres_db psql -U ${DB_USER} -d ${DB_NAME} < db_backup_YYYYMMDD_HHMMSS.sql
   ```

## 🚀 Deployment Checklist

### Pre-deployment
- [ ] Environment variables configured
- [ ] SSL certificates obtained
- [ ] Database credentials set
- [ ] CORS origins updated
- [ ] Firewall configured
- [ ] Domain DNS configured

### Deployment
- [ ] Application deployed
- [ ] Database initialized
- [ ] Nginx configured
- [ ] SSL enabled
- [ ] Health checks passing
- [ ] Monitoring setup
- [ ] Backups configured

### Post-deployment
- [ ] All endpoints tested
- [ ] Performance verified
- [ ] Logs monitored
- [ ] Security scan completed
- [ ] Documentation updated
- [ ] Team notified

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit a pull request

## 📞 Support

For deployment issues or questions:
- Check the troubleshooting section
- Review application logs
- Contact the development team

---

**Note**: This documentation assumes Ubuntu/Debian-based systems. Adjust commands accordingly for other Linux distributions. 