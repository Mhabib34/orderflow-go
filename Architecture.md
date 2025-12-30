# Arsitektur Microservice dengan Message Broker

## Overview

Sistem ini terdiri dari 2 microservices yang berkomunikasi menggunakan message broker (RabbitMQ/Kafka):

1. **Order Service** - Mengelola orders
2. **Notification Service** - Mengelola notifications

## Arsitektur Flow

```
┌─────────────────┐
│     Client      │
└────────┬────────┘
         │
         │ HTTP Request
         │
         ▼
┌─────────────────────────────────────┐
│        Order Service (Port 8001)    │
│  ┌─────────────────────────────┐   │
│  │  API Endpoints:             │   │
│  │  - POST   /orders           │   │
│  │  - GET    /orders           │   │
│  │  - GET    /orders/{id}      │   │
│  │  - PATCH  /orders/{id}      │   │
│  └─────────────────────────────┘   │
│               │                      │
│               │ Publish Event        │
│               ▼                      │
│  ┌─────────────────────────────┐   │
│  │    Message Producer          │   │
│  └─────────────────────────────┘   │
└────────────┬────────────────────────┘
             │
             │ Events:
             │ - order.created
             │ - order.updated
             │ - order.completed
             │ - order.cancelled
             │
             ▼
┌─────────────────────────────────────┐
│     Message Broker (RabbitMQ)       │
│                                     │
│  Exchange: order.events             │
│  Queue: notification.queue          │
│  Routing Key: order.*               │
└────────────┬────────────────────────┘
             │
             │ Consume Event
             │
             ▼
┌─────────────────────────────────────┐
│  Notification Service (Port 8002)   │
│  ┌─────────────────────────────┐   │
│  │    Message Consumer          │   │
│  └─────────────────────────────┘   │
│               │                      │
│               │ Save to DB           │
│               ▼                      │
│  ┌─────────────────────────────┐   │
│  │  API Endpoints:             │   │
│  │  - GET    /notifications    │   │
│  │  - GET    /notifications/{id}│  │
│  │  - PATCH  /notifications/{id}│  │
│  │  - POST   /mark-all-read    │   │
│  └─────────────────────────────┘   │
└─────────────────────────────────────┘
```

## Event Structure

### Order Events

Semua event yang dipublish oleh Order Service memiliki struktur:

```json
{
  "event_type": "order.created | order.updated | order.completed | order.cancelled",
  "event_id": "uuid",
  "timestamp": "2025-12-30T10:30:00Z",
  "data": {
    "order_id": "uuid",
    "status": "pending | processing | completed | cancelled",
    "total_amount": 150000.0,
    "previous_status": "pending" // hanya untuk order.updated
  }
}
```

## Message Broker Setup

### RabbitMQ Configuration

```yaml
# Exchange
Name: order.events
Type: topic
Durable: true

# Queue
Name: notification.queue
Durable: true
Auto-delete: false

# Binding
Exchange: order.events
Queue: notification.queue
Routing Key: order.*
```

### Event Routing Keys

- `order.created` - Ketika order baru dibuat
- `order.updated` - Ketika status order berubah
- `order.completed` - Ketika order selesai
- `order.cancelled` - Ketika order dibatalkan

## Databases

### Order Service Database

```
Database: order_service_db
Table: orders
- Menyimpan data orders
```

### Notification Service Database

```
Database: notification_service_db
Table: notifications
- Menyimpan data notifications yang dibuat dari event
```

## Flow Scenario

### 1. Create Order

1. Client mengirim POST request ke `/orders`
2. Order Service membuat order baru di database
3. Order Service publish event `order.created` ke message broker
4. Notification Service consume event dari queue
5. Notification Service membuat notification baru di database
6. Client dapat mengakses notification melalui GET `/notifications`

### 2. Update Order Status

1. Client mengirim PATCH request ke `/orders/{id}`
2. Order Service update status order di database
3. Order Service publish event `order.updated` ke message broker
4. Notification Service consume event dari queue
5. Notification Service membuat notification baru di database

## Technology Stack (Rekomendasi)

### Order Service

- **Runtime**: Node.js / Go / Python
- **Framework**: Express / Gin / FastAPI
- **Database**: PostgreSQL (install local)
- **Message Broker Client**: amqplib / amqp091-go / pika

### Notification Service

- **Runtime**: Node.js / Go / Python
- **Framework**: Express / Gin / FastAPI
- **Database**: PostgreSQL (install local)
- **Message Broker Client**: amqplib / amqp091-go / pika

### Message Broker

- **Option 1**: RabbitMQ (Recommended for beginners) - install local
- **Option 2**: Apache Kafka (For high throughput)

### Infrastructure

- **Local Development**: Semua service berjalan di local machine

## Environment Variables

### Order Service (.env)

```
PORT=8001
DB_HOST=localhost
DB_PORT=5432
DB_NAME=order_service_db
DB_USER=postgres
DB_PASSWORD=password

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_EXCHANGE=order.events
```

### Notification Service (.env)

```
PORT=8002
DB_HOST=localhost
DB_PORT=5432
DB_NAME=notification_service_db
DB_USER=postgres
DB_PASSWORD=password

RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=notification.queue
```

## Next Steps untuk Development

1. **Setup Infrastructure**

   - Install Docker & Docker Compose
   - Setup PostgreSQL containers
   - Setup RabbitMQ container

2. **Develop Order Service**

   - Implementasi API endpoints
   - Setup database connection
   - Implementasi message producer

3. **Develop Notification Service**

   - Setup message consumer
   - Implementasi API endpoints
   - Setup database connection

4. **Testing**

   - Unit testing
   - Integration testing
   - End-to-end testing dengan message flow

5. **Monitoring & Logging**
   - Add logging untuk setiap service
   - Monitor message queue
   - Track event processing

## Keuntungan Arsitektur Ini

✅ **Decoupling** - Services tidak saling bergantung secara langsung
✅ **Scalability** - Setiap service dapat di-scale secara independent
✅ **Resilience** - Jika satu service down, yang lain tetap jalan
✅ **Asynchronous** - Processing tidak blocking
✅ **Extensibility** - Mudah menambah service baru yang consume event yang sama
