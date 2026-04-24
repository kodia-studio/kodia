# Database Setup Guide

Kodia Framework mendukung multiple database backends dengan konfigurasi yang fleksibel.

## 🚀 Quick Start (SQLite - Recommended for Development)

SQLite adalah database yang paling sederhana untuk memulai. **Tidak perlu setup apapun** — file database otomatis dibuat di `storage/app.db`.

Konfigurasi default `.env`:
```env
APP_DATABASE_DRIVER=sqlite
APP_DATABASE_PATH=storage/app.db
```

Itu saja! Sekarang Anda bisa langsung:
```bash
kodia migrate
kodia dev
```

### Keuntungan SQLite untuk Development:
- ✅ Tidak perlu install database server
- ✅ Satu file (mudah di-backup, share, version control)
- ✅ Setup instant, langsung bisa develop
- ✅ Perfekt untuk prototyping & testing

---

## PostgreSQL Setup (Production-Ready)

Untuk production atau project yang lebih kompleks, gunakan PostgreSQL.

### 1. Install PostgreSQL

**macOS:**
```bash
brew install postgresql
brew services start postgresql
```

**Ubuntu/Debian:**
```bash
sudo apt-get install postgresql postgresql-contrib
sudo service postgresql start
```

**Docker:**
```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:latest
```

### 2. Create Database & User

```bash
# Masuk ke PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE kodia_db;

# Create user
CREATE USER kodia WITH PASSWORD 'postgres';

# Grant privileges
ALTER ROLE kodia CREATEDB;
GRANT ALL PRIVILEGES ON DATABASE kodia_db TO kodia;

# Exit
\q
```

### 3. Update `.env`

```env
APP_DATABASE_DRIVER=postgres
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=kodia_db
APP_DATABASE_USER=kodia
APP_DATABASE_PASSWORD=postgres
APP_DATABASE_SSL_MODE=disable
APP_DATABASE_TIMEZONE=UTC
```

### 4. Run Migrations

```bash
kodia migrate
```

### 5. Verify Connection

```bash
kodia db:status
```

---

## MySQL Setup (Alternative)

### 1. Install MySQL

**macOS:**
```bash
brew install mysql
brew services start mysql
```

**Ubuntu/Debian:**
```bash
sudo apt-get install mysql-server
sudo service mysql start
```

### 2. Create Database & User

```bash
# Connect to MySQL
mysql -u root

# Create database
CREATE DATABASE kodia_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# Create user
CREATE USER 'kodia'@'localhost' IDENTIFIED BY 'password';

# Grant privileges
GRANT ALL PRIVILEGES ON kodia_db.* TO 'kodia'@'localhost';
FLUSH PRIVILEGES;

# Exit
EXIT;
```

### 3. Update `.env`

```env
APP_DATABASE_DRIVER=mysql
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=3306
APP_DATABASE_NAME=kodia_db
APP_DATABASE_USER=kodia
APP_DATABASE_PASSWORD=password
APP_DATABASE_TIMEZONE=UTC
```

### 4. Run Migrations

```bash
kodia migrate
```

---

## Docker Compose Setup (Recommended)

Untuk development dengan Docker, gunakan `kodia sail`:

```bash
kodia sail init:docker
kodia sail up
```

Ini akan start PostgreSQL, Redis, dan services lainnya secara otomatis.

Atau, setup manual dengan docker-compose:

```yaml
version: '3'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: kodia_db
      POSTGRES_USER: kodia
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Run:
```bash
docker-compose up -d
kodia migrate
```

---

## Configuration Reference

### SQLite
```env
APP_DATABASE_DRIVER=sqlite
APP_DATABASE_PATH=storage/app.db
```

**Variables:**
- `APP_DATABASE_PATH` — Path ke file SQLite (relative to project root)

### PostgreSQL
```env
APP_DATABASE_DRIVER=postgres
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=kodia_db
APP_DATABASE_USER=kodia
APP_DATABASE_PASSWORD=secure-password
APP_DATABASE_SSL_MODE=disable    # disable, allow, prefer, require, verify-ca, verify-full
APP_DATABASE_TIMEZONE=UTC
```

**Variables:**
- `APP_DATABASE_HOST` — PostgreSQL server host
- `APP_DATABASE_PORT` — PostgreSQL server port (default: 5432)
- `APP_DATABASE_NAME` — Database name
- `APP_DATABASE_USER` — Database user
- `APP_DATABASE_PASSWORD` — Database password
- `APP_DATABASE_SSL_MODE` — SSL connection mode
- `APP_DATABASE_TIMEZONE` — Timezone untuk timestamps

### MySQL
```env
APP_DATABASE_DRIVER=mysql
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=3306
APP_DATABASE_NAME=kodia_db
APP_DATABASE_USER=kodia
APP_DATABASE_PASSWORD=password
APP_DATABASE_TIMEZONE=UTC
```

---

## Migration Commands

Setelah configure database:

```bash
# Run all pending migrations
kodia migrate

# Rollback last batch of migrations
kodia migrate:rollback

# Show migration status
kodia migrate:status

# Reset database (rollback all, then migrate)
kodia db:reset

# Fresh start (drop all, then migrate + seed)
kodia db:fresh
```

---

## Switching Databases

Untuk switch dari SQLite ke PostgreSQL:

1. **Backup data (optional)**
   ```bash
   cp storage/app.db storage/app.db.backup
   ```

2. **Update `.env`**
   ```env
   APP_DATABASE_DRIVER=postgres
   APP_DATABASE_HOST=localhost
   # ... other postgres config
   ```

3. **Re-run migrations**
   ```bash
   kodia migrate:rollback
   kodia migrate
   ```

---

## Environment-Specific Configuration

**Development:**
```env
APP_ENV=development
APP_DATABASE_DRIVER=sqlite
APP_DATABASE_PATH=storage/app.db
```

**Testing:**
```env
APP_ENV=testing
APP_DATABASE_DRIVER=sqlite
APP_DATABASE_PATH=storage/test.db
```

**Production:**
```env
APP_ENV=production
APP_DATABASE_DRIVER=postgres
APP_DATABASE_HOST=db.example.com
# ... use strong credentials!
```

---

## Troubleshooting

### "database is locked" (SQLite)
SQLite tidak cocok untuk concurrent writes. Gunakan PostgreSQL untuk production.

### "connection refused" (PostgreSQL)
- Pastikan PostgreSQL server running: `psql -U postgres`
- Check host & port di `.env`
- Verify credentials

### "Unknown database" 
- Pastikan database sudah created: `CREATE DATABASE kodia_db;`
- Check APP_DATABASE_NAME di `.env`

### "Access denied"
- Verify username & password
- Check user privileges: `SHOW GRANTS FOR kodia@localhost;` (MySQL)

---

## Best Practices

✅ **Do:**
- Use SQLite untuk development (instant, no setup)
- Use PostgreSQL untuk staging/production
- Keep `.env` with default values safe
- Regular backups untuk production data
- Use SSL_MODE=require untuk production PostgreSQL

❌ **Don't:**
- Commit `.env` ke Git (use `.env.example` instead)
- Use weak passwords di production
- Change database driver frequently (test before switching)
- Share database credentials

---

## See Also

- [Migrations Guide](MIGRATIONS.md) — How to create & manage schema changes
- [ORM Guide](ORM_GUIDE.md) — Database queries dengan GORM
- [Seeding Guide](SEEDING.md) — Populate database dengan test data
