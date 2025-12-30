#!/bin/bash

# Скрипт для быстрого тестирования Refresh Token функциональности

echo "🚀 Запуск тестирования Refresh Token..."

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Регистрация пользователя
echo -e "\n${BLUE}1. Регистрация нового пользователя...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }')
echo -e "${GREEN}Response:${NC} $REGISTER_RESPONSE"

# 2. Вход в систему
echo -e "\n${BLUE}2. Вход в систему...${NC}"
LOGIN_RESPONSE=$(curl -s -X GET http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
echo -e "${GREEN}Response:${NC} $LOGIN_RESPONSE"

# Извлечение токенов
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)

echo -e "${YELLOW}Access Token:${NC} $ACCESS_TOKEN"
echo -e "${YELLOW}Refresh Token:${NC} $REFRESH_TOKEN"

# 3. Доступ к защищенному эндпоинту
echo -e "\n${BLUE}3. Доступ к защищенному эндпоинту (GET /blogs)...${NC}"
BLOGS_RESPONSE=$(curl -s -X GET http://localhost:8080/blogs \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo -e "${GREEN}Response:${NC} $BLOGS_RESPONSE"

# 4. Тест с невалидным токеном (для проверки refresh механизма)
echo -e "\n${BLUE}4. Тест с невалидным access token и refresh cookie...${NC}"
INVALID_TOKEN="invalid.token.here"
REFRESH_TEST=$(curl -s -X GET http://localhost:8080/blogs \
  -H "Authorization: Bearer $INVALID_TOKEN" \
  -b "refresh_token=$REFRESH_TOKEN" \
  -c /tmp/cookies.txt \
  -v 2>&1)
echo -e "${GREEN}Response:${NC} $REFRESH_TEST"

echo -e "\n${GREEN}✅ Тестирование завершено!${NC}"

