# **gRPC SSO Service**

🔐 **gRPC SSO Service** — это сервис единого входа (**SSO**) на основе **gRPC**, предназначенный для централизованной аутентификации и авторизации в распределённых системах. Основная цель — упростить управление доступом пользователей к различным сервисам внутри организации.

## **📌 Основные возможности**
- 📢 **Централизованное управление учетными записями**
- 🔑 **Генерация и валидация токенов доступа** (JWT)
- 🔗 **Поддержка внешних провайдеров аутентификации** (OAuth2, OpenID Connect)
- 🔒 **Защита передачи данных через TLS**
- ⚡ **Высокая производительность с gRPC**

---

## **🛠 Используемые технологии**
| Технология  | Описание  |
|------------|-----------|
| **Go** | Основной язык разработки |
| **gRPC** | Коммуникация между сервисами |
| **Protocol Buffers (Protobuf)** | Описание API и сериализация данных |
| **JWT (JSON Web Token)** | Аутентификация пользователей |
| **TLS** | Защищённая передача данных |
| **SQLite** | Лёгкая база данных для хранения пользователей |
| **Docker** | Контейнеризация сервиса |
| **YAML** | Конфигурация сервиса |

---

## **📥 Установка**
### **1️⃣ Клонирование репозитория**
```sh
git clone https://github.com/21Timofei/gRPC-SSO-Service.git
cd gRPC-SSO-Service
```

### **2️⃣ Установка зависимостей**
```sh
go mod tidy
```

### **3️⃣ Конфигурация сервиса**
Создайте `config.yaml` и добавьте:
```yaml
db:
  type: sqlite
  path: "./sso.db"

tls:
  enabled: true
  cert_file: "server.crt"
  key_file: "server.key"
```

---

## **🚀 Запуск**
### **Локальный запуск**
```sh
go run main.go
```
Сервис запустится на `localhost:50051`.

### **Запуск в Docker**
```sh
docker build -t grpc-sso .
docker run -p 50051:50051 --env-file .env grpc-sso
```

---

## **📡 gRPC API Методы**
| Метод  | Описание |
|--------|----------|
| **Login** | Аутентификация пользователя |
| **ValidateToken** | Проверка токена доступа |
| **RefreshToken** | Обновление токена |
| **RevokeToken** | Аннулирование токена |

Пример запроса:
```proto
rpc Login(LoginRequest) returns (LoginResponse);
```

---

## **🛠 Структура проекта**
```
gRPC-SSO-Service/
│── server/
│   ├── config/          # Конфигурация сервиса
│   ├── handlers/        # gRPC-методы
│   ├── storage/         # Работа с базой данных
│── proto/               # Определение API (Protobuf)
│── main.go              # Точка входа
│── Dockerfile           # Конфигурация для Docker
│── go.mod               # Модуль Go
```

---

## **📝 To-Do**
- [ ] Поддержка MySQL / PostgreSQL
- [ ] Ролевое управление доступом (RBAC)
- [ ] Подключение к LDAP / Active Directory
- [ ] Поддержка OAuth2 и OpenID Connect

---

## **🤝 Контакты**
📧 Email: [timverhos@gmail.com](mailto:timverhos@gmail.com)  
🟦 Linkedin: [Timofei Verkhososov](https://www.linkedin.com/feed/)

**✨ Если проект вам понравился — ставьте ⭐ на GitHub!** 🚀

