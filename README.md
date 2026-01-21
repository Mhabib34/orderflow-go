# OrderFlow-Go - Microservices Learning Project

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.8+-FF6600?style=flat&logo=rabbitmq)](https://www.rabbitmq.com/)

Project sederhana untuk belajar **microservice architecture** dan **message broker** menggunakan Go, PostgreSQL, dan RabbitMQ. Project ini mengimplementasikan Order Management System dengan event-driven architecture.

## 📚 Table of Contents

- [Arsitektur](#-arsitektur)
- [Komponen](#-komponen)
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Service Details](#-service-details)
- [API Documentation](#-api-documentation)
- [Development](#-development)
- [Testing](#-testing)
- [Monitoring](#-monitoring)
- [Troubleshooting](#-troubleshooting)
- [Learning Objectives](#-learning-objectives)
- [Next Steps](#-next-steps)

## 🏗 Arsitektur

```
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│  Order Service  │         │ Payment Service │         │Notification Svc │
│    (Port 3000)  │         │   (Port 8000)   │         │   (Port 8080)   │
└────────┬────────┘         └────────┬────────┘         └────────┬────────┘
         │                           │                           │
         │  Publish Events           │  Publish Events           │  Consume
         │  ┌────────────────────────┼───────────────────────────┘  Events
         │  │                        │
         ▼  ▼                        ▼
    ┌────────────────────────────────────┐
    │         RabbitMQ Broker            │
    │                                    │
    │  Exchanges: order_exchanges        │
    │             payment_exchanges      │
    │                                    │
    │  Queues: notification.order.created│
    │         payment.order.created      │
    │         order.payment.status.updated│
    └────────────────────────────────────┘
         │           │           │
         ▼           ▼           ▼
    ┌─────────┐ ┌─────────┐ ┌─────────┐
    │PostgreSQL│ │PostgreSQL│ │PostgreSQL│
    │ order_  │ │ payment_│ │notification_│
    │ service │ │ service │ │  service   │
    └─────────┘ └─────────┘ └─────────┘
```

Lihat detail arsitektur di [Architecture.md](Architecture.md)

## 📋 Komponen

| Service | Port | Database | Fungsi |
|---------|------|----------|--------|
| **Order Service** | 8001 | `order_service` | Mengelola pembuatan dan pembaruan order |
| **Payment Service** | 8003 | `payment_service` | Mengelola proses pembayaran |
| **Notification Service** | 8002 | `notification_service` | Mengelola notifikasi kepada user |
| **RabbitMQ** | 5672, 15672 | - | Message broker untuk komunikasi async |
| **PostgreSQL** | 5432 | Multiple DBs | Database untuk setiap service |

## 🔧 Prerequisites

Pastikan sudah terinstall:

- **Go** v1.21 atau lebih baru - [Download](https://go.dev/dl/)
- **PostgreSQL** v15 atau lebih baru - [Download](https://www.postgresql.org/download/)
- **RabbitMQ** v3.8 atau lebih baru - [Download](https://www.rabbitmq.com/download.html)
- **Git** - [Download](https://git-scm.com/downloads)
- **Postman** atau **cURL** - untuk testing API

## 🚀 Quick Start

### 1. Clone Repository

```bash
git clone https://github.com/Mhabib34/orderflow-go.git
cd orderflow-go
```

### 2. Setup PostgreSQL

#### Windows

```bash
# Buka Command Prompt atau PowerShell sebagai Administrator
# Masuk ke PostgreSQL
psql -U postgres

# Buat databases
CREATE DATABASE order_service;
CREATE DATABASE payment_service;
CREATE DATABASE notification_service;
\q

# Import schema (jika ada file schema)
# psql -U postgres -d order_service -f order-service/schema.sql
# psql -U postgres -d payment_service -f payment_service/schema.sql
# psql -U postgres -d notification_service -f notification-service/schema.sql
```

#### macOS

```bash
# Install PostgreSQL (jika belum)
brew install postgresql@15
brew services start postgresql@15

# Buat databases
createdb order_service
createdb payment_service
createdb notification_service

# Import schema (jika ada)
# psql -d order_service -f order-service/schema.sql
# psql -d payment_service -f payment_service/schema.sql
# psql -d notification_service -f notification-service/schema.sql
```

#### Linux (Ubuntu/Debian)

```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Buat databases
sudo -u postgres createdb order_service
sudo -u postgres createdb payment_service
sudo -u postgres createdb notification_service

# Import schema (jika ada)
# sudo -u postgres psql -d order_service -f order-service/schema.sql
# sudo -u postgres psql -d payment_service -f payment_service/schema.sql
# sudo -u postgres psql -d notification_service -f notification-service/schema.sql
```

### 3. Setup RabbitMQ

#### Windows

```bash
# Download dan install Erlang dari: https://www.erlang.org/downloads
# Download dan install RabbitMQ dari: https://www.rabbitmq.com/download.html

# Enable Management Plugin (Command Prompt as Administrator)
rabbitmq-plugins enable rabbitmq_management

# Service akan start otomatis
```

#### macOS

```bash
# Install RabbitMQ
brew install rabbitmq

# Start service
brew services start rabbitmq

# Enable Management Plugin
rabbitmq-plugins enable rabbitmq_management
```

#### Linux (Ubuntu/Debian)

```bash
# Install RabbitMQ
sudo apt install rabbitmq-server

# Start service
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server

# Enable Management Plugin
sudo rabbitmq-plugins enable rabbitmq_management
```

### 4. Setup RabbitMQ Exchange & Queue

Akses RabbitMQ Management UI di [http://localhost:15672](http://localhost:15672)
- Username: `guest`
- Password: `guest`

**Setup untuk Order Service:**
1. **Create Exchange**: `order_exchanges` (type: `topic`, durable: `true`)
2. **Create Queues**:
   - `notification.order.created` (durable: `true`)
   - `payment.order.created` (durable: `true`)
3. **Create Bindings**:
   - Exchange: `order_exchanges` → Queue: `notification.order.created`, Routing Key: `order.created`
   - Exchange: `order_exchanges` → Queue: `payment.order.created`, Routing Key: `order.created`

**Setup untuk Payment Service:**
1. **Create Exchange**: `payment_exchanges` (type: `topic`, durable: `true`)
2. **Create Queue**: `order.payment.status.updated` (durable: `true`)
3. **Create Binding**:
   - Exchange: `payment_exchanges` → Queue: `order.payment.status.updated`, Routing Key: `payment.status.updated`

### 5. Setup Environment Variables

Buat file `.env` di masing-masing service directory:

**order-service/.env**
```env
PORT=8001
DB_HOST=localhost
DB_PORT=5432
DB_NAME=order_service
DB_USER=postgres
DB_PASSWORD=your_postgres_password

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_EXCHANGE=order_exchanges
```

**payment_service/.env**
```env
PORT=8003
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payment_service
DB_USER=postgres
DB_PASSWORD=your_postgres_password

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=payment.order.created
RABBITMQ_EXCHANGE=payment_exchanges
```

**notification-service/.env**
```env
PORT=8002
DB_HOST=localhost
DB_PORT=5432
DB_NAME=notification_service
DB_USER=postgres
DB_PASSWORD=your_postgres_password

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=notification.order.created
```

### 6. Install Dependencies & Run Services

```bash
# Terminal 1 - Order Service
cd order-service
go mod download
go run main.go

# Terminal 2 - Payment Service
cd payment_service
go mod download
go run main.go

# Terminal 3 - Notification Service
cd notification-service
go mod download
go run main.go
```

## 🔍 Service Details

### Order Service

**Struktur:**
```
order-service/
├── internal/
│   ├── controllers/    # HTTP handlers
│   ├── models/         # Data models
│   ├── repository/     # Database layer
│   ├── usecase/        # Business logic
│   ├── routes/         # API routes
│   └── broker/         # RabbitMQ publisher
├── main.go
├── .env
└── go.mod
```

**Responsibilities:**
- Menerima request pembuatan order
- Menyimpan order ke database
- Publish event `order.created` ke RabbitMQ
- Update status order berdasarkan payment status
- Menyediakan API untuk query orders

### Payment Service

**Struktur:**
```
payment_service/
├── internal/
│   ├── controllers/    # HTTP handlers
│   ├── models/         # Data models
│   ├── repository/     # Database layer
│   ├── usecase/        # Business logic
│   ├── routes/         # API routes
│   ├── broker/         # RabbitMQ publisher
│   └── consumers/      # RabbitMQ consumer
├── main.go
├── .env
└── go.mod
```

**Responsibilities:**
- Consume event `order.created` dari RabbitMQ
- Proses pembayaran (simulasi)
- Update payment status
- Publish event `payment.status.updated` ke RabbitMQ

### Notification Service

**Struktur:**
```
notification-service/
├── internal/
│   ├── controllers/    # HTTP handlers
│   ├── models/         # Data models
│   ├── repository/     # Database layer
│   ├── usecase/        # Business logic
│   ├── routes/         # API routes
│   └── consumers/      # RabbitMQ consumer
├── main.go
├── .env
└── go.mod
```

**Responsibilities:**
- Consume events dari RabbitMQ
- Membuat notifikasi untuk user
- Menyimpan notifikasi ke database
- Menyediakan API untuk query notifications

## 📖 API Documentation

### Order Service API

**Base URL:** `http://localhost:8001/api/v1`

#### Create Order
```bash
POST /orders
Content-Type: application/json

{
  "total_amount": 150000.00
}
```

**Response:**
```json
{
  "id": "uuid",
  "total_amount": 150000.00,
  "status": "pending",
  "created_at": "2024-01-21T10:00:00Z"
}
```

#### Get All Orders
```bash
GET /orders
```

#### Get Order by ID
```bash
GET /orders/{orderId}
```

#### Update Order Status
```bash
PATCH /orders/{orderId}
Content-Type: application/json

{
  "status": "completed"
}
```

### Payment Service API

**Base URL:** `http://localhost:8003/api/v1`

#### Get Payment by Order ID
```bash
GET /payments/order/{orderId}
```

#### Get All Payments
```bash
GET /payments
```

### Notification Service API

**Base URL:** `http://localhost:8002/api/v1`

#### Get All Notifications
```bash
GET /notifications
```

#### Get Notification by ID
```bash
GET /notifications/{notificationId}
```

#### Mark as Read
```bash
PATCH /notifications/{notificationId}
```

## 🧪 Testing

### End-to-End Flow

1. **Create Order**
```bash
curl -X POST http://localhost:8001/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"total_amount": 150000.00}'
```

**Expected Results:**
- Order created in `order_service` database
- Event published to `order_exchanges`
- Payment created in `payment_service` database
- Notification created in `notification_service` database

2. **Check Payment**
```bash
curl http://localhost:8003/api/v1/payments
```

3. **Check Notifications**
```bash
curl http://localhost:8002/api/v1/notifications
```

4. **Verify in RabbitMQ**
- Open [http://localhost:15672](http://localhost:15672)
- Check queues for message counts
- Monitor message rates

### Unit Testing

```bash
# Run tests for each service
cd order-service
go test ./... -v

cd ../payment_service
go test ./... -v

cd ../notification-service
go test ./... -v
```

## 📊 Monitoring

### RabbitMQ Management UI

- **URL:** [http://localhost:15672](http://localhost:15672)
- **Credentials:** guest / guest
- **Features:**
  - Monitor queue depth
  - Check message rates
  - View active connections
  - Inspect exchanges and bindings

### PostgreSQL Monitoring

**Tools yang direkomendasikan:**
- [pgAdmin](https://www.pgadmin.org/) - Full-featured GUI
- [DBeaver](https://dbeaver.io/) - Universal database tool
- [TablePlus](https://tableplus.com/) - Modern database GUI

**Connection Details:**
```
Host: localhost
Port: 5432
Databases: order_service, payment_service, notification_service
Username: postgres
Password: [your_password]
```

### Application Logs

Monitor console output dari setiap service untuk:
- HTTP request/response logs
- Database query logs
- RabbitMQ connection status
- Error messages dan stack traces

## 🔧 Troubleshooting

### PostgreSQL Issues

**Problem: Cannot connect to database**
```bash
# Check PostgreSQL status
# Windows: services.msc → Look for PostgreSQL
# macOS: brew services list | grep postgresql
# Linux: sudo systemctl status postgresql

# Restart if needed
# Windows: Restart via services.msc
# macOS: brew services restart postgresql@15
# Linux: sudo systemctl restart postgresql
```

**Problem: Database doesn't exist**
```bash
# List databases
psql -U postgres -l

# Create if missing
createdb order_service
createdb payment_service
createdb notification_service
```

**Problem: Password authentication failed**
- Verify password in `.env` files
- Check `pg_hba.conf` for authentication method
- Reset password: `ALTER USER postgres PASSWORD 'new_password';`

### RabbitMQ Issues

**Problem: Cannot connect to RabbitMQ**
```bash
# Check RabbitMQ status
# Windows: services.msc → Look for RabbitMQ
# macOS: brew services list | grep rabbitmq
# Linux: sudo systemctl status rabbitmq-server

# Restart if needed
# Windows: Restart via services.msc
# macOS: brew services restart rabbitmq
# Linux: sudo systemctl restart rabbitmq-server
```

**Problem: Management UI not accessible**
```bash
# Enable management plugin
rabbitmq-plugins enable rabbitmq_management

# Restart RabbitMQ after enabling
```

**Problem: Messages not being consumed**
- Verify consumer is running (check service logs)
- Check queue bindings in Management UI
- Verify routing keys match
- Check for errors in consumer service logs

### Application Issues

**Problem: Port already in use**
```bash
# Find process using port
# Windows: netstat -ano | findstr :8001
# macOS/Linux: lsof -i :8001

# Change PORT in .env or kill the process
```

**Problem: Go module issues**
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy
```

## 🎯 Learning Objectives

Project ini dirancang untuk mempelajari:

### 1. Microservice Architecture ✅
- Service independence dan isolation
- Database per service pattern
- Service communication patterns
- API design dan versioning

### 2. Event-Driven Architecture ✅
- Asynchronous communication
- Event publishing dan consuming
- Message broker integration
- Event-driven workflows

### 3. Message Broker (RabbitMQ) ✅
- Exchange types (topic, direct, fanout)
- Queue management
- Routing keys dan bindings
- Publisher/Consumer pattern
- Message durability

### 4. Clean Architecture ✅
- Separation of concerns
- Dependency injection
- Repository pattern
- Use case layer
- Controller/Handler layer

### 5. Database Design ✅
- Schema design
- Relationships
- Indexing strategies
- Database migrations

### 6. Go Best Practices ✅
- Project structure
- Error handling
- Environment configuration
- Dependency management

## 📚 Resources

### Documentation
- [Go Documentation](https://go.dev/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [Microservices Patterns](https://microservices.io/patterns/)

### Learning Materials
- [Event-Driven Architecture by Martin Fowler](https://martinfowler.com/articles/201701-event-driven.html)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [RabbitMQ Best Practices](https://www.rabbitmq.com/production-checklist.html)

### Tools
- [Postman](https://www.postman.com/) - API testing
- [pgAdmin](https://www.pgadmin.org/) - PostgreSQL GUI
- [DBeaver](https://dbeaver.io/) - Database tool
- [Docker](https://www.docker.com/) - Containerization

## 🤝 Contributing

Contributions are welcome! Ini adalah project pembelajaran, silakan:

1. Fork repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Mhabib34**
- GitHub: [@Mhabib34](https://github.com/Mhabib34)

---

**Happy Learning! 🚀**

*If you find this project helpful, please give it a ⭐*
