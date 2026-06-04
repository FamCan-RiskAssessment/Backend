# Environment Variables

## Server Configuration

- `SERVER_PORT=` - Port number for the server
- `SERVER_MODE=` - Server mode (development, production, etc.)

## Database Configuration

- `DB_HOST=` - Database host address
- `DB_PORT=` - Database port number
- `DB_USER=` - Database username
- `DB_PASSWORD=` - Database password
- `DB_NAME=` - Database name

## Redis Configuration

- `RDB_PORT=` - Redis port number
- `RDB_ADDRESS=` - Redis server address
- `RDB_PASSWORD=` - Redis password
- `RDB_NUMBER=` - Redis database number

## PostgreSQL Configuration

- `POSTGRES_USER=` - PostgreSQL username
- `POSTGRES_PASSWORD=` - PostgreSQL password
- `POSTGRES_DB=` - PostgreSQL database name

## External Services

- `GITHUB_TOKEN=` - GitHub API token

## SMS Service Configuration

- `SMS_PROVIDER=` - SMS service provider
- `SMS_GATEWAY_API_KEY=` - SMS gateway API key
- `SMS_USERNAME=` - SMS service username
- `SMS_PASSWORD=` - SMS service password
- `SMS_SOURCE=` - SMS source identifier

## MinIO Configuration

- `MINIO_REGION=` - MinIO region (e.g. `us-east-1`)
- `MINIO_ACCESS_KEY=` - MinIO access key
- `MINIO_SECRET_KEY=` - MinIO secret key
- `MINIO_ENDPOINT=` - MinIO server URL (e.g. `http://localhost:9000`)
- `MINIO_BUCKET=` - Single bucket for all file storage

## Admin Configuration

- `SUPER_ADMIN_PHONE=` - Super admin phone number
- `SUPER_ADMIN_PASSWORD=` - Super admin password

## Third-party APIs

- `PREMM5_API_URL=` - PREMM5 API endpoint
- `BCRA_API_URL=` - BCRA API endpoint
- `GAIL_API_URL=` - GAIL API endpoint
- `PLCO_API_URL=` - PLCO API endpoint

## MinIO Object Prefixes

Former per-type buckets are directories (key prefixes) inside `MINIO_BUCKET`. Defaults apply when unset.

- `MINIO_PREFIX_MAMOGRAPHY=` - Default: `mamography`
- `MINIO_PREFIX_CANCER=` - Default: `cancer`
- `MINIO_PREFIX_GENETIC_TEST=` - Default: `genetic-test`
- `MINIO_PREFIX_FATHER_GENETIC_TEST=` - Default: `father-genetic-test`
- `MINIO_PREFIX_MOTHER_GENETIC_TEST=` - Default: `mother-genetic-test`
