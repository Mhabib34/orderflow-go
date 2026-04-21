# OrderFlow-Go - Microservices Learning Project

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.8+-FF6600?style=flat&logo=rabbitmq)](https://www.rabbitmq.com/)

A project for learning **microservice architecture** and **message brokers** using Go, PostgreSQL, and RabbitMQ. This project implements an Order Management System built on an event-driven architecture.

## 📚 Table of Contents

- [Architecture](#-architecture)
- [Components](#-components)
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

## 🏗 Architecture

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

For detailed architecture information, see [Architecture.md](Architecture.md)

## 📋 Components

| Service | Port | Database | Responsibilities |
|---------|------|----------|-----------------|
| **Order Service** | 8001 | `order_service` | Manages order creation and status updates |
| **Payment Service** | 8003 | `payment_service` | Handles payment processing |
| **Notification Service** | 8002 | `notification_service` | Manages user notifications |
| **RabbitMQ** | 5672, 15672 | - | Message broker for async communication |
| **PostgreSQL** | 5432 | Multiple DBs | Per-service relational database |

## 🔧 Prerequisites

Make sure the following are installed on your machine:

- **Go** v1.21 or later — [Download](https://go.dev/dl/)
- **PostgreSQL** v15 or later — [Download](https://www.postgresql.org/download/)
- **RabbitMQ** v3.8 or later — [Download](https://www.rabbitmq.com/download.html)
- **Git** — [Download](https://git-scm.com/downloads)
- **Postman** or **cURL** — for API testing

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/Mhabib34/orderflow-go.git
cd orderflow-go
```

### 2. Set Up PostgreSQL

#### Windows

```bash
# Open Command Prompt or PowerShell as Administrator
# Connect to PostgreSQL
psql -U postgres

# Create the required databases
CREATE DATABASE order_service;
CREATE DATABASE payment_service;
CREATE DATABASE notification_service;
\q

# Import schema (if schema files are available)
# psql -U postgres -d order_service -f order-service/schema.sql
# psql -U postgres -d payment_service -f payment_service/schema.sql
# psql -U postgres -d notification_service -f notification-service/schema.sql
```

#### macOS

```bash
# Install PostgreSQL (if not already installed)
brew install postgresql@15
brew services start postgresql@15

# Create the required databases
createdb order_service
createdb payment_service
createdb notification_service

# Import schema (if available)
# psql -d order_service -f order-service/schema.sql
# psql -d payment_service -f payment_service/schema.sql
# psql -d notification_service -f notification-service/schema.sql
```

#### Linux (Ubuntu/Debian)

```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start the service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create the required databases
sudo -u postgres createdb order_service
sudo -u postgres createdb payment_service
sudo -u postgres createdb notification_service

# Import schema (if available)
# sudo -u postgres psql -d order_service -f order-service/schema.sql
# sudo -u postgres psql -d payment_service -f payment_service/schema.sql
# sudo -u postgres psql -d notification_service -f notification-service/schema.sql
```

### 3. Set Up RabbitMQ

#### Windows

```bash
# Download and install Erlang from: https://www.erlang.org/downloads
# Download and install RabbitMQ from: https://www.rabbitmq.com/download.html

# Enable the Management Plugin (run Command Prompt as Administrator)
rabbitmq-plugins enable rabbitmq_management

# The service will start automatically
```

#### macOS

```bash
# Install RabbitMQ
brew install rabbitmq

# Start the service
brew services start rabbitmq

# Enable the Management Plugin
rabbitmq-plugins enable rabbitmq_management
```

#### Linux (Ubuntu/Debian)

```bash
# Install RabbitMQ
sudo apt install rabbitmq-server

# Start the service
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server

# Enable the Management Plugin
sudo rabbitmq-plugins enable rabbitmq_management
```

### 4. Configure RabbitMQ Exchanges & Queues

Access the RabbitMQ Management UI at [http://localhost:15672](http://localhost:15672)

- **Username:** `guest`
- **Password:** `guest`

**Order Service Setup:**

1. **Create Exchange:** `order_exchanges` (type: `topic`, durable: `true`)
2. **Create Queues:**
   - `notification.order.created` (durable: `true`)
   - `payment.order.created` (durable: `true`)
3. **Create Bindings:**
   - Exchange: `order_exchanges` → Queue: `notification.order.created`, Routing Key: `order.created`
   - Exchange: `order_exchanges` → Queue: `payment.order.created`, Routing Key: `order.created`

**Payment Service Setup:**

1. **Create Exchange:** `payment_exchanges` (type: `topic`, durable: `true`)
2. **Create Queue:** `order.payment.status.updated` (durable: `true`)
3. **Create Binding:**
   - Exchange: `payment_exchanges` → Queue: `order.payment.status.updated`, Routing Key: `payment.status.updated`

### 5. Configure Environment Variables

Create a `.env` file inside each service directory:

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
# Terminal 1 — Order Service
cd order-service
go mod download
go run main.go

# Terminal 2 — Payment Service
cd payment_service
go mod download
go run main.go

# Terminal 3 — Notification Service
cd notification-service
go mod download
go run main.go
```

## 🔍 Service Details

### Order Service

**Directory Structure:**

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
- Accepts incoming order creation requests
- Persists orders to the database
- Publishes `order.created` events to RabbitMQ
- Updates order status based on payment events
- Provides a query API for orders

### Payment Service

**Directory Structure:**

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
- Consumes `order.created` events from RabbitMQ
- Processes payments (simulated)
- Updates payment status
- Publishes `payment.status.updated` events to RabbitMQ

### Notification Service

**Directory Structure:**

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
- Consumes events from RabbitMQ
- Creates user notifications
- Persists notifications to the database
- Provides a query API for notifications

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

#### Mark Notification as Read

```bash
PATCH /notifications/{notificationId}
```

## 🧪 Testing

### End-to-End Flow

**1. Create an Order**

```bash
curl -X POST http://localhost:8001/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"total_amount": 150000.00}'
```

**Expected results:**
- Order created in the `order_service` database
- Event published to `order_exchanges`
- Payment record created in the `payment_service` database
- Notification created in the `notification_service` database

**2. Check Payments**

```bash
curl http://localhost:8003/api/v1/payments
```

**3. Check Notifications**

```bash
curl http://localhost:8002/api/v1/notifications
```

**4. Verify in RabbitMQ**
- Open [http://localhost:15672](http://localhost:15672)
- Inspect queue message counts
- Monitor message throughput rates

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
  - Track message rates
  - View active connections
  - Inspect exchanges and bindings

### PostgreSQL Monitoring

**Recommended tools:**
- [pgAdmin](https://www.pgadmin.org/) — Full-featured GUI client
- [DBeaver](https://dbeaver.io/) — Universal database management tool
- [TablePlus](https://tableplus.com/) — Modern, lightweight database GUI

**Connection details:**

```
Host:      localhost
Port:      5432
Databases: order_service, payment_service, notification_service
Username:  postgres
Password:  [your_password]
```

### Application Logs

Monitor the console output of each service for:
- HTTP request and response logs
- Database query logs
- RabbitMQ connection status
- Error messages and stack traces

## 🔧 Troubleshooting

### PostgreSQL Issues

**Problem: Cannot connect to the database**

```bash
# Check PostgreSQL status
# Windows:  services.msc → Look for PostgreSQL
# macOS:    brew services list | grep postgresql
# Linux:    sudo systemctl status postgresql

# Restart if needed
# Windows:  Restart via services.msc
# macOS:    brew services restart postgresql@15
# Linux:    sudo systemctl restart postgresql
```

**Problem: Database does not exist**

```bash
# List all databases
psql -U postgres -l

# Create the missing databases
createdb order_service
createdb payment_service
createdb notification_service
```

**Problem: Password authentication failed**
- Verify the password in each `.env` file
- Check `pg_hba.conf` for the configured authentication method
- Reset the password: `ALTER USER postgres PASSWORD 'new_password';`

---

### RabbitMQ Issues

**Problem: Cannot connect to RabbitMQ**

```bash
# Check RabbitMQ status
# Windows:  services.msc → Look for RabbitMQ
# macOS:    brew services list | grep rabbitmq
# Linux:    sudo systemctl status rabbitmq-server

# Restart if needed
# Windows:  Restart via services.msc
# macOS:    brew services restart rabbitmq
# Linux:    sudo systemctl restart rabbitmq-server
```

**Problem: Management UI is not accessible**

```bash
# Enable the management plugin
rabbitmq-plugins enable rabbitmq_management

# Restart RabbitMQ after enabling the plugin
```

**Problem: Messages are not being consumed**
- Confirm the consumer service is running (check service logs)
- Verify queue bindings in the Management UI
- Ensure routing keys match what the publisher is sending
- Check for errors in the consumer service logs

---

### Application Issues

**Problem: Port already in use**

```bash
# Find the process using the port
# Windows:      netstat -ano | findstr :8001
# macOS/Linux:  lsof -i :8001

# Update the PORT value in .env, or terminate the conflicting process
```

**Problem: Go module errors**

```bash
# Clear the module cache
go clean -modcache

# Re-download and tidy dependencies
go mod download
go mod tidy
```

## 🎯 Learning Objectives

This project is designed to teach the following concepts:

### 1. Microservice Architecture ✅
- Service independence and isolation
- Database-per-service pattern
- Inter-service communication patterns
- API design and versioning

### 2. Event-Driven Architecture ✅
- Asynchronous communication
- Event publishing and consuming
- Message broker integration
- Event-driven workflows

### 3. Message Brokers (RabbitMQ) ✅
- Exchange types (topic, direct, fanout)
- Queue management
- Routing keys and bindings
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
- Entity relationships
- Indexing strategies
- Database migrations

### 6. Go Best Practices ✅
- Project structure
- Error handling
- Environment configuration
- Dependency management

## 📚 Resources

### Official Documentation
- [Go Documentation](https://go.dev/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [Microservices Patterns](https://microservices.io/patterns/)

### Learning Materials
- [Event-Driven Architecture by Martin Fowler](https://martinfowler.com/articles/201701-event-driven.html)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [RabbitMQ Best Practices](https://www.rabbitmq.com/production-checklist.html)

### Recommended Tools
- [Postman](https://www.postman.com/) — API testing
- [pgAdmin](https://www.pgadmin.org/) — PostgreSQL GUI
- [DBeaver](https://dbeaver.io/) — Universal database tool
- [Docker](https://www.docker.com/) — Containerization

## 🤝 Contributing

Contributions are welcome! Since this is a learning project, feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License — see the LICENSE file for details.

## 👨‍💻 Author

**Mhabib34** — GitHub: [@Mhabib34](https://github.com/Mhabib34)

---

**Happy Learning! 🚀**

*If you find this project helpful, please give it a ⭐*
