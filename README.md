# Microservice Learning Project

Project sederhana untuk belajar microservice architecture dan message broker menggunakan Order Service dan Notification Service.

## 📋 Komponen

- **Order Service**: Service untuk mengelola orders
- **Notification Service**: Service untuk mengelola notifications
- **Message Broker**: RabbitMQ untuk komunikasi antar service
- **Database**: PostgreSQL untuk masing-masing service

## 🚀 Quick Start

### 1. Prerequisites

Pastikan sudah terinstall:

- **PostgreSQL** (versi 12 atau lebih baru)
- **RabbitMQ** (versi 3.8 atau lebih baru)
- **Go**
- **Git**
- **Postman** (untuk testing API)

### 2. Setup PostgreSQL

#### Windows

```bash
# Download dan install PostgreSQL dari: https://www.postgresql.org/download/windows/

# Setelah install, buat database untuk Order Service
psql -U postgres
CREATE DATABASE order_service;
\q

# Buat database untuk Notification Service
psql -U postgres
CREATE DATABASE notification_service;
\q

# Import schema ke Order Service database
psql -U postgres -d order_service -f database-schema.sql

# Import schema ke Notification Service database
psql -U postgres -d notification_service -f database-schema.sql
```

#### macOS

```bash
# Install PostgreSQL menggunakan Homebrew
brew install postgresql@15

# Start PostgreSQL service
brew services start postgresql@15

# Buat databases
createdb order_service
createdb notification_service

# Import schema
psql -d order_service -f database-schema.sql
psql -d notification_service -f database-schema.sql
```

#### Linux (Ubuntu/Debian)

```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Buat databases
sudo -u postgres createdb order_service
sudo -u postgres createdb notification_service

# Import schema
sudo -u postgres psql -d order_service -f database-schema.sql
sudo -u postgres psql -d notification_service -f database-schema.sql
```

### 3. Setup RabbitMQ

#### Windows

```bash
# Download dan install Erlang terlebih dahulu dari:
# https://www.erlang.org/downloads

# Download dan install RabbitMQ dari:
# https://www.rabbitmq.com/download.html

# Enable Management Plugin
rabbitmq-plugins enable rabbitmq_management

# Start RabbitMQ
# (Service akan start otomatis setelah install)
```

#### macOS

```bash
# Install RabbitMQ menggunakan Homebrew
brew install rabbitmq

# Start RabbitMQ service
brew services start rabbitmq

# Enable Management Plugin
rabbitmq-plugins enable rabbitmq_management
```

#### Linux (Ubuntu/Debian)

```bash
# Install RabbitMQ
sudo apt update
sudo apt install rabbitmq-server

# Start RabbitMQ service
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server

# Enable Management Plugin
sudo rabbitmq-plugins enable rabbitmq_management
```

### 4. Verifikasi Installation

**PostgreSQL:**

```bash
# Cek PostgreSQL berjalan
psql -U postgres -l

# Seharusnya melihat order_service dan notification_service dalam list
```

**RabbitMQ:**

```bash
# Cek RabbitMQ berjalan
sudo rabbitmqctl status

# Akses Management UI
# Buka browser: http://localhost:15672
# Username: guest
# Password: guest
```

### 5. Setup RabbitMQ Exchange & Queue

**Via RabbitMQ Management UI (Recommended untuk pemula):**

1. Buka http://localhost:15672
2. Login dengan username: `guest`, password: `guest`
3. Pergi ke tab **Exchanges**
   - Click "Add a new exchange"
   - Name: `order_exchanges`
   - Type: `topic`
   - Durability: `Durable`
   - Click "Add exchange"
4. Pergi ke tab **Queues**
   - Click "Add a new queue"
   - Name: `notification_queue`
   - Durability: `Durable`
   - Click "Add queue"
5. Kembali ke tab **Exchanges**, click `order.events`
   - Scroll ke bawah ke bagian "Bindings"
   - To queue: `notification_queue`
   - Routing key: `order.*`
   - Click "Bind"

**Via Command Line (Alternative):**

```bash
# Untuk Linux/macOS
rabbitmqadmin declare exchange name=order.events type=topic durable=true
rabbitmqadmin declare queue name=notification.queue durable=true
rabbitmqadmin declare binding source=order.events destination=notification.queue routing_key="order.*"

# Untuk Windows (gunakan rabbitmqadmin.bat jika sudah di-setup)
# Atau lebih mudah pakai Management UI
```

## 🛠️ Development

### Struktur Service yang Perlu Dibuat

```
order-service/
├── src/
│   ├── routes/          # API routes
│   ├── controllers/     # Business logic
│   ├── models/          # Database models
│   ├── usecase/        # Use Case layer'
│   ├── repository/      # Repository layer
│   └── broker/          # RabbitMQ producer
├── .env

notification-service/
├── src/
│   ├── routes/          # API routes
│   ├── controllers/     # Business logic
│   ├── models/          # Database models
│   ├── usecase/         # Use Case layer
│   ├── repository/      # Repository layer
│   └── consumers/       # RabbitMQ consumer
├── .env
```

### Environment Variables

**Order Service (.env)**

```env
PORT=8001
DB_HOST=localhost
DB_PORT=5432
DB_NAME=order_service
DB_USER=postgres
DB_PASSWORD=postgres

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_EXCHANGE=order_exchanges
```

**Notification Service (.env)**

```env
PORT=8002
DB_HOST=localhost
DB_PORT=5432
DB_NAME=notification_service
DB_USER=postgres
DB_PASSWORD=postgres

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=notification_)queue
```

**Note:**

- Sesuaikan `DB_PASSWORD` dengan password PostgreSQL yang kamu set saat instalasi
- Default PostgreSQL port adalah 5432 untuk kedua database (berbeda database name)

## 📖 API Documentation

### Order Service API

Base URL: `http://localhost:3000/api/v1`

**Endpoints:**

- `POST /orders` - Membuat order baru
- `GET /orders` - Mendapatkan daftar orders
- `GET /orders/{orderId}` - Mendapatkan detail order
- `PATCH /orders/{orderId}` - Update status order

Lihat detail lengkap di `order-service-api.yaml`

### Notification Service API

Base URL: `http://localhost:8000/api/v1`

**Endpoints:**

- `GET /notifications` - Mendapatkan daftar notifications
- `GET /notifications/{notificationId}` - Mendapatkan detail notification
- `PATCH /notifications/{notificationId}` - Mark notification as read

Lihat detail lengkap di `notification-service-api.yaml`

## 🧪 Testing Flow

### 1. Create Order

```bash
curl -X POST http://localhost:3000/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"total_amount": 150000.00}'
```

**Expected:**

- Order created in order_service_db
- Event published to RabbitMQ
- Notification created in notification_service_db

### 2. Check Notifications

```bash
curl http://localhost:3000/api/v1/notifications
```

### 3. Update Order Status

```bash
curl -X PATCH http://localhost:3000/api/v1/orders/{orderId} \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

**Expected:**

- Order status updated in order_service_db
- Event published to RabbitMQ
- New notification created in notification_service_db

### 4. Mark Notification as Read

```bash
curl -X PATCH http://localhost:3000/api/v1/notifications/{notificationId}
```

## 📊 Monitoring

### RabbitMQ Management UI

- URL: http://localhost:15672
- Username: `guest` / Password: `guest`
- Monitor queues, exchanges, dan message flow
- Lihat message rate dan throughput
- Check connections dari services

### PostgreSQL Database

Gunakan tools seperti:

- **pgAdmin** - GUI tool yang powerful
- **DBeaver** - Universal database tool
- **psql CLI** - Command line interface
- **TablePlus** - Modern database GUI

Contoh koneksi untuk Order Service:

```
Host: localhost
Port: 5432
Database: order_service_db
Username: postgres
Password: [your_password]
```

### Application Logs

Monitor console output dari masing-masing service untuk:

- Request/Response logs
- RabbitMQ connection status
- Database queries
- Error messages

## 🔧 Troubleshooting

### PostgreSQL Issues

**Problem: Cannot connect to PostgreSQL**

```bash
# Check jika PostgreSQL service berjalan
# Windows
services.msc (cari PostgreSQL)

# macOS
brew services list | grep postgresql

# Linux
sudo systemctl status postgresql
```

**Problem: Password authentication failed**

- Pastikan password yang digunakan benar
- Cek file `pg_hba.conf` untuk authentication method
- Restart PostgreSQL service setelah perubahan

**Problem: Database not found**

```bash
# List semua databases
psql -U postgres -l

# Buat database jika belum ada
createdb order_service_db
createdb notification_service_db
```

### RabbitMQ Issues

**Problem: Cannot connect to RabbitMQ**

```bash
# Check jika RabbitMQ service berjalan
# Windows
services.msc (cari RabbitMQ)

# macOS
brew services list | grep rabbitmq

# Linux
sudo systemctl status rabbitmq-server
```

**Problem: Management UI tidak bisa diakses**

```bash
# Enable management plugin
rabbitmq-plugins enable rabbitmq_management

# Restart RabbitMQ service
# Windows: restart via services.msc
# macOS: brew services restart rabbitmq
# Linux: sudo systemctl restart rabbitmq-server
```

**Problem: Exchange/Queue tidak ada**

- Login ke Management UI (http://localhost:15672)
- Buat manual sesuai langkah di section "Setup RabbitMQ Exchange & Queue"

### Service Development Issues

**Problem: Port already in use**

```bash
# Check process menggunakan port
# Windows
netstat -ano | findstr :3000

# macOS/Linux
lsof -i :3000

# Kill process jika perlu atau ganti PORT di .env
```

**Problem: Cannot connect to database dari service**

- Pastikan PostgreSQL berjalan
- Cek connection string di .env file
- Test connection manual menggunakan psql
- Pastikan firewall tidak block koneksi

**Problem: Message tidak diterima di Notification Service**

- Cek RabbitMQ Management UI, apakah message masuk ke queue?
- Verifikasi binding antara exchange dan queue
- Cek routing key sesuai dengan yang di-publish
- Pastikan consumer service berjalan

## 🎯 Learning Objectives

Dari project ini kamu akan belajar:

1. ✅ **Microservice Architecture**

   - Service independence
   - Database per service pattern
   - API design

2. ✅ **Message Broker**

   - Asynchronous communication
   - Event-driven architecture
   - Publisher/Subscriber pattern
   - RabbitMQ setup dan configuration

3. ✅ **Database Design**

   - Schema design
   - Indexing
   - Triggers
   - Multi-database management

4. ✅ **Local Development Environment**
   - PostgreSQL installation dan configuration
   - RabbitMQ setup
   - Service orchestration tanpa containerization

## 🔄 Next Steps

Setelah basic implementation, kamu bisa tambahkan:

1. **Containerization (Recommended next step!)**

   - Dockerize Order Service
   - Dockerize Notification Service
   - Docker Compose orchestration
   - Multi-stage builds

2. **Error Handling & Retry**

   - Dead letter queue
   - Retry mechanism
   - Circuit breaker

3. **Authentication & Authorization**

   - JWT tokens
   - API Gateway

4. **Logging & Monitoring**

   - Centralized logging (ELK stack)
   - Metrics (Prometheus + Grafana)
   - Distributed tracing (Jaeger)

5. **Additional Services**

   - Payment Service
   - Email Service
   - User Service

6. **Advanced Features**
   - CQRS pattern
   - Event sourcing
   - Saga pattern

## 📚 Resources

### Installation Guides

- [PostgreSQL Download & Installation](https://www.postgresql.org/download/)
- [RabbitMQ Installation Guide](https://www.rabbitmq.com/download.html)
- [Erlang Download (required for RabbitMQ)](https://www.erlang.org/downloads)

### Learning Resources

- [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
- [RabbitMQ Management Plugin](https://www.rabbitmq.com/management.html)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Microservices Patterns](https://microservices.io/patterns/)
- [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html)

### Tools

- [Postman](https://www.postman.com/downloads/) - API testing
- [pgAdmin](https://www.pgadmin.org/download/) - PostgreSQL GUI
- [DBeaver](https://dbeaver.io/download/) - Universal database tool
- [TablePlus](https://tableplus.com/) - Modern database GUI

## 🤝 Contributing

Ini adalah project pembelajaran, feel free untuk:

- Menambahkan fitur baru
- Memperbaiki bugs
- Menambahkan dokumentasi

## 📝 License

MIT License - silakan digunakan untuk belajar dan eksperimen!

---

**Happy Learning! 🚀**
