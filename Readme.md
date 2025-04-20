# Workmate Task API

Асинхронный сервис для создания и отслеживания долгих I/O-bound задач на Go с использованием PostgreSQL.
Позволяет создавать задачи через HTTP API и получать их статус и результат по уникальному идентификатору.

---

## Технологии

- Go
- PostgreSQL
- sqlx
- chi (REST API)
- zap (логирование)
- docker-compose
- golang-migrate (миграции)
- testify (unit-тесты)
- godotenv (конфиги)

---

## Инструкция по запуску

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/Egorpalan/workmate-test.git
cd workmate-test
```


### 2. Настройте переменные окружения

Создайте `.env` на основе `.env.example` (или проверьте значения по умолчанию):

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tasks_db
DB_SSLMODE=disable
SERVER_PORT=8080
```


### 3. Запустите сервис и базу данных

```bash
make compose-up
```


### 4. Примените миграции

В новом терминале:

```bash
make migrate-up
```


### 5. Запустите тесты

```bash
make test
```

---

## API

### 1. Создать задачу

**POST** `/api/tasks`

- **Описание:** Создаёт новую долгую задачу.
- **Ответ:**

```json
{
  "id": "c9e8b5c7-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "status": "pending",
  "result": {},
  "error": "",
  "created_at": "2025-04-20T19:00:00Z",
  "updated_at": "2025-04-20T19:00:00Z"
}
```

---

### 2. Получить задачу по ID

**GET** `/api/tasks/{id}`

- **Описание:** Получает статус и результат задачи по её идентификатору.
- **Ответ:**

```json
{
  "id": "c9e8b5c7-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "status": "completed",
  "result": {
    "message": "Task completed successfully",
    "timestamp": "2025-04-20T19:03:00Z"
  },
  "error": "",
  "created_at": "2025-04-20T19:00:00Z",
  "updated_at": "2025-04-20T19:03:00Z"
}
```

---

### 3. Получить список задач

**GET** `/api/tasks?limit=10&offset=0`

- **Описание:** Получает список задач с пагинацией.
- **Ответ:** Массив задач.

---

## Примеры запросов

### Создать задачу

```bash
curl -X POST http://localhost:8080/api/tasks
```


### Получить задачу по ID

```bash
curl http://localhost:8080/api/tasks/<task_id>
```


### Получить список задач

```bash
curl http://localhost:8080/api/tasks
```

---

## Особенности

- **Асинхронные задачи:** задачи выполняются в фоне, статус можно отслеживать по ID.
- **REST API:** простые и понятные эндпоинты.
- **Логирование:** все события и ошибки логируются через zap.
- **Graceful shutdown:** сервис корректно завершает работу по SIGINT/SIGTERM.
- **Тесты:** покрытие бизнес-логики unit-тестами (testify).
- **Миграции:** структура БД управляется через golang-migrate.
- **Контейнеризация:** всё запускается через docker-compose.

---

## Makefile команды

- `make run` — локальный запуск приложения
- `make build` — сборка бинарника
- `make compose-up` — запуск docker-compose (приложение + БД)
- `make compose-down` — остановка docker-compose
- `make migrate-up` — применить миграции к БД
- `make migrate-down` — откатить миграции
- `make test` — запустить unit-тесты



