# API платежной системы

## Содержание
- [Обзор](#обзор)
- [Запуск](#запуск)
- [Структура проекта](#структура)
- [Эндпоинты](#эндпоинты)
- [Используемые библиотеки](#библиотеки)


## Обзор
  API платежной системы представляет собой RESTful сервис, позволяющий управлять транзакциями и пользовательскими данными. Поддерживаются функции получения информации о пользователях, просмотра истории транзакций, пополнения баланса и перевода средств.
>[!CAUTION]
>Данные о деньгах хранятся в целочисленом int формате для безопасности, так как float не точен.
>Так же они хранятся не в рублях или долларах, а в копейках или центах. То есть при выводить их на фронт, нужно разделить на 100


## Запуск
 > **Docker запуск**  
> ```sh
> make docker-build
> ```
>  
> **Локальный запуск**  
> Отредактируйте `.env` и выполните команду:  
> ```sh
> make local-run
> ```



## Структура
```
├── cmd
│   └── payment
│       └── main.go #точка входа
├── config 
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── config 
│   │   ├── config.go
│   │   └── provider.go
│   ├── delivery
│   │   └── http
│   │       └── handler 
│   │           ├── provider.go
│   │           ├── transaction_handler.go
│   │           └── user_handler.go
│   ├── domain
│   │   ├── entity
│   │   │   ├── transaction.go
│   │   │   └── user.go
│   │   └── service
│   │       ├── provider.go
│   │       ├── transaction_service.go
│   │       └── user_service.go
│   ├── repository
│   │   ├── provider.go
│   │   ├── transaction_repository.go
│   │   └── user_repository.go
│   └── server
│       ├── server.go
│       ├── wire.go
│       └── wire_gen.go
├── logs
│   └── iq-testtask-error.log
├── migrations
│   ├── 20250124165437_create_user_table.sql
│   └── 20250124165452_create_transaction_table.sql
├── pkg
│   └── logger
│       └── logger.go
└── storage
    └── postgres
        ├── postgres.go
        └── provider.go
```

## Эндпоинты

### Пользовательские эндпоинты
#### Получение пользователя по ID
**Запрос:**
```
GET /api/user/:id
```

**Ответ:**
- `200 OK` – Возвращает информацию о пользователе
- `400 Bad Request` – Некорректный ID пользователя
- `404 Not Found` – Пользователь не найден
- `500 Internal Server Error` – Ошибка сервера

**Пример ответа:**
```json
{
  "id": 1,
  "name": "Иван Иванов",
  "balance": 1000
}
```

---

### Эндпоинты транзакций
#### Получение последних транзакций
**Запрос:**
```
GET /api/transaction/:user_id?limit={limit}
```

**Параметры запроса:**
- `limit` (необязательно) – Количество транзакций (по умолчанию: 10)

**Ответ:**
- `200 OK` – Возвращает массив транзакций
- `400 Bad Request` – Некорректный ID пользователя или limit
- `500 Internal Server Error` – Ошибка сервера

**Пример ответа:**
```json
[
  {
    "id": 101,
    "user_id": 1,
    "amount": 500,
    "type": "credit",
    "created_at": "2024-01-30T12:00:00Z"
  }
]
```

---

#### Пополнение баланса
**Запрос:**
```
POST /api/transaction/replenish/:user_id/:amount
```

**Ответ:**
- `200 OK` – Баланс успешно пополнен
- `400 Bad Request` – Некорректный ID пользователя или сумма
- `500 Internal Server Error` – Ошибка при пополнении

**Пример ответа:**
```json
{
  "status": "успешно пополнено"
}
```

---

#### Перевод средств
**Запрос:**
```
POST /api/transaction/transfer?from_user_id={id}&to_user_id={id}&amount={amount}
```

**Ответ:**
- `200 OK` – Перевод успешно выполнен
- `400 Bad Request` – Некорректные ID пользователей или сумма
- `500 Internal Server Error` – Ошибка при переводе

**Пример ответа:**
```json
{
  "status": "успешно переведено"
}
```

## Библиотеки

| Library      | Usage          |
| ------------ | -------------- |
| gin          | Base framework |
| pgx          | SQL library    |
| postgres     | Database       |
| custom logger| My Logger Setup| 
| viper        | Config library |
| wire    | Dependency injection|
