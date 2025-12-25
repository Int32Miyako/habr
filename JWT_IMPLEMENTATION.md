# JWT Authentication Implementation Guide

## 📋 Обзор

Реализация JWT аутентификации с Access и Refresh токенами для микросервиса авторизации.

## 🔑 Концепция

### Access Token
- **Назначение**: Доступ к защищенным ресурсам
- **Срок жизни**: 15 минут (настраивается)
- **Формат**: JWT (JSON Web Token)
- **Передача**: В заголовке `Authorization: Bearer <token>`

### Refresh Token
- **Назначение**: Обновление Access Token
- **Срок жизни**: 30 дней (настраивается)
- **Формат**: Случайная строка (32 байта в base64)
- **Хранение**: В БД в виде SHA-256 хэша
- **Передача**: Как обычная строка (в gRPC) или через HTTP-Only Cookie (в HTTP API)

## 🏗️ Архитектура

```
Client → gRPC Server → UserService → UserRepository → PostgreSQL
                  ↓
              JWTManager
```

## 📦 Компоненты

### 1. JWT Manager (`internal/auth/core/jwt/jwt_manager.go`)
- Генерация Access Token (JWT)
- Генерация Refresh Token (random)
- Валидация Access Token

### 2. User Repository (`internal/auth/app/repositories/user.go`)
- Работа с пользователями
- Сохранение/проверка/удаление Refresh Token в БД

### 3. User Service (`internal/auth/app/services/user.go`)
- Login - аутентификация и выдача токенов
- Register - регистрация с автоматической выдачей токенов
- Refresh - обновление пары токенов
- Logout - удаление Refresh Token
- ValidateAccessToken - проверка Access Token

### 4. gRPC Server (`internal/auth/grpc/server/server.go`)
- Обработка gRPC запросов
- Валидация входных данных
- Логирование

## 🔄 Потоки данных

### Регистрация
```
Client → Register(email, username, password)
Server → Хэширует пароль
Server → Создает пользователя в БД
Server → Генерирует Access Token (JWT)
Server → Генерирует Refresh Token (random)
Server → Сохраняет хэш Refresh Token в БД
Server → Возвращает: user_id, access_token, refresh_token
```

### Вход
```
Client → Login(email, password)
Server → Проверяет пользователя в БД
Server → Проверяет пароль (bcrypt)
Server → Генерирует Access Token
Server → Генерирует Refresh Token
Server → Сохраняет хэш Refresh Token в БД
Server → Возвращает: user_id, access_token, refresh_token
```

### Обновление токенов
```
Client → Refresh(refresh_token)
Server → Проверяет Refresh Token в БД
Server → Проверяет срок действия
Server → Генерирует новый Access Token
Server → Генерирует новый Refresh Token
Server → Удаляет старый Refresh Token из БД
Server → Сохраняет новый Refresh Token в БД
Server → Возвращает: access_token, refresh_token
```

### Валидация
```
Client → Validate(access_token)
Server → Парсит JWT
Server → Проверяет подпись
Server → Проверяет срок действия
Server → Возвращает: user_id, email
```

### Выход
```
Client → Logout(refresh_token)
Server → Удаляет Refresh Token из БД
Server → Возвращает: success
```

## 🚀 Установка

### 1. Установите зависимости
```bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
```

### 2. Настройте .env
```env
JWT_SECRET_KEY=your-super-secret-key-here
JWT_ACCESS_TOKEN_DURATION_MINUTES=15
JWT_REFRESH_TOKEN_DURATION_DAYS=30
```

### 3. Запустите миграции
```bash
make migrate-auth-up
```

### 4. Сгенерируйте protobuf файлы
```bash
make proto
```

### 5. Запустите сервис
```bash
make run-auth
```

## 🧪 Тестирование

### Регистрация
```bash
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "username": "testuser",
  "password": "password123"
}' localhost:50051 auth.Auth/Register
```

### Вход
```bash
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "password123"
}' localhost:50051 auth.Auth/Login
```

### Обновление токенов
```bash
grpcurl -plaintext -d '{
  "refresh_token": "YOUR_REFRESH_TOKEN_HERE"
}' localhost:50051 auth.Auth/Refresh
```

### Валидация
```bash
grpcurl -plaintext -d '{
  "access_token": "YOUR_ACCESS_TOKEN_HERE"
}' localhost:50051 auth.Auth/Validate
```

### Выход
```bash
grpcurl -plaintext -d '{
  "refresh_token": "YOUR_REFRESH_TOKEN_HERE"
}' localhost:50051 auth.Auth/Logout
```

## 🛡️ Безопасность

### Реализовано:
✅ Пароли хэшируются с помощью bcrypt
✅ Refresh Token хранится в БД как SHA-256 хэш
✅ Access Token имеет короткий срок жизни
✅ JWT подписан секретным ключом
✅ Валидация всех входных данных

### Рекомендации:
⚠️ Используйте сложный JWT_SECRET_KEY в продакшене
⚠️ Используйте HTTPS для всех запросов
⚠️ В HTTP API храните Refresh Token в HTTP-Only Cookie
⚠️ Настройте rate limiting для предотвращения брутфорса
⚠️ Добавьте логирование подозрительной активности

## 📊 База данных

### Таблица users
```sql
id SERIAL PRIMARY KEY
username TEXT UNIQUE NOT NULL
password_hash TEXT NOT NULL
email TEXT UNIQUE NOT NULL
created_at TIMESTAMP
updated_at TIMESTAMP
```

### Таблица refresh_tokens
```sql
id SERIAL PRIMARY KEY
user_id INTEGER REFERENCES users(id)
token_hash VARCHAR(255) UNIQUE NOT NULL  -- SHA-256 хэш токена
expires_at TIMESTAMP NOT NULL
created_at TIMESTAMP
```

## 🔍 Типичные сценарии

### Клиент делает запрос к API
```
1. Клиент отправляет Access Token в заголовке
2. Сервер валидирует токен через Validate()
3. Если токен валиден → обрабатывает запрос
4. Если токен истек → возвращает 401
5. Клиент вызывает Refresh()
6. Клиент повторяет запрос с новым Access Token
```

### Истек Refresh Token
```
1. Клиент пытается вызвать Refresh()
2. Сервер возвращает ошибку "invalid or expired refresh token"
3. Клиент перенаправляет пользователя на страницу входа
4. Пользователь вводит email и пароль
5. Клиент получает новую пару токенов
```

## 📝 Примечания

- **Rotation**: При каждом Refresh старый токен удаляется и создается новый (Token Rotation)
- **Cleanup**: Можно добавить cronjob для очистки истекших токенов: `CleanupExpiredTokens()`
- **Multiple devices**: Каждое устройство получает свой Refresh Token
- **Logout from all devices**: Используйте `DeleteAllUserRefreshTokens(userID)`

## 🎯 Что дальше?

1. Добавить HTTP middleware для автоматической валидации токенов
2. Реализовать rate limiting
3. Добавить email верификацию
4. Реализовать "Remember me" функционал
5. Добавить 2FA (Two-Factor Authentication)
6. Логирование всех действий с токенами

