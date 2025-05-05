# Recommendation 推薦系統

後端實作，實現使用者註冊、登入、Email 驗證與推薦商品 API，並結合 Redis 快取機制

## Features

- 使用者註冊，含密碼強度驗證
- 使用者登入，採用 JWT 驗證機制
- Email 驗證（模擬發送驗證碼）
- 推薦商品 API，需登入後使用，並模擬資料庫延遲與快取
- 支援 graceful shutdown (包含 MySQL 與 Redis 連線)
- 所有設定統一管理於 `config.yaml`

## Project Structure

```
recommendation
├── cmd
│   └── main.go               # 應用程式進入點
├── config
│   └── config.example.yaml   # 設定檔範例
├── internal
│   ├── auth
│   │   ├── handler.go         # 認證 API HTTP handler
│   │   ├── request.go         # 認證請求參數定義
│   │   ├── service.go         # 認證業務邏輯處理
│   │   └── validator.go       # 密碼與欄位驗證邏輯
│   ├── config
│   │   └── config.go          # 使用 Viper 讀取設定檔
│   ├── email
│   │   └── email.go           # 模擬 Email 發送行為
│   ├── middleware
│   │   └── auth_middleware.go # JWT 驗證中介層
│   ├── recommendation
│   │   ├── handler.go         # 推薦商品 API handler
│   │   ├── model.go           # 推薦商品資料模型
│   │   ├── repository.go      # 推薦商品資料存取
│   │   └── service.go         # 推薦邏輯與 Redis 快取
│   └── user
│       ├── model.go           # 使用者資料模型
│       └── repository.go      # 使用者資料存取邏輯
├── pkg
│   ├── database
│   │   ├── seed.go            # 初始推薦商品資料 seeder
│   │   └── mysql.go           # MySQL 初始化邏輯
│   ├── cache
│   │   └── redis.go           # Redis 客戶端初始化
│   ├── email
│   │   └── mock.go            # 模擬 Email 寄送服務實作
│   ├── logger
│   │   └── logger.go          # logger 初始化
│   └── utils
│       └── password.go        # 密碼加密與驗證工具
├── go.mod                     # Go 模組定義
├── go.sum                     # 相依套件版本資訊
└── README.md                  # 專案說明文件
```
## Dockerization

### Requirement

#### `Dockerfile`
```Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080
CMD ["./main"]
```

#### `docker-compose.yml`
```yaml
version: '3.8'

services:
  app:
    build: .
    container_name: app
    depends_on:
      - mysql
      - redis
    ports:
      - "8080:8080"
    volumes:
      - .:/app
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    container_name: db
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: recommendation
      MYSQL_USER: user
      MYSQL_PASSWORD: password
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql

  redis:
    image: redis:7
    container_name: cache
    restart: always
    ports:
      - "6379:6379"

volumes:
  mysql-data:
```

### Run

```bash
docker-compose up --build
```

## Manual Run with Local Go

1. 啟動服務
```bash
docker-compose up -d mysql redis
```

2. 本地執行:
```bash
go run cmd/main.go
```

須確認 `config.yaml` 使用 `localhost` 連接 MySQL 及 Redis

## 🚀 Install and Run

1. 下載並進入專案：

```bash
git clone https://github.com/jerry1993tw/recommendation.git
cd recommendation
```

2. 安裝相依套件：

```bash
go mod tidy
```

3. 建立設定檔：

```bash
cp config/config.example.yaml config/config.yaml
```

4. 執行專案：

```bash
go run cmd/main.go
```

## API Endpoint

### 註冊
- **POST /register**：使用 email 與密碼註冊

```json
{
  "email": "test@example.com",
  "password": "Abc123$"
}
```

### 驗證 Email
- **POST /verify**：提交 email 與驗證碼進行驗證

```json
{
  "email": "test@example.com",
  "code": "123456"
}
```

### 登入
- **POST /login**：登入並取得 JWT Token

```json
{
  "email": "test@example.com",
  "password": "Abc123$"
}
```

### 取得推薦商品
- **GET /recommendation**：需帶入 Bearer Token

```http
Authorization: Bearer <your-token>
```

