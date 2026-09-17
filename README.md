# Авито.Кухня - MVP Backend Service
## Описание работы 
Я уже реализовывал похожий проект предоставляющий API для определенного функционала на REST API, поэтому решил адаптировать существующие наработки под это тестовое
Использовал ИИ для написания кода в некоторых местах, для корректировки разработанной архитектуры, создания рутинных спецификаций типа OpenAPI.

## Описание проекта

**Авито.Кухня** - это MVP веб-сервис для заказа и доставки еды, который предоставляет API для:
- **Ресторанов/кафе** - управление меню, прием и обработка заказов
- **Пользователей** - просмотр заведений, оформление заказов, отслеживание статуса

---

## Архитектура

Система состоит из двух основных компонентов:

### 1. Main Service (Основной API сервис)
- Управление заведениями и меню
- Обработка заказов пользователей
- Интеграция с ресторанами через API
- PostgreSQL для хранения данных

### 2. Restaurant Service (Сервис-пример заведения)
- Демонстрация интеграции с основным API
- Автоматическая обработка заказов
- Управление меню через API Main Service
- Webhook для получения уведомлений о заказах

```
┌──────────────┐         HTTP API        ┌──────────────────┐
│              │◄───────────────────────►│                  │
│   Client     │                         │  Main Service    │
│   (User)     │                         │  (REST API)      │
│              │                         │                  │
└──────────────┘                         └────────┬─────────┘
                                                  │
                                         ┌────────▼─────────┐
                                         │   PostgreSQL     │
                                         │    Database      │
                                         └────────▲─────────┘
                                                  │
                      HTTP API + Webhooks         │
                           ┌─────────────────────┘
                           │
                  ┌────────▼─────────┐
                  │   Restaurant     │
                  │    Service       │
                  │  (Integration)   │
                  └──────────────────┘
```

---

## Основные бизнес-сценарии

### Сценарий пользователя:
1. Просмотр списка доступных заведений
2. Выбор заведения и просмотр меню
3. Добавление блюд в корзину
4. Оформление заказа с указанием адреса доставки
5. Отслеживание статуса заказа (создан → подтверждён → готовится → доставляется → выполнен)

### Сценарий ресторана:
1. Регистрация заведения в системе
2. Управление меню (добавление/обновление/удаление блюд)
3. Получение уведомлений о новых заказах
4. Обновление статуса заказа
5. Управление доступностью блюд (в наличии/нет в наличии)

---

## Схема базы данных

### Основные таблицы:

**restaurants** - Информация о заведениях
```
id              SERIAL PRIMARY KEY
name            VARCHAR(255) NOT NULL
description     TEXT
address         VARCHAR(500)
phone           VARCHAR(20)
is_active       BOOLEAN DEFAULT true
webhook_url     VARCHAR(500)
created_at      TIMESTAMP DEFAULT NOW()
updated_at      TIMESTAMP DEFAULT NOW()
```

**dishes** - Меню блюд
```
id              SERIAL PRIMARY KEY
restaurant_id   INTEGER REFERENCES restaurants(id)
name            VARCHAR(255) NOT NULL
description     TEXT
price           DECIMAL(10,2) NOT NULL
category        VARCHAR(100)
is_available    BOOLEAN DEFAULT true
created_at      TIMESTAMP DEFAULT NOW()
updated_at      TIMESTAMP DEFAULT NOW()
```

**orders** - Заказы пользователей
```
id              SERIAL PRIMARY KEY
user_id         INTEGER NOT NULL
restaurant_id   INTEGER REFERENCES restaurants(id)
status          VARCHAR(50) NOT NULL
total_price     DECIMAL(10,2) NOT NULL
delivery_address TEXT NOT NULL
phone           VARCHAR(20)
created_at      TIMESTAMP DEFAULT NOW()
updated_at      TIMESTAMP DEFAULT NOW()
```

**order_items** - Состав заказа
```
id              SERIAL PRIMARY KEY
order_id        INTEGER REFERENCES orders(id)
dish_id         INTEGER REFERENCES dishes(id)
quantity        INTEGER NOT NULL
price           DECIMAL(10,2) NOT NULL
```

**users** - Пользователи (упрощенная версия для MVP)
```
id              SERIAL PRIMARY KEY
name            VARCHAR(255) NOT NULL
phone           VARCHAR(20) UNIQUE
address         TEXT
created_at      TIMESTAMP DEFAULT NOW()
```

---

## Технологический стек

- **Язык:** Go 1.21+
- **Веб-фреймворк:** Gorilla Mux
- **База данных:** PostgreSQL 15
- **ORM:** GORM
- **Миграции:** golang-migrate
- **Контейнеризация:** Docker, Docker Compose
- **Линтинг:** golangci-lint
- **Документация API:** OpenAPI 3.0
- **Диаграммы:** PlantUML

---

## Быстрый старт

### Предварительные требования
- Docker и Docker Compose
- Git

### Запуск проекта

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd backend-trainee-assignment-autumn-2026-flow-2-sondernce-22088573
```

2. Запустите все сервисы через Docker Compose:
```bash
docker-compose up --build
```

3. Дождитесь инициализации (миграции применятся автоматически)

4. API будет доступно:
   - Main Service: `http://localhost:8080`
   - Restaurant Service: `http://localhost:8081`

### Проверка работоспособности

```bash
# Получить список заведений
curl http://localhost:8080/api/v1/restaurants

# Получить меню заведения
curl http://localhost:8080/api/v1/restaurants/1/menu

# Создать заказ
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "restaurant_id": 1,
    "delivery_address": "ул. Пушкина, д. 10",
    "phone": "+79991234567",
    "items": [
      {"dish_id": 1, "quantity": 2},
      {"dish_id": 2, "quantity": 1}
    ]
  }'
```

---

## Документация

### API Документация
OpenAPI спецификация доступна в файле [`api/openapi.yaml`](./api/openapi.yaml)

### Диаграммы
- [Customer Journey Map - Пользователь](./docs/diagrams/cjm-user.puml)
- [Customer Journey Map - Ресторан](./docs/diagrams/cjm-restaurant.puml)
- [C4 Architecture Diagram](./docs/diagrams/c4-architecture.puml)
- [ER-диаграмма БД](./docs/diagrams/er-diagram.puml)

---

## Структура проекта

```
.
├── README.md
├── docker-compose.yml
├── .gitignore
│
├── api/                          # OpenAPI спецификации
│   └── openapi.yaml
│
├── docs/                         # Документация
│   │   │   └── diagrams/
│       ├── cjm-user.puml
│       ├── cjm-restaurant.puml
│       ├── c4-architecture.puml
│       └── er-diagram.puml
│
├── main-service/                 # Основной API сервис
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   │   │
│   ├── cmd/
│   │   └── main.go              # Entry point
│   │
│   ├── internal/
│   │   ├── handlers/            # HTTP handlers (аналог вашего проекта)
│   │   │   ├── restaurant.go
│   │   │   ├── dish.go
│   │   │   ├── order.go
│   │   │   ├── dto.go
│   │   │   └── server.go
│   │   │
│   │   ├── models/              # Domain модели
│   │   │   ├── restaurant.go
│   │   │   ├── dish.go
│   │   │   ├── order.go
│   │   │   └── user.go
│   │   │
│   │   ├── repository/          # Работа с БД
│   │   │   ├── restaurant.go
│   │   │   ├── dish.go
│   │   │   └── order.go
│   │   │
│   │   └── service/             # Бизнес-логика
│   │       ├── restaurant.go
│   │       ├── dish.go
│   │       └── order.go
│   │
│   └── migrations/              # SQL миграции
│       ├── 000001_init.up.sql
│       └── 000001_init.down.sql
│
└── restaurant-service/          # Сервис-пример ресторана
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    │
    ├── cmd/
    │   └── main.go
    │
    └── internal/
        ├── handlers/            # Webhook handlers
        ├── client/              # HTTP client для Main API
        └── models/
```

---

## Разработка

### Запуск линтера
```bash
cd main-service
golangci-lint run
```

### Применение миграций вручную
```bash
migrate -path main-service/migrations -database "postgresql://user:password@localhost:5432/avito_kitchen?sslmode=disable" up
```

### Запуск тестов
```bash
cd main-service
go test ./...
```

---
