# Development Docker Setup

This setup provides an isolated Docker deployment for frontend development with additional seeding data not available in production.

## Features

- **Isolated Environment**: Separate containers with different ports and names from production
- **Development Seeding**: Includes sample users, organizations, and venues for testing
- **Persistent Data**: Database and MinIO data persist between container restarts
- **Easy Access**: Different ports to avoid conflicts with production setup

## Ports

- Backend API: `http://localhost:8081`
- PostgreSQL: `localhost:5433`
- MinIO API: `localhost:9002`
- MinIO Console: `http://localhost:9003`

## Sample Data

When `APP_ENV=development`, the seeder creates:

### Users

- **Admin**: `admin@inspacemap.dev` / `admin123` (Owner role)
- **Editor**: `editor@inspacemap.dev` / `editor123` (Editor role)
- **Viewer**: `viewer@inspacemap.dev` / `viewer123` (Viewer role)

### Organization

- **Demo Organization** (slug: `demo-org`)

### Venue

- **Demo Venue** (slug: `demo-venue`) - Sample venue in Jakarta with draft revision

## Usage

1. Ensure you have a `.env` file with required environment variables (same as production)
2. Start the development environment:
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```
3. The backend will be available at `http://localhost:8081`
4. MinIO console at `http://localhost:9003` (credentials from .env)

## Stopping

```bash
docker-compose -f docker-compose.dev.yml down
```

## Cleaning Up

To remove containers and volumes (including data):

```bash
docker-compose -f docker-compose.dev.yml down -v
```

## Frontend Development

Use this backend setup while developing the frontend locally. The frontend can connect to `http://localhost:8081` for API calls.

For full frontend setup, refer to `FRONTEND_GUIDE.md`.
