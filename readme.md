
migrate down
```bash 
    goose -dir internal/migrations/postgres/sql postgres 'postgres://postgres:SecurePass123!@localhost:5432/appdb?sslmode=disable' down
```

migrate up
```bash 
    goose -dir internal/migrations/postgres/sql postgres 'postgres://postgres:SecurePass123!@localhost:5432/appdb?sslmode=disable' up
```

docker compose -f deployment/docker/docker-compose.yml up




Запуск проекта
docker compose  -f deployment/docker/docker-compose.yml  --env-file .env up  -d --build --force-recreate

Необходимые переменные окружения указываются в .env файле
APP_NAME="TEST_APP"
APP_VERSION=1
HTTP_PORT=:8000
LOG_LEVEL=debug
PG_POOL_MAX=10
PG_POOL_IDLE=1
PG_URL=postgres://postgres:SecurePass123!@postgres:5432/appdb?sslmode=disable --in docker-compose
# PG_URL=postgres://postgres:SecurePass123!@localhost:5432/appdb?sslmode=disable  -- make dev
SWAGGER_ENABLED=true

После запуска запустите следующую команду чтобы в БД создались первоначальные таблицы (Я более чем уверен что это можно делать автоматически через какой-то скрипт но не придумал как)
```bash 
    goose -dir internal/migrations/postgres/sql postgres 'postgres://postgres:SecurePass123!@localhost:5432/appdb?sslmode=disable' up
```


после запуска можно найти документацию сваггера по 
http://localhost:8000/swagger/
сделана через swaggo, не придумал как сделать чтобы в докере было достаточно подхватить бинарь чтобы не качать зависимости в контейнер и дока обновлялась поэтому она часть репозитория
