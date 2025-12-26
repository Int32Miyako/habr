# Тестирование Refresh Token

## Шаги для тестирования:

### 1. Запустите сервисы

```bash
# Запустите базу данных
docker-compose up -d

# Запустите Auth сервис
make run-auth

# В другом терминале запустите Blog сервис
make run-blog
```

### 2. Зарегистрируйте пользователя

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }'
```

### 3. Войдите в систему

```bash
curl -X GET http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Сохраните `access_token` и `refresh_token` из ответа.

### 4. Используйте access token для доступа к защищенному эндпоинту

```bash
curl -X GET http://localhost:8080/blogs \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. Проверьте Refresh (если access token истек)

Для тестирования refresh механизма в middleware:
- Подождите 15 минут (или измените JWT_ACCESS_TOKEN_DURATION_MINUTES в .env на 1 минуту)
- Добавьте cookie с refresh_token
- Попробуйте снова обратиться к защищенному эндпоинту

```bash
curl -X GET http://localhost:8080/blogs \
  -H "Authorization: Bearer EXPIRED_ACCESS_TOKEN" \
  -b "refresh_token=YOUR_REFRESH_TOKEN"
```

Middleware автоматически обновит access token.

## Что реализовано:

### В Auth Service (gRPC):
- ✅ `Register` - регистрация пользователя
- ✅ `Login` - вход и выдача access + refresh токенов
- ✅ `Refresh` - обновление access token по refresh token
- ✅ `Validate` - валидация access token
- ✅ `Logout` - выход из системы

### В Blog Service (HTTP):
- ✅ Auth middleware с поддержкой автоматического refresh
- ✅ Защищенные эндпоинты для работы с блогами

### База данных:
- ✅ Таблица refresh_tokens с индексами
- ✅ Связь с пользователями через foreign key
- ✅ Автоматическое удаление при удалении пользователя (CASCADE)

## Логика работы Refresh:

1. Клиент отправляет запрос с access token в заголовке Authorization
2. Middleware проверяет валидность access token через gRPC Validate
3. Если токен невалидный:
   - Middleware ищет refresh_token в cookie
   - Отправляет запрос на gRPC Refresh
   - Получает новый access_token
   - Устанавливает новый access_token в cookie
   - Пропускает запрос дальше
4. Если токен валидный - сразу пропускает запрос

## Безопасность:

- Access token живет 15 минут (настраивается в .env)
- Refresh token живет 30 дней (настраивается в .env)
- Refresh token хранится в базе данных
- При refresh проверяется срок действия токена
- Пароли хешируются с помощью bcrypt
- JWT подписывается секретным ключом

