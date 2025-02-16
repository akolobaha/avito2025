### API сервис по покупке мерча

#### Поднять контейнеры
```
docker-compose up
```

#### Миграции
Подключиться к контейнеру с приложением
```
docker exec -it app-merch bash 
```

Установить goose, если не установлен
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Применить миграции
```
goose -dir=migrations postgres "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up
```
Откатить миграции
```
goose -dir=migrations postgres "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" down
```

Статус миграций
```
goose -dir=migrations postgres "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" status
```


[openapi.yaml](api/openapi/schema.json)

docker run -d -p 8081:8080 -e SWAGGER_JSON=/api/openapi/schema.json -v $(pwd)/api/openapi:/api/openapi swaggerapi/swagger-ui
