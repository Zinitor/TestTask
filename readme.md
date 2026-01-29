Stack:
- **Database**: PostgreSQL with GORM ORM
- **HTTP Framework**: Standard net/http
- **Logging**: zerolog
- **Configuration**: caarlos0/env
- **Documentation**: Swaggo
- **Containerization**: Docker & Docker Compose
- **Migrations**: Goose

Запуск проекта
docker compose  -f deployment/docker/docker-compose.yml  --env-file .env up  -d --build --force-recreate

Необходимые переменные окружения указываются в .env файле

Минимально необходимые переменные окружения:
```
APP_NAME="TEST_APP"
APP_VERSION=1
HTTP_PORT=:8000
LOG_LEVEL=debug
PG_POOL_MAX=10
PG_POOL_IDLE=1
PG_URL=postgres://postgres:SecurePass123!@postgres:5432/appdb?sslmode=disable
SWAGGER_ENABLED=true
```
После запуска запустите следующую команду чтобы в БД создались первоначальные таблицы
```bash 
    goose -dir internal/migrations/postgres/sql postgres 'postgres://postgres:SecurePass123!@localhost:5432/appdb?sslmode=disable' up
```


После этого приложение полностью функционирует и должно работать :) 
документацию сваггера можно найти по 
http://localhost:8000/swagger/
сделана через swaggo, не придумал как сделать чтобы в докере было достаточно подхватить бинарь чтобы не качать зависимости в контейнер и дока при этом обновлялась поэтому она часть репозитория


Команды для применения миграций

migrate down
```bash 
    goose -dir internal/migrations/postgres/sql postgres 'postgres://postgres:SecurePass123!@localhost:5432/appdb?sslmode=disable' down
```

migrate up
```bash 
    goose -dir internal/migrations/postgres/sql postgres 'postgres://postgres:SecurePass123!@localhost:5432/appdb?sslmode=disable' up
```

для запуска тестов репозиториев
```bash
go test ./internal/repos/chats/
```


для запуска бенчмарк тестов репозиториев (сделан только на get с message, т.к. он самый нагруженный, остальные показались лишними)
```bash
go test -tags=bench ./internal/repos/chats -bench=BenchmarkRepository_Get_WithPreload  -benchmem -memprofile=allocs.out -cpuprofile=cpu.out
```

Изначальный бенчмарк показал с использованием preload
```
goos: linux
goarch: amd64
pkg: ReviewTask/internal/repos/chats
cpu: AMD Ryzen 7 7700 8-Core Processor
BenchmarkRepository_Get_WithPreload-16    	    2800	    395428 ns/op	   59999 B/op	    1066 allocs/op
PASS
ok  	ReviewTask/internal/repos/chats	1.135s
```
После некоторых оптимизаций смог прийти к такому, что конечно все ещё грустно, но все равно лучше

```
goos: linux
goarch: amd64
pkg: ReviewTask/internal/repos/chats
cpu: AMD Ryzen 7 7700 8-Core Processor
BenchmarkRepository_Get_WithPreload-16    	    3808	    265084 ns/op	   46515 B/op	     728 allocs/op
PASS
ok  	ReviewTask/internal/repos/chats	1.211s
```
