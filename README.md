# Агрегатор новостей 🚀

## Endpoint: Получение новостей

`GET /news/{n}`

### Описание

Возвращает `n` последних новостных постов в формате JSON, отсортированных по дате публикации (новые первыми).

### Параметры

| Параметр | Тип    | Обязательный   | Описание                          | Валидация    |
|----------|--------|----------------|-----------------------------------|--------------|
| n        | integer| Да             | Количество возвращаемых записей   | 1 ≤ n ≤ 1000 |

### Пример запроса

```bash
curl -X GET 'http://localhost:8088/news/5' -H 'Accept: application/json'
```

### Пример ответа

```json
[
    {
        "id": "uuidv5",
        "title": "Заголовок новости",
        "content": "Текст новости (может включать html разметку)",
        "published": "2025-03-29T12:00:00Z",
        "link": "https://source.com/full-article"
    }
]
```

## Тестирование

Для запуска некоторых тестов требуется поднятый контейнер с Postgres.

```bash
export POSTGRES_PASSWORD='some_pass'
export POSTGRES_HOST='localhost'
export POSTGRES_PORT=5432
chmod +x cmd/run_postgres.sh
./cmd/run_postgres.sh
# Check if the Postgres container is ready to accept connections.
while ! pg_isready -h localhost -p 5432; do sleep 1; done
go test -v -cover ./...
```

## Запуск

### Запуск  в режиме разработки

В режиме разработки приложение использует in-memory базу данных

```bash
go run cmd/server/main.go -dev
```

### Запуск в обычном режиме

Для запуска в обычном режиме необходимо:

1. Установить переменные окружения.

    ```bash
    export POSTGRES_PASSWORD='some_pass'
    export POSTGRES_HOST='localhost'
    export POSTGRES_PORT=5432
    ```

2. Запустить контейнер с Postgres.

    ```bash
    chmod +x cmd/run_postgres.sh
    ./cmd/run_postgres.sh
    ```

3. Запустить сервер.

    ```bash
    go run cmd/server/main.go
    ```

### Развертывание в контейнерах вместе с фронтендом

1. Установить переменные окружения.

    ```bash
    export POSTGRES_PASSWORD='some_pass'
    export POSTGRES_HOST='db'
    export POSTGRES_PORT=5432
    ```

    *Пароль можно установить любой. Хосты и порты пока не конфигурируются, лучше оставить как в примере, чтобы не нарушить взаимосвязь компонентов системы.*

2. Поднять контейнеры при помощи **Docker Compose**.

    ```bash
    docker compose up --build
    ```

[Веб-морда](http://localhost:8080/)

## TODO

1. ~~Верификация постов перед пакетной записью в БД.~~
2. ~~Ограничить максимальное кол-во запрашиваемых постов.~~
