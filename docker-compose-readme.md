# Directus Docker Setup

This Docker Compose configuration runs Directus with a PostgreSQL database.

## Quick Start

1. **Start the services:**
   ```bash
   docker-compose up -d
   ```

2. **Access Directus:**
   - URL: http://localhost:8055
   - Email: admin@example.com
   - Password: admin123

3. **Stop the services:**
   ```bash
   docker-compose down
   ```

## Configuration

### Environment Variables
Copy `.env.example` to `.env` and customize:

```bash
cp .env.example .env
```

**Important:** Replace the default `KEY` and `SECRET` values with secure random strings:

```bash
# Generate secure keys
openssl rand -hex 32  # for KEY
openssl rand -hex 64  # for SECRET
```

### Services

- **PostgreSQL** (port 5432): Database backend
- **Directus** (port 8055): Headless CMS

### Volumes
- `postgres_data`: PostgreSQL data persistence
- `directus_uploads`: User upload files
- `directus_extensions`: Custom extensions

## Development

### View logs
```bash
docker-compose logs -f directus
docker-compose logs -f postgres
```

### Reset data
```bash
docker-compose down -v  # Removes all volumes
```

### Update Directus
```bash
docker-compose pull directus
docker-compose up -d --force-recreate
```

## Production Considerations

1. **Change default credentials** in environment variables
2. **Use secure KEY/SECRET** values
3. **Set up backups** for postgres_data volume
4. **Configure reverse proxy** (nginx/caddy) for HTTPS
5. **Set appropriate memory limits** in docker-compose.yml