# Kodia Framework — Production Deployment Guide

This guide covers deploying Kodia (backend + dependencies) to production using Docker, docker-compose, or Kubernetes.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start with Docker Compose](#quick-start-with-docker-compose)
3. [Kubernetes Deployment](#kubernetes-deployment)
4. [Bare Metal / VPS Deployment](#bare-metal--vps-deployment)
5. [Environment Variables Reference](#environment-variables-reference)
6. [Database Migrations](#database-migrations)
7. [SSL/TLS Setup](#ssltls-setup)
8. [Health Checks & Monitoring](#health-checks--monitoring)
9. [Scaling & Performance](#scaling--performance)
10. [Backup & Disaster Recovery](#backup--disaster-recovery)
11. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### System Requirements

- **Docker & Docker Compose** (for docker-compose deployment)
  - Docker Engine 20.10+
  - Docker Compose 2.0+

- **Kubernetes** (for K8s deployment)
  - kubectl 1.24+
  - Kubernetes cluster 1.24+
  - helm (optional, for package management)

- **VPS / Bare Metal**
  - Ubuntu 20.04 LTS or later
  - 2+ CPU cores
  - 4GB+ RAM
  - 20GB+ disk space

### Pre-deployment Checklist

- [ ] Go binary built and tested locally
- [ ] Environment variables prepared (database credentials, JWT secrets, etc.)
- [ ] PostgreSQL database provisioned or plan to use docker-compose postgres service
- [ ] Redis instance provisioned or plan to use docker-compose redis service
- [ ] Domain name registered and DNS configured
- [ ] SSL certificate (or use Let's Encrypt with cert-manager)
- [ ] Backup strategy planned (PostgreSQL + Redis)

---

## Quick Start with Docker Compose

### 1. Prepare Environment Variables

Create `.env.prod` in the project root:

```bash
# Database
APP_DATABASE_HOST=postgres
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=kodia_prod
APP_DATABASE_USER=kodia
APP_DATABASE_PASSWORD=your-secure-password-here
APP_DATABASE_SSL_MODE=disable

# Redis
APP_REDIS_HOST=redis
APP_REDIS_PORT=6379
APP_REDIS_PASSWORD=your-redis-password-here
APP_REDIS_DB=0

# Application
APP_NAME=Kodia
APP_ENV=production
APP_PORT=8080
APP_TIMEZONE=UTC

# JWT Secrets
APP_JWT_ACCESS_SECRET=your-access-secret-here
APP_JWT_REFRESH_SECRET=your-refresh-secret-here
APP_JWT_ACCESS_EXPIRES_IN=3600
APP_JWT_REFRESH_EXPIRES_IN=2592000

# SMTP (for email)
MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USER=no-reply@example.com
MAIL_PASSWORD=your-smtp-password
MAIL_FROM_ADDR=no-reply@example.com
MAIL_FROM_NAME=Kodia

# Observability (optional)
APP_OBSERVABILITY_SENTRY_DSN=

# Search (optional)
# APP_SEARCH_MASTER_KEY=

# Storage (optional)
# APP_STORAGE_ACCESS_ID=
# APP_STORAGE_SECRET_KEY=
```

**⚠️ Security**: Never commit `.env.prod` to git. Use secret management tools (e.g., HashiCorp Vault, AWS Secrets Manager).

### 2. Build & Start Services

```bash
# Build the app image
docker-compose -f docker-compose.prod.yml build

# Start all services (postgres, redis, app)
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f app

# Stop services
docker-compose -f docker-compose.prod.yml down
```

### 3. Run Database Migrations

Once services are running:

```bash
# Connect to the app container and run migrations
docker-compose -f docker-compose.prod.yml exec app ./server migrate:latest

# Or seed data (optional)
docker-compose -f docker-compose.prod.yml exec app ./server seed:run
```

### 4. Access the Application

- **API**: `http://localhost:8080`
- **Health check**: `http://localhost:8080/health`
- **Metrics**: `http://localhost:9090/metrics` (Prometheus format)

---

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster running (1.24+)
- kubectl configured
- cert-manager installed (for TLS)
- NGINX Ingress Controller installed
- Metrics Server installed (for HPA)

### 1. Install Dependencies

```bash
# cert-manager (for Let's Encrypt)
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# NGINX Ingress Controller
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install nginx-ingress ingress-nginx/ingress-nginx -n ingress-nginx --create-namespace

# Metrics Server (if not already installed)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

### 2. Create Secrets

```bash
# Generate JWT secrets
export JWT_ACCESS_SECRET=$(openssl rand -base64 32)
export JWT_REFRESH_SECRET=$(openssl rand -base64 32)
export DB_PASSWORD=$(openssl rand -base64 16)
export REDIS_PASSWORD=$(openssl rand -base64 16)

# Create secret in cluster
kubectl create secret generic kodia-secret \
  --from-literal=APP_DATABASE_PASSWORD=$DB_PASSWORD \
  --from-literal=APP_REDIS_PASSWORD=$REDIS_PASSWORD \
  --from-literal=APP_JWT_ACCESS_SECRET=$JWT_ACCESS_SECRET \
  --from-literal=APP_JWT_REFRESH_SECRET=$JWT_REFRESH_SECRET \
  --from-literal=MAIL_HOST=smtp.example.com \
  --from-literal=MAIL_PORT=587 \
  --from-literal=MAIL_USER=no-reply@example.com \
  --from-literal=MAIL_PASSWORD=your-smtp-password \
  --from-literal=MAIL_FROM_ADDR=no-reply@example.com \
  --from-literal=MAIL_FROM_NAME=Kodia \
  --dry-run=client -o yaml | kubectl apply -f -

# Create configmap
kubectl create configmap kodia-config \
  --from-literal=APP_NAME=Kodia \
  --from-literal=APP_ENV=production \
  --from-literal=APP_PORT=8080 \
  --from-literal=APP_DATABASE_HOST=postgres \
  --from-literal=APP_DATABASE_PORT=5432 \
  --from-literal=APP_DATABASE_NAME=kodia_prod \
  --from-literal=APP_DATABASE_USER=kodia \
  --from-literal=APP_REDIS_HOST=redis \
  --from-literal=APP_REDIS_PORT=6379 \
  --from-literal=APP_REDIS_DB=0 \
  --dry-run=client -o yaml | kubectl apply -f -
```

### 3. Deploy with kustomize

```bash
# Review what will be deployed
kubectl apply -k k8s/ --dry-run=client -o yaml

# Apply the manifests
kubectl apply -k k8s/

# Watch deployment rollout
kubectl rollout status deployment/kodia-api -w

# Verify pods are running
kubectl get pods -l app.kubernetes.io/name=kodia
```

### 4. Run Database Migrations

```bash
# Get pod name
POD=$(kubectl get pod -l app.kubernetes.io/name=kodia -o jsonpath='{.items[0].metadata.name}')

# Run migrations
kubectl exec -it $POD -- ./server migrate:latest

# Seed data (optional)
kubectl exec -it $POD -- ./server seed:run
```

### 5. Configure Ingress

Edit `k8s/ingress.yaml` and update:
- `api.example.com` → your actual domain
- `admin@example.com` → your email for Let's Encrypt

```bash
kubectl apply -f k8s/ingress.yaml

# Check ingress status
kubectl get ingress kodia-api
kubectl describe ingress kodia-api

# View certificate status
kubectl describe certificate kodia-api-tls
```

### 6. Access the Application

```bash
# Get the ingress external IP
kubectl get ingress kodia-api -o wide

# Health check
curl https://api.example.com/health

# Metrics
curl https://api.example.com:9090/metrics
```

### Kubernetes Scaling

The HPA automatically scales between 2–10 replicas based on CPU (70%) and memory (80%) utilization.

```bash
# View current HPA status
kubectl get hpa kodia-api -w

# View HPA events
kubectl describe hpa kodia-api

# Manually scale (overrides HPA)
kubectl scale deployment kodia-api --replicas=5
```

---

## Bare Metal / VPS Deployment

### 1. Server Setup

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

### 2. Clone & Configure

```bash
# Clone repository
git clone https://github.com/yourusername/kodia.git /opt/kodia
cd /opt/kodia

# Create environment file
nano .env.prod

# Build Docker image
docker-compose -f docker-compose.prod.yml build

# Start services
docker-compose -f docker-compose.prod.yml up -d
```

### 3. Reverse Proxy (Nginx)

```bash
# Install Nginx
sudo apt install -y nginx

# Create Nginx config
sudo tee /etc/nginx/sites-available/kodia > /dev/null <<'EOF'
upstream kodia_app {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name api.example.com;

    client_max_body_size 100M;

    location / {
        proxy_pass http://kodia_app;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
EOF

# Enable site
sudo ln -sf /etc/nginx/sites-available/kodia /etc/nginx/sites-enabled/kodia
sudo rm -f /etc/nginx/sites-enabled/default

# Test and restart
sudo nginx -t
sudo systemctl restart nginx
```

### 4. SSL/TLS with Certbot

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Get certificate (automated renewal)
sudo certbot --nginx -d api.example.com --email admin@example.com --agree-tos --non-interactive --redirect

# Verify auto-renewal
sudo certbot renew --dry-run
```

### 5. System Service (systemd)

Create `/etc/systemd/system/kodia.service`:

```ini
[Unit]
Description=Kodia Application
After=network.target docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/kodia
ExecStart=/usr/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/bin/docker-compose -f docker-compose.prod.yml down
Restart=always
RestartSec=10s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable kodia
sudo systemctl start kodia
sudo systemctl status kodia
```

---

## Environment Variables Reference

### Required

| Variable | Example | Description |
|----------|---------|-------------|
| `APP_NAME` | `Kodia` | Application name |
| `APP_ENV` | `production` | Environment: `production`, `staging`, `development` |
| `APP_PORT` | `8080` | Application HTTP port |
| `APP_DATABASE_HOST` | `postgres` | PostgreSQL host |
| `APP_DATABASE_PORT` | `5432` | PostgreSQL port |
| `APP_DATABASE_NAME` | `kodia_prod` | Database name |
| `APP_DATABASE_USER` | `kodia` | Database user |
| `APP_DATABASE_PASSWORD` | `***` | Database password (min 16 chars) |
| `APP_REDIS_HOST` | `redis` | Redis host |
| `APP_REDIS_PORT` | `6379` | Redis port |
| `APP_REDIS_PASSWORD` | `***` | Redis password (min 16 chars) |
| `APP_JWT_ACCESS_SECRET` | `***` | JWT access token secret (32+ chars) |
| `APP_JWT_REFRESH_SECRET` | `***` | JWT refresh token secret (32+ chars) |

### Optional

| Variable | Example | Default | Description |
|----------|---------|---------|-------------|
| `APP_TIMEZONE` | `UTC` | `UTC` | Server timezone |
| `APP_DATABASE_SSL_MODE` | `require` | `disable` | PostgreSQL SSL: `disable`, `require`, `verify-full` |
| `APP_REDIS_DB` | `0` | `0` | Redis DB number |
| `APP_JWT_ACCESS_EXPIRES_IN` | `3600` | `3600` | Access token TTL (seconds) |
| `APP_JWT_REFRESH_EXPIRES_IN` | `2592000` | `2592000` | Refresh token TTL (seconds) |
| `APP_CORS_ALLOWED_ORIGINS` | `https://example.com` | `*` | Comma-separated CORS origins |
| `APP_RATE_LIMIT_REQUESTS` | `100` | `100` | Requests per minute |
| `MAIL_HOST` | `smtp.example.com` | — | SMTP server host |
| `MAIL_PORT` | `587` | `587` | SMTP port |
| `MAIL_USER` | `no-reply@example.com` | — | SMTP username |
| `MAIL_PASSWORD` | `***` | — | SMTP password |
| `MAIL_FROM_ADDR` | `no-reply@example.com` | — | Default From address |
| `MAIL_FROM_NAME` | `Kodia` | `Kodia` | Default From name |
| `APP_OBSERVABILITY_SENTRY_DSN` | `https://***@sentry.io/***` | — | Sentry error tracking DSN |
| `APP_SEARCH_MASTER_KEY` | `***` | — | Meilisearch master key |
| `APP_STORAGE_ACCESS_ID` | `***` | — | S3 access key |
| `APP_STORAGE_SECRET_KEY` | `***` | — | S3 secret key |
| `APP_STORAGE_BUCKET` | `kodia-files` | — | S3 bucket name |
| `APP_STORAGE_REGION` | `us-east-1` | — | S3 region |

---

## Database Migrations

### Running Migrations

In docker-compose:

```bash
# List pending migrations
docker-compose -f docker-compose.prod.yml exec app ./server migrate:status

# Run all pending migrations
docker-compose -f docker-compose.prod.yml exec app ./server migrate:latest

# Rollback one migration
docker-compose -f docker-compose.prod.yml exec app ./server migrate:rollback
```

In Kubernetes:

```bash
POD=$(kubectl get pod -l app.kubernetes.io/name=kodia -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it $POD -- ./server migrate:status
kubectl exec -it $POD -- ./server migrate:latest
```

### Initial Setup

After first deployment:

```bash
# Run migrations
docker-compose -f docker-compose.prod.yml exec app ./server migrate:latest

# Seed initial data (users, roles, etc.)
docker-compose -f docker-compose.prod.yml exec app ./server seed:run
```

---

## SSL/TLS Setup

### Kubernetes with cert-manager (Recommended)

cert-manager automatically provisions and renews certificates via Let's Encrypt.

1. **Install cert-manager**:
   ```bash
   kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
   ```

2. **Create ClusterIssuer** (already in `k8s/ingress.yaml`):
   ```bash
   kubectl apply -f k8s/ingress.yaml
   ```

3. **Verify certificate status**:
   ```bash
   kubectl get certificates
   kubectl describe certificate kodia-api-tls
   ```

Certificate is auto-renewed 30 days before expiration.

### Docker Compose with Let's Encrypt

```bash
# Install certbot
sudo apt install -y certbot

# Get certificate
sudo certbot certonly --standalone -d api.example.com

# Update docker-compose.prod.yml to mount certs:
# volumes:
#   - /etc/letsencrypt/live/api.example.com:/etc/certs:ro

# Renew automatically (cron)
sudo certbot renew --quiet --cron
```

### Manual Certificate Upload (Self-Signed)

```bash
# Generate self-signed cert (development only)
openssl req -x509 -newkey rsa:4096 -nodes -out cert.pem -keyout key.pem -days 365

# Mount in Kubernetes secret
kubectl create secret tls kodia-api-tls --cert=cert.pem --key=key.pem
```

---

## Health Checks & Monitoring

### Health Check Endpoint

```bash
# Local
curl http://localhost:8080/health

# Expected response
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Prometheus Metrics

Metrics available at `/metrics` (port 9090 in K8s, 8080 in docker-compose):

```bash
# View all metrics
curl http://localhost:9090/metrics

# Key metrics:
# - http_requests_total (by method, path, status)
# - http_request_duration_seconds
# - database_query_duration_seconds
# - cache_hits_total, cache_misses_total
```

### Kubernetes Liveness & Readiness Probes

Already configured in `k8s/deployment.yaml`:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

### Monitoring with Prometheus + Grafana

```bash
# Install Prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

# Port-forward to Grafana
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80

# Access at http://localhost:3000 (admin/prom-operator)
```

---

## Scaling & Performance

### Kubernetes Auto-Scaling

HPA automatically scales based on metrics:

```bash
# Current status
kubectl get hpa kodia-api

# Scale manually (if needed)
kubectl scale deployment kodia-api --replicas=5

# View HPA events
kubectl logs -n kube-system deployment/metrics-server
```

### Performance Tuning

1. **Database Connection Pooling**:
   ```
   APP_DATABASE_MAX_CONNS=20
   APP_DATABASE_MIN_CONNS=5
   ```

2. **Redis Connection Pool**:
   ```
   APP_REDIS_MAX_RETRIES=3
   ```

3. **HTTP Timeouts**:
   ```
   APP_HTTP_READ_TIMEOUT=15s
   APP_HTTP_WRITE_TIMEOUT=15s
   ```

4. **Resource Limits** (K8s):
   Edit `k8s/deployment.yaml`:
   ```yaml
   resources:
     requests:
       memory: "256Mi"
       cpu: "200m"
     limits:
       memory: "512Mi"
       cpu: "1000m"
   ```

---

## Backup & Disaster Recovery

### PostgreSQL Backups

#### Docker Compose

```bash
# Full backup
docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U kodia kodia_prod > backup.sql

# Restore
cat backup.sql | docker-compose -f docker-compose.prod.yml exec -T postgres psql -U kodia kodia_prod

# Automated daily backup (cron)
0 2 * * * cd /opt/kodia && docker-compose -f docker-compose.prod.yml exec -T postgres pg_dump -U kodia kodia_prod | gzip > backups/kodia_$(date +\%Y\%m\%d).sql.gz
```

#### Kubernetes

```bash
# Port-forward PostgreSQL
kubectl port-forward svc/postgres 5432:5432 &

# Backup via psql client
pg_dump -h localhost -U kodia kodia_prod > backup.sql

# Schedule with CronJob:
kubectl apply -f - <<'EOF'
apiVersion: batch/v1
kind: CronJob
metadata:
  name: kodia-postgres-backup
spec:
  schedule: "0 2 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: postgres
            image: postgres:16-alpine
            command:
            - /bin/sh
            - -c
            - pg_dump -h postgres kodia_prod | gzip > /backups/kodia_$(date +%Y%m%d).sql.gz
            env:
            - name: PGUSER
              value: kodia
            - name: PGPASSWORD
              valueFrom:
                secretKeyRef:
                  name: kodia-secret
                  key: APP_DATABASE_PASSWORD
          restartPolicy: OnFailure
EOF
```

### Redis Backups

```bash
# Docker Compose
docker-compose -f docker-compose.prod.yml exec redis redis-cli BGSAVE

# Kubernetes
POD=$(kubectl get pod -l app=redis -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it $POD -- redis-cli BGSAVE

# Verify (creates dump.rdb)
docker-compose -f docker-compose.prod.yml exec redis redis-cli LASTSAVE
```

### Disaster Recovery Plan

1. **RTO (Recovery Time Objective)**: < 1 hour
2. **RPO (Recovery Point Objective)**: < 1 day
3. **Backup Location**: External (S3, GCS, or off-site)
4. **Test Recovery**: Monthly restore drill

```bash
# Test PostgreSQL recovery
pg_restore -d test_db backup.sql

# Test Redis recovery
redis-cli --rdb /tmp/dump.rdb
```

---

## Troubleshooting

### Logs

#### Docker Compose

```bash
# All logs
docker-compose -f docker-compose.prod.yml logs -f

# App only
docker-compose -f docker-compose.prod.yml logs -f app

# Last 100 lines
docker-compose -f docker-compose.prod.yml logs --tail=100 app
```

#### Kubernetes

```bash
# Pod logs
kubectl logs deployment/kodia-api
kubectl logs pod/kodia-api-xyz123 -f

# Previous logs (if crashed)
kubectl logs pod/kodia-api-xyz123 --previous

# Describe pod for events
kubectl describe pod/kodia-api-xyz123
```

### Common Issues

#### Application Won't Start

```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs app

# Verify database connectivity
docker-compose -f docker-compose.prod.yml exec app ./server db:test

# Check environment variables
docker-compose -f docker-compose.prod.yml config | grep APP_
```

#### Database Connection Fails

```bash
# Verify PostgreSQL is running
docker-compose -f docker-compose.prod.yml exec postgres pg_isready

# Test connection
docker-compose -f docker-compose.prod.yml exec postgres psql -U kodia -d kodia_prod -c "SELECT 1"

# Check password and host in .env.prod
```

#### Redis Connection Fails

```bash
# Test Redis
docker-compose -f docker-compose.prod.yml exec redis redis-cli ping

# Check AUTH
docker-compose -f docker-compose.prod.yml exec redis redis-cli -a "$REDIS_PASSWORD" ping
```

#### Out of Memory

```bash
# Check pod memory usage
kubectl top pods
kubectl top node

# Increase container limits in k8s/deployment.yaml
# or
# Increase docker memory in docker-compose.prod.yml

# Restart to clear cache
kubectl rollout restart deployment/kodia-api
```

#### Certificate Not Renewing

```bash
# Check cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager

# Manually trigger renewal
kubectl annotate certificate kodia-api-tls cert-manager.io/issue-temporary-certificate=true --overwrite

# Debug certificate
kubectl describe certificate kodia-api-tls
```

---

## Support & Questions

For issues, check:
- Application logs
- Database connectivity
- Environment variables
- Kubernetes events: `kubectl get events -A`
- cert-manager issues: `kubectl describe certificate kodia-api-tls`

See main [README.md](../README.md) for community support channels.
