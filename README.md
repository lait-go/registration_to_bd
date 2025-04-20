REST-сервис на Go, который принимает ФИО через API, обогащает данные из открытых внешних API (возраст, пол, национальность), сохраняет в PostgreSQL и предоставляет интерфейс для управления данными (получение, удаление, обновление, фильтрация и пагинация).

---

## 🚀 Возможности

- 🔎 Получение информации о людях с фильтрацией и пагинацией
- ➕ Добавление нового человека с автоматическим обогащением
- ♻️ Обновление информации по идентификатору
- ❌ Удаление человека по идентификатору
- 🧠 Обогащение:
  - **Возраст** — [`agify.io`](https://api.agify.io/?name=...)
  - **Пол** — [`genderize.io`](https://api.genderize.io/?name=...)
  - **Национальность** — [`nationalize.io`](https://api.nationalize.io/?name=...)
- 🐘 Хранение в PostgreSQL (через `sqlx`)
- 🧪 Поддержка логгирования:
  - Общие логи в `log_storage.json`
  - Ошибки в `error_log.json`
- ⚙️ Конфигурация через `.env`
- 📄 Swagger-документация

---

## 📦 Пример JSON-запроса на добавление

```json
POST

{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"
}
```

📤 Примеры запросов через curl
➕ Добавить 
```
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"name": "Dmitriy", "surname": "Ushakov", "patronymic": "Vasilevich"}'
```

🔍 Получить всех (с фильтрацией и пагинацией)
```
curl "http://localhost:8080?name=Dmitriy&limit=5&offset=0"
```

🔍 Получить по ID
```
curl http://localhost:8080/3
```

♻️ Обновить по ID
```
curl -X PUT http://localhost:8080/3 \
  -H "Content-Type: application/json" \
  -d '{"name": "Ivan", "surname": "Petrov", "patronymic": "Nikolaevich"}'
```

❌ Удалить по ID
```
curl -X DELETE http://localhost:8080/3
```

🗂 Структура проекта
.
├── cmd/                    # Точка входа в приложение
├── config/                 # Загрузка конфигурации из .env
├── db/                     # Подключение и миграции PostgreSQL
├── internal/
│   ├── api/                # Обработчики REST-запросов
│   ├── err/                # Обработка ошибок
│   └── utils/
│       └── log/            # JSON-логгирование
├── migrations/             # SQL-файлы (CREATE, INSERT, UPDATE, DELETE)
├── swagger/                # Swagger-документация
├── .env                    # Конфигурация проекта
├── README.md               # Документация проекта

🗄 Пример .env файла
```
ENV=local
HOST=:8080
DATE_CONF=postgres://name:password@hash:port/mydb?sslmode=disable
LOG_PATH=../internal/utils/log/log_storage.json
ERROR_PATH=../internal/err/error_log.json
```

🗃 Структура БД (PostgreSQL)
Таблица создаётся через SQL миграции:
```
CREATE TABLE IF NOT EXISTS person (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    age INT CHECK (age > 0 AND age <= 100),
    gender TEXT CHECK (gender IN ('male', 'female')),
    nationality TEXT
);
```
Также должны быть миграции для:

Вставки (insert_to_tables.sql)

Получения по ID (receiving_by_id.sql)

Обновления (update_by_id.sql)

Удаления (delete_by_id.sql)



🛠 Запуск проекта
go run cmd/main.go
Убедитесь, что PostgreSQL запущен и доступен, а таблицы созданы миграциями.


