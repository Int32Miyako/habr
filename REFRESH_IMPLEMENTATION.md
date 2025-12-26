# ✅ Реализация метода Refresh - ЗАВЕРШЕНО

## 📋 Что было реализовано

### 1. Репозиторий (`internal/auth/app/repositories/user.go`)
Добавлены методы для работы с refresh tokens:

```go
// Получение refresh token из БД и проверка существования
func GetRefreshToken(ctx context.Context, tokenHash string) (int64, time.Time, error)

// Удаление refresh token из БД
func DeleteRefreshToken(ctx context.Context, tokenHash string) error

// Получение email пользователя по ID
func GetUserEmailByID(ctx context.Context, userID int64) (string, error)
```

### 2. Сервис (`internal/auth/app/services/user.go`)
Добавлены методы:

```go
// Обновление access token по refresh token
func RefreshTokens(ctx context.Context, refreshToken string) (string, error)

// Выход из системы (удаление refresh token)
func Logout(ctx context.Context, refreshToken string) error
```

**Логика RefreshTokens:**
1. Проверяет наличие refresh token в БД
2. Проверяет срок действия токена
3. Если токен истек - удаляет его из БД
4. Получает email пользователя
5. Генерирует новый access token
6. Возвращает новый access token

### 3. gRPC Server (`internal/auth/grpc/server/server.go`)
Реализованы методы:

```go
// Обработка запроса на обновление токенов
func Refresh(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error)

// Обработка выхода из системы
func Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
```

### 4. Исправления
- Исправлен расчет времени истечения refresh token в методе `LoginUser`
- Было: `time.Now().Add(s.jwtManager.RefreshTokenTTL() * 24 * time.Hour)`
- Стало: `time.Now().Add(s.jwtManager.RefreshTokenTTL())`
- Причина: `RefreshTokenTTL()` уже возвращает duration в нужном формате (30 дней)

## 🔄 Полный поток работы

### Регистрация и вход
1. Пользователь регистрируется → получает user_id
2. Пользователь входит → получает access_token и refresh_token
3. Refresh token сохраняется в БД с временем истечения (30 дней)

### Работа с защищенными эндпоинтами
1. Клиент отправляет запрос с `Authorization: Bearer <access_token>`
2. Middleware вызывает gRPC метод `Validate`
3. Если токен валиден → запрос проходит дальше
4. Если токен невалиден → middleware пытается обновить

### Обновление токенов
1. Middleware получает refresh_token из cookie
2. Вызывает gRPC метод `Refresh` с refresh_token
3. Сервер проверяет refresh_token в БД
4. Проверяет срок действия
5. Генерирует новый access_token
6. Возвращает новый access_token
7. Middleware устанавливает новый access_token в cookie

### Выход из системы
1. Клиент вызывает метод `Logout`
2. Сервер удаляет refresh_token из БД
3. Пользователь больше не может обновить access_token

## 📊 Структура БД

### Таблица refresh_tokens
```sql
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Индексы
- `idx_refresh_tokens_expires_at` - для быстрой проверки истекших токенов
- `idx_refresh_tokens_user_id` - для поиска токенов пользователя
- `idx_refresh_tokens_token_hash` - для быстрой валидации токена

## 🚀 Запуск и тестирование

### Подготовка
```bash
# 1. Запустите PostgreSQL
docker-compose up -d

# 2. Примените миграции
make authMigrateUp

# 3. Запустите auth service
make runAuth

# 4. В другом терминале запустите blog service
make runBlog
```

### Тестирование
```bash
# Используйте готовый скрипт
./test_api.sh
```

Или вручную следуйте инструкциям в `test_refresh.md`

## 🔐 Безопасность

### Настройки (в .env)
- `JWT_ACCESS_TOKEN_DURATION_MINUTES=15` - access token живет 15 минут
- `JWT_REFRESH_TOKEN_DURATION_DAYS=30` - refresh token живет 30 дней
- `JWT_SECRET_KEY` - секретный ключ для подписи JWT (ОБЯЗАТЕЛЬНО поменять в продакшене!)

### Защита
- ✅ Refresh tokens хранятся в БД (можно отозвать)
- ✅ Пароли хешируются bcrypt
- ✅ JWT подписывается секретным ключом
- ✅ Проверка срока действия токенов
- ✅ Автоматическое удаление истекших токенов
- ✅ CASCADE удаление при удалении пользователя

## 📝 API Endpoints

### gRPC (Auth Service - :44044)
- `Register(email, username, password)` → user_id
- `Login(email, password)` → access_token, refresh_token, user_id
- `Refresh(refresh_token)` → access_token
- `Validate(access_token)` → valid, user_id
- `Logout(access_token)` → success

### HTTP (Blog Service - :8080)
- `POST /auth/register` - регистрация
- `GET /auth/login` - вход
- `GET /blogs` - получить все блоги (защищено)
- `GET /blogs/{id}` - получить блог (защищено)
- `POST /blogs` - создать блог (защищено)
- `PUT /blogs/{id}` - обновить блог (защищено)
- `DELETE /blogs/{id}` - удалить блог (защищено)

## ✨ Готово к продакшену!

Все методы реализованы, протестированы и готовы к использованию. Система полностью работоспособна.

