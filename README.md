# Marketplace REST API (Go)

## Технологии

- Go (net/http, стандартная библиотека)
- GORM ORM (https://gorm.io/)
- **PostgreSQL** - основная БД (используется в Docker)
- **SQLite** - in-memory БД для тестов (без внешних зависимостей)
- JWT для авторизации

## Запуск

### 1. Клонировать репозиторий

```sh
git clone https://github.com/Yvetlin/MarketVK.git
cd MarketVK
```

### 2. Установить зависимости Go

```sh
go mod download
```

### 3. Запуск через Docker

```sh
docker-compose up --build
```
По умолчанию будет поднят сервер Go (/cmd/main.go) и контейнер с PostgreSQL.

### 4. Перемененые окружения
Создать файл .env или задать переменную DB_DSN для подключения к Postgres, например:
```env
DB_DSN=host=db user=postgres password=postgres dbname=marketplace port=5432 sslmode=disable
```
В docker-compose переменные уже прописаны.

## Запуск тестов

Для тестов используется SQLite in-memory.
Должен быть установлен GCC (MinGW/TDN-GCC для Windows), и CGO должен быть включен.
```sh 
#Windows
set CGO_ENABLED=1
go test ./...
```
```sh 
#Linux/macOS
CGO_ENABLED=1
go test ./...
```
Тесты не трогают основную БД и запускаются независимо

## Основные команды

- POST /register - регистрация пользователя
- POST /login - вход, возврат JWT
- POST /ads - создать объявление (авторизация)
- GET /ads - список объявлений (с фильтрами и пагинацией)
- PUT /ads - обновить объявление (только автор)
- DELETE /ads - удалить объявление (только автор)

## API
- POST /register - регистрация пользователя
```raw 
{
  "login": "Имя пользователя (более 4х символов)",
  "password": "Пароль (более 6-ти символов)"
  }
```
- POST /login - вход, возврат JWT
```raw
{
  "login": "Имя пользователя",
  "password": "Пароль"
}
```
- POST /ads - создать объявление (авторизация)
```raw 
{
  "title": "Название",
  "description": "Описание",
  "image_url": "Ссылка на фото",
  "price": Стоимость (int)
}
```
```Headers
Authorization - Bearer <token>
```
- GET /ads - список объявлений (с фильтрами и пагинацией)
```Headers
Authorization - Bearer <token>
```
- PUT /ads - обновить объявление (только автор)
```raw 
{
    "id": 1,
    "title": "Новое название",
    "description": "Новое описание",
    "image_url": "Новое изображение",
    "price": 324234
}
```
```Headers
Authorization - Bearer <token>
```
- DELETE /ads - удалить объявление (только автор)
```raw 
{
  "id": Айди удаляемого объявления
}

```
```Headers
Authorization - Bearer <token>
```

## Примечания
- Все пароли хранятся в базе только в захешированном виде (bcrypt)
- В тестах - изолированная база на SQLite, не нужен отдельный сервер
- Для запуска тестов на Windows обязательно наличие GCC
