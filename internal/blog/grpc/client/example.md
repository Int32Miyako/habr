# Пример использования AuthClient

## Инициализация клиента

```go
import (
    "context"
    "habr/internal/blog/grpc/client"
    "log"
)

// Создание клиента
authClient, err := client.NewAuthClient("localhost:50051")
if err != nil {
    log.Fatalf("Failed to create auth client: %v", err)
}
defer authClient.Close() // Обязательно закрыть соединение при завершении

ctx := context.Background()
```

## Регистрация пользователя

```go
registerResp, err := authClient.Register(
    ctx,
    "user@example.com",
    "username",
    "password123",
)
if err != nil {
    log.Printf("Registration failed: %v", err)
    return
}
log.Printf("User registered with ID: %d", registerResp.UserId)
```

## Вход пользователя

```go
loginResp, err := authClient.Login(
    ctx,
    "user@example.com",
    "password123",
)
if err != nil {
    log.Printf("Login failed: %v", err)
    return
}
log.Printf("Login successful, token: %s", loginResp.Token)
```

## Использование в HTTP хендлерах

```go
func LoginHandler(authClient *client.AuthClient) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        resp, err := authClient.Login(r.Context(), req.Email, req.Password)
        if err != nil {
            http.Error(w, "Authentication failed", http.StatusUnauthorized)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{
            "token": resp.Token,
        })
    }
}
```

## Конфигурация

Добавьте в `.env` файл:

```env
AUTH_GRPC_ADDRESS=localhost:50051
```

## Интеграция в приложение

Клиент уже интегрирован в `cmd/blog/main.go` и доступен через router. Вы можете использовать его в middleware для проверки аутентификации или в хендлерах для вызова методов auth сервиса.

