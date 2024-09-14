# rest-notes

## Описание

REST API для управления заметками с возможностью создания и просмотра заметок.

## Требования

- Go 1.18+
- PostgreSQL

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/MaximInnopolis/rest-notes.git
cd rest-notes
```

2. Соберите докер-билд:
```bash
make up-all
```

3. Проведите миграцию:
```bash
make migrate
```

4. Выполняйте сетевые запросы с помощью Postman или вручную через консоль. Примеры запросов:

Регистрация пользователя:
```bash
curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{
           "username": "existinguser",
           "password": "password123"
         }'
```

Логин пользователя:
```bash
curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{
           "username": "existinguser",
           "password": "password123"
         }'

```

Создание заметки:
```bash
curl -X POST http://localhost:8080/notes/new \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <your-token-here>" \
     -d '{
           "title": "Meeting Notes",
           "description": "Notes from the meeting",
           "due_date": "2024-09-30T10:00:00Z"
         }'
```

Получение списка заметок:
```bash
curl -X GET http://localhost:8080/notes/list \
     -H "Authorization: Bearer <your-token-here>"
```