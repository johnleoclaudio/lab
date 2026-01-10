# Database

- brew install golang-migrate
- migrate create -seq -ext=.sql -dir=./migrations create_users_table
- there's CHECK constraints in Postgres


```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: postgres_dev
    restart: unless-stopped
    environment:
      POSTGRES_USER: devuser
      POSTGRES_PASSWORD: devpassword
      POSTGRES_DB: devdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```


```bash
# Install PostgreSQL client only
brew install libpq
brew link --force libpq

# Connect
psql -h localhost -p 5432 -U devuser -d devdb
```

**Connection string for your application:**
```
postgresql://devuser:devpassword@localhost:5432/devdb
```
