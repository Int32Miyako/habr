CREATE DATABASE auth_db;
CREATE DATABASE blog_db;
CREATE DATABASE notification_db;

-- Создание пользователя для notification_db
CREATE USER notification_user WITH PASSWORD 'notification_pass';
GRANT ALL PRIVILEGES ON DATABASE notification_db TO notification_user;
