# Тестовое задание Medods
## Краткое описание 
RESTful API-сервис на Go (Gin) для авторизации пользователей, генерации JWT-токенов, и отправки уведомлений через Webhook. Использует PostgreSQL в качестве хранилища.
## Технологический стек
- Go (GIN и GORM)
- PostgreSQL
- Docker и Docker Compose
- Swagger

## Структура проекта
```
medods_test_task/
├── cmd/
│   └── main.go                   # Точка входа
├── internal/
│   ├── config/                   # Загрузка конфигурации
│   ├── db/
│   │   ├── impl/                 # Реализация взаимодействия с БД
│   │   └── intf/                 # Интерфейс работы с БД
│   ├── dto/                      # Структуры запросов/ответов
│   ├── handler/                  # HTTP-обработчик 
│   ├── middleware/               # Middleware для Gin
│   ├── model/                    # Бизнес-модель
│   ├── repository/
│   │   ├── impl/                 # Реализация репозитория
│   │   └── intf/                 # Интерфейс репозитория
│   ├── service/
│   │   ├── impl/                 # Реализация бизнес-логики
│   │   └── intf/                 # Интерфейс сервиса
│   └── utils/                    # Вспомогательные функции
├── scripts/                      # SQL-инициализация и healthchecks
├── .env                          # Переменные окружения
├── docker-compose.yml            # Docker Compose
├── Dockerfile                    # Dockerfile для сборки Go-приложения
├── go.mod                        # Go-модули
├── go.sum
└── README.md
```
## Переменные окружения
```
ACCESS_TOKEN_TTL=15m                                                # Время жизни access токена
JWT_SECRET=jwt-secret                                               # JWT-secret
WEBHOOK=https://webhook.site/08de5a48-8337-436f-8410-4bc4d94b440f   # Ссылка на WebHook

POSTGRES_USER=user                                                  # Пользователь БД
POSTGRES_PASSWORD=password                                          # Пароль БД
POSTGRES_DB=db                                                      # Имя БД
POSTGRES_PORT=5432                                                  # Порт БД
```

## Запуск проекта
- Убедись, что установлен Docker
- Клонируй репозиторий
```
git clone https://github.com/umarov-work/medods_test_task.git
cd medods_test_task
```
- Собери и запусти контейнеры
```
docker compose -f docker-compose.yml up -d 
```

## Документация
Swagger-документация будет доступна по адресу:
```
http://localhost:8080/swagger/index.html
```