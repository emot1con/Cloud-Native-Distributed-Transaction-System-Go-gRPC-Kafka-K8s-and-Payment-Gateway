#Cloud-Native-Distributed-Transaction-System-Go-gRPC-Kafka-K8s-and-Payment-Gateway

> Production-grade distributed transaction system implementing event-driven architecture with gRPC inter-service communication, Kafka message streaming, Redis caching layer, and Kubernetes orchestration for order and payment processing at scale.

**Distributed System** showcasing:

- ✅ **Event-Driven Architecture** with Kafka for async order-payment workflow
- ✅ **gRPC Service Mesh** for high-performance inter-service communication
- ✅ **Distributed Caching Strategy** with Redis TTL and invalidation patterns
- ✅ **Token Bucket Rate Limiting** using Lua scripts in Redis
- ✅ **Saga Pattern Implementation** for distributed transactions
- ✅ **Cloud-Native Deployment** with Kubernetes health checks and resource limits
- ✅ **Real Payment Gateway Integration** (Midtrans) with webhook verification
- ✅ **Clean Architecture** with proper layering across all microservices

##  Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                    CLIENT LAYER                              │
│              Next.js 16 + TypeScript + React Query           │
│                     http://localhost:3000                    │
└────────────────────────┬─────────────────────────────────────┘
                         │ HTTP/REST
                         ▼
┌────────────────────────────────────────────────────────────┐
│                   API GATEWAY (Broker Service)             │
│                         Port: 8080 / 30005 (K8s)           │
│  • JWT Authentication & Authorization                      │
│  • Token Bucket Rate Limiter (Redis + Lua)                 │
│  • HTTP → gRPC Translation Layer                           │
│  • Request/Response Orchestration                          │
└──┬────────┬────────┬────────┬──────────────────────────────┘
   │        │        │        │
   │ gRPC   │ gRPC   │ gRPC   │ gRPC
   │        │        │        │
   ▼        ▼        ▼        ▼
┌──────┐ ┌────────┐ ┌──────┐ ┌────────┐
│ USER │ │PRODUCT │ │ORDER │ │PAYMENT │◄─── Midtrans Webhook
│50001 │ │ 40001  │ │30001 │ │ 60001  │     (SHA512 Signature)
└──┬───┘ └───┬────┘ └──┬───┘ └───┬────┘
   │         │         │         │
   │         │         │         │
   │         │         ▼         │
   │         │    ┌─────────┐    │
   │         │    │  KAFKA  │─── ┼ ────── order.created
   │         │    │  :9092  │    │       (Event Streaming)
   │         │    └─────────┘    │
   │         │                   │
   │         │         ▼         │
   │         │    ┌─────────┐    │
   │         │    │  REDIS  │─── ┘──── Cache + Rate Limit
   │         │    │  :6379  │         
   │         │    └─────────┘
   │         │
   ▼         ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL 13 Database (:5432)                 │
│  • users • products • orders • order_items • payments       │
└─────────────────────────────────────────────────────────────┘
```

##  Project Structure

```
go_microservice/
├── broker/                 # API Gateway Service
│   ├── auth/              # JWT middleware
│   ├── cmd/               # Application entry points
│   ├── handler/           # HTTP handlers for all services
│   ├── middleware/        # Rate limiter, JWT validation
│   ├── proto/             # gRPC proto files
│   ├── repository/        # gRPC client connections
│   └── transport/         # Kafka producer
│
├── user/                   # User Authentication Service
│   ├── auth/              # OAuth, JWT, Password utilities
│   ├── cmd/               # Application & database setup
│   ├── proto/             # gRPC definitions
│   ├── repository/        # Database operations
│   ├── service/           # Business logic
│   ├── transport/         # gRPC & HTTP handlers
│   └── types/             # Data structures
│
├── product/                # Product Catalog Service
│   ├── cmd/               # Application & Redis setup
│   ├── helper/            # Utility functions
│   ├── migrations/        # Database migrations
│   ├── proto/             # gRPC definitions
│   ├── repository/        # Database operations with caching
│   ├── service/           # Business logic + cache layer
│   └── transport/         # gRPC handlers
│
├── order/                  # Order Management Service
│   ├── cmd/               # Application & Redis setup
│   │   └── db/            # Redis cache helpers
│   ├── helper/            # Utility functions
│   ├── proto/             # gRPC definitions
│   ├── repository/        # Database operations
│   ├── service/           # Business logic + cache invalidation
│   └── transport/         # gRPC & Kafka producer
│
├── payment/                # Payment Processing Service
│   ├── client/            # Midtrans SDK wrapper
│   ├── cmd/               # Application entry
│   ├── helper/            # Utility functions
│   ├── proto/             # gRPC definitions
│   ├── repository/        # Database operations + gRPC clients
│   ├── service/           # Payment workflow + webhook handling
│   └── transport/         # gRPC & Kafka consumer
│
├── fe/                     # Frontend Application
│   ├── src/
│   │   ├── app/           # Next.js App Router pages
│   │   ├── components/    # Reusable UI components
│   │   ├── services/      # API client services
│   │   ├── store/         # Zustand state management
│   │   └── types/         # TypeScript types
│
└── project/                # Infrastructure Configuration
    ├── docker-compose.yaml
    ├── Makefile
    ├── migrations/         # Database schema
    ├── kafka.yaml          # Kubernetes Kafka config
    ├── postgres.yaml       # Kubernetes PostgreSQL config
    └── README.md           # This file
```

##  Services Overview

### 1. Broker Service (API Gateway)
The central entry point for all client requests with advanced middleware.

**Features:**
- RESTful API endpoints
- JWT-based authentication with user context propagation
- **Token bucket rate limiter** (capacity: 20, refill: 10 req/s per user)
- CORS configuration
- Request routing to microservices via gRPC
- Detailed logging with user tracking

**Tech Stack:** Go, gRPC Client, Redis (rate limiting), Kafka Producer

**Port:** `8080` (Docker) / `30005` (Kubernetes NodePort)

**Docker Image:** `numpyh/broker:v2.4`

---

### 2. User Service
Handles user authentication and management with OAuth 2.0 support.

**Features:**
- User registration & login
- JWT token generation & refresh (15min access, 7 days refresh)
- OAuth2 integration (Google, Facebook, GitHub) with PKCE
- Password hashing with bcrypt
- User profile retrieval by ID

**Tech Stack:** Go, gRPC Server, PostgreSQL, JWT, OAuth 2.0

**Port:** `50001` (gRPC)

**Docker Image:** `numpyh/user:v2.0`

---

### 3. Product Service
Manages the product catalog with intelligent caching.

**Features:**
- CRUD operations for products
- Paginated product listing
- Stock management and validation
- **Redis caching with 10-minute TTL**
- Cache invalidation on write operations

**Tech Stack:** Go, gRPC Server, PostgreSQL, Redis Cache

**Port:** `40001` (gRPC)

**Docker Image:** `numpyh/product:v2.0`

---

### 4. Order Service
Handles order creation and lifecycle management with event streaming.

**Features:**
- Create orders with multiple items and stock validation
- Order listing with pagination (15 items per page)
- **Cache-aside pattern**: order lists (5min TTL) + single orders (10min TTL)
- **GetOrderById** with user validation for secure access
- Order status updates (pending → paid/failed)
- **Smart cache invalidation**: both list and single order caches
- Kafka event publishing (`order.created`)
- gRPC server for status updates from payment service

**Tech Stack:** Go, gRPC Server, PostgreSQL, Redis Cache, Kafka Producer

**Port:** `30001` (gRPC)

**Docker Image:** `numpyh/order:v2.2`

---

### 5. Payment Service
Handles payment gateway integration and transaction orchestration.

**Features:**
- **Midtrans Snap API integration** for multiple payment methods
- Payment initiation with token generation and expiry (24 hours)
- Support for various payment channels:
  - Virtual Account (BCA, BNI, Mandiri, Permata, BRI)
  - E-Wallet (GoPay, OVO, DANA, ShopeePay, LinkAja)
  - Credit Card (Visa, Mastercard, JCB)
  - QRIS, Akulaku, Kredivo, Indomaret, Alfamart
- **Webhook signature verification** (SHA512 with server key)
- Payment status mapping (capture, settlement, pending, deny, expire, cancel)
- **Automatic order status sync** via gRPC to order service
- Kafka consumer for `order.created` events (async payment creation)
- Idempotency support (reuse existing gateway token if pending)

**Tech Stack:** Go, gRPC Server/Client, PostgreSQL, Kafka Consumer, Midtrans SDK

**Port:** `60001` (gRPC)

**Docker Image:** `numpyh/payment:v2.0`

---

### 6. Frontend (Next.js)
Modern, responsive e-commerce web application.

**Features:**
- Product catalog with search and pagination
- Shopping cart management
- Order history with pagination (synchronized with backend: 15 items/page)
- **Order detail view with direct ID lookup** (fixes pagination issue)
- Midtrans Snap payment popup integration
- JWT-based authentication with refresh token
- Toast notifications for user feedback
- Optimistic UI updates with React Query
- TypeScript for type safety

**Tech Stack:** Next.js 16, React 18, TanStack React Query, Zustand, TailwindCSS, Axios

**Port:** `3000`

##  Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Products Table
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DOUBLE PRECISION NOT NULL,
    stock INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Orders Table
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    total_price DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, paid, failed, cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Order Items Table
```sql
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Payments Table
```sql
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    amount DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, paid, failed, expired
    payment_method VARCHAR(50),            -- gopay, bank_transfer, credit_card, etc.
    payment_channel VARCHAR(50),           -- bca_va, gopay, credit_card, etc.
    gateway_order_id VARCHAR(100),         -- Midtrans order ID (PAY-{id}-{timestamp})
    gateway_token TEXT,                    -- Midtrans Snap token
    gateway_redirect_url TEXT,             -- Midtrans payment page URL
    gateway_transaction_id VARCHAR(100),   -- Midtrans transaction ID
    va_number VARCHAR(50),                 -- Virtual Account number
    qr_code_url TEXT,                      -- QR code for QRIS
    expired_at VARCHAR(50),                -- Payment expiration time
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔌 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/register` | Register new user | ❌ |
| POST | `/login` | Login user | ❌ |
| POST | `/refresh` | Refresh JWT token | ❌ |
| GET | `/oauth/google` | Google OAuth redirect URL | ❌ |
| GET | `/oauth/facebook` | Facebook OAuth redirect URL | ❌ |
| GET | `/oauth/github` | GitHub OAuth redirect URL | ❌ |
| GET | `/oauth/google/callback` | Google OAuth callback | ❌ |
| GET | `/oauth/facebook/callback` | Facebook OAuth callback | ❌ |
| GET | `/oauth/github/callback` | GitHub OAuth callback | ❌ |

### Products
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/product?page={n}&page_size={m}` | List products (paginated, cached) | ❌ |
| GET | `/product/{id}` | Get product by ID | ❌ |
| POST | `/product` | Create product | ✅ (+ Rate Limited) |
| PUT | `/product/{id}` | Update product | ✅ |
| DELETE | `/product/{id}` | Delete product | ✅ |

### Orders
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/order` | Create order | ✅ (+ Rate Limited) |
| GET | `/order?page={n}&page_size={m}` | Get user orders (paginated, cached) | ✅ |
| GET | `/order/{id}` | Get order by ID (cached) | ✅ |

### Payments
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/payment/order/{order_id}` | Get payment by order ID | ✅ |
| POST | `/payment/initiate` | Initiate Midtrans payment | ✅ |
| POST | `/payment/webhook` | Midtrans webhook (signature verified) | ❌ (Signature) |

##  gRPC Services

### AuthService (User Service)
```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (EmptyResponse);
  rpc Login(LoginRequest) returns (TokenResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (TokenResponse);
  rpc GetUserByID(GetUserRequest) returns (User);
  rpc GoogleOauth(EmptyRequest) returns (URLResponse);
  rpc FacebookOauth(EmptyRequest) returns (URLResponse);
  rpc GithubOauth(EmptyRequest) returns (URLResponse);
  rpc GoogleOauthCallback(CodeRequest) returns (TokenResponse);
  rpc FacebookOauthCallback(CodeRequest) returns (TokenResponse);
  rpc GithubOauthCallback(CodeRequest) returns (TokenResponse);
}
```

### ProductService
```protobuf
service ProductService {
  rpc CreateProduct(ProductRequest) returns (Empty);
  rpc GetProduct(GetProductRequest) returns (Product);
  rpc ListProducts(Offset) returns (ProductList);
  rpc UpdateProduct(Product) returns (Product);
  rpc DeleteProduct(GetProductRequest) returns (Empty);
}
```

### OrderService
```protobuf
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (OrderResponse);
  rpc GetOrder(GetOrderRequest) returns (OrdersResponse);          // Paginated list
  rpc GetOrderById(GetOrderByIdRequest) returns (OrderResponse);   // Single order
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (EmptyOrder);
}

message GetOrderByIdRequest {
  int32 order_id = 1;
  int32 user_id = 2;  // Security: validate order ownership
}
```

### PaymentService
```protobuf
service PaymentService {
  rpc CreatePayment(CreatePaymentRequest) returns (PaymentResponse);
  rpc GetPaymentByOrderId(GetPaymentByOrderIdRequest) returns (PaymentResponse);
  rpc InitiatePayment(InitiatePaymentRequest) returns (InitiatePaymentResponse);
  rpc HandleWebhook(WebhookRequest) returns (EmptyPayment);
}
```

##  Getting Started

### Prerequisites
- **Go** 1.24+
- **Docker** 24+ & Docker Compose
- **Kubernetes** (Minikube/Kind/K3s) for K8s deployment
- **Protocol Buffers** compiler (`protoc`) with Go plugin
- **Node.js** 20+ (for frontend)
- **kubectl** CLI tool

### Environment Setup

Each service requires specific environment variables. Create `.env` files:

**Broker Service** (`broker/.env`)
```env
PORT=8080
JWT_SECRET=your-jwt-secret-key-min-32-chars
JWT_REFRESH_SECRET=your-refresh-secret-key

# Service addresses (Docker Compose)
USER_SERVICE_ADDR=user-service:50001
PRODUCT_SERVICE_ADDR=product-service:40001
ORDER_SERVICE_ADDR=order-service:30001
PAYMENT_SERVICE_ADDR=payment-service:60001

# For Kubernetes, use service names:
# USER_SERVICE_ADDR=user-service:50001
# PRODUCT_SERVICE_ADDR=product-service:40001
# etc.

# Redis
REDIS_URL=redis:6379

# Rate Limiter Config
RATE_LIMIT_CAPACITY=20
RATE_LIMIT_REFILL_RATE=10
```

**User Service** (`user/.env`)
```env
PORT=50001
DATABASE_URL=postgres://postgres:password@postgres:5432/microservice_db?sslmode=disable

# JWT
JWT_SECRET=your-jwt-secret-key-min-32-chars
JWT_REFRESH_SECRET=your-refresh-secret-key

# OAuth 2.0
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/oauth/google/callback

FACEBOOK_CLIENT_ID=your-facebook-app-id
FACEBOOK_CLIENT_SECRET=your-facebook-app-secret
FACEBOOK_REDIRECT_URL=http://localhost:8080/oauth/facebook/callback

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/oauth/github/callback
```

**Product Service** (`product/.env`)
```env
PORT=40001
DATABASE_URL=postgres://postgres:password@postgres:5432/microservice_db?sslmode=disable
REDIS_URL=redis:6379
```

**Order Service** (`order/.env`)
```env
PORT=30001
DATABASE_URL=postgres://postgres:password@postgres:5432/microservice_db?sslmode=disable
REDIS_URL=redis:6379
KAFKA_BROKER_URL=kafka:9092
KAFKA_ORDER_TOPIC=order.created
PRODUCT_SERVICE_ADDR=product-service:40001
```

**Payment Service** (`payment/.env`)
```env
PORT=60001
DATABASE_URL=postgres://postgres:password@postgres:5432/microservice_db?sslmode=disable
KAFKA_BROKER_URL=kafka:9092
KAFKA_ORDER_TOPIC=order.created
ORDER_SERVICE_ADDR=order-service:30001

# Midtrans Configuration
MIDTRANS_SERVER_KEY=your-midtrans-server-key
MIDTRANS_CLIENT_KEY=your-midtrans-client-key
MIDTRANS_ENVIRONMENT=sandbox  # or "production"
```

**Frontend** (`fe/.env.local`)
```env
# Docker Compose
NEXT_PUBLIC_API_URL=http://localhost:8080

# Kubernetes (get from: minikube service broker-service --url)
# NEXT_PUBLIC_API_URL=http://192.168.49.2:30005
```

---

### Local Development with Docker Compose

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/distributed-order-payment-system.git
   cd distributed-order-payment-system/project
   ```

2. **Configure environment variables**
   ```bash
   # Create .env files for each service (see Environment Setup above)
   # Or copy from examples if available
   ```

3. **Start infrastructure services**
   ```bash
   # Start PostgreSQL, Redis, Kafka, Zookeeper
   docker-compose up -d postgres redis kafka zookeeper
   
   # Wait for services to be healthy
   docker-compose ps
   ```

4. **Run database migrations**
   ```bash
   docker-compose up migrate
   ```

5. **Build and start microservices**
   ```bash
   make up_build
   ```

6. **Verify services are running**
   ```bash
   docker-compose ps
   curl http://localhost:8080/health  # Broker health check
   ```

7. **Access the application**
   - **API Gateway:** `http://localhost:8080`
   - **Frontend:** `http://localhost:3000` (if running)
   - **PostgreSQL:** `localhost:5432`
   - **Redis:** `localhost:6379`
   - **Kafka:** `localhost:9092`

---

### Kubernetes Deployment

1. **Start Minikube cluster**
   ```bash
   minikube start --cpus=4 --memory=8192
   ```

2. **Create Kubernetes secrets**
   ```bash
   # Create secret for broker service
   kubectl create secret generic broker-service-secret \
     --from-literal=JWT_SECRET=your-jwt-secret \
     --from-literal=JWT_REFRESH_SECRET=your-refresh-secret \
     --from-literal=REDIS_URL=redis:6379 \
     --from-literal=USER_SERVICE_ADDR=user-service:50001 \
     --from-literal=PRODUCT_SERVICE_ADDR=product-service:40001 \
     --from-literal=ORDER_SERVICE_ADDR=order-service:30001 \
     --from-literal=PAYMENT_SERVICE_ADDR=payment-service:60001

   # Create secret for user service
   kubectl create secret generic user-service-secret \
     --from-literal=DATABASE_URL=postgres://postgres:password@postgres:5432/microservice_db?sslmode=disable \
     --from-literal=JWT_SECRET=your-jwt-secret \
     --from-literal=GOOGLE_CLIENT_ID=your-google-id \
     --from-literal=GOOGLE_CLIENT_SECRET=your-google-secret
   
   # Similar for other services...
   ```

3. **Deploy infrastructure (PostgreSQL, Redis, Kafka)**
   ```bash
   kubectl apply -f postgres.yaml
   kubectl apply -f kafka.yaml
   
   # Wait for pods to be ready
   kubectl get pods -w
   ```

4. **Deploy microservices**
   ```bash
   # Deploy in dependency order
   kubectl apply -f ../user/deployment.yaml
   kubectl apply -f ../product/deployment.yaml
   kubectl apply -f ../order/deployment.yaml
   kubectl apply -f ../payment/deployment.yaml
   kubectl apply -f ../broker/deployment.yaml
   
   # Check deployment status
   kubectl get deployments
   kubectl get pods
   ```

5. **Access via NodePort**
   ```bash
   # Get broker service URL
   minikube service broker-service --url
   # Output: http://192.168.49.2:30005
   
   # Test API
   curl http://$(minikube ip):30005/health
   ```

6. **View logs**
   ```bash
   # Stream logs from a service
   kubectl logs -f deployment/broker-service
   kubectl logs -f deployment/order-service
   
   # View all pods in namespace
   kubectl get pods -o wide
   ```

7. **Update deployment (after code changes)**
   ```bash
   # Build new image
   cd ../broker
   docker build -f Dockerfile.prod -t numpyh/broker:v2.5 .
   docker push numpyh/broker:v2.5
   
   # Update deployment.yaml with new tag
   # Then apply
   kubectl apply -f deployment.yaml
   
   # Force pod recreation
   kubectl delete pod -l app=broker-service
   ```

---

### Makefile Commands

```bash
# Docker Compose
make up              # Start all services (no rebuild)
make up_build        # Build and start all services
make down            # Stop all services

# Build individual services
make build_broker
make build_user
make build_product
make build_order
make build_payment

# Run migrations
make migrate_up
make migrate_down
```

---

### Frontend Development

```bash
# Navigate to frontend directory
cd ../fe

# Install dependencies
npm install

# Configure API endpoint
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local

# For Kubernetes
echo "NEXT_PUBLIC_API_URL=http://$(minikube ip):30005" > .env.local

# Start development server
npm run dev

# Open browser to http://localhost:3000

# Build for production
npm run build
npm start
```

## ⚡ **Key Technical Implementations**

### 1. Event-Driven Saga Pattern
```go
// Order Service publishes event to Kafka
func (s *OrderService) CreateOrder(payload *proto.CreateOrderRequest) {
    // Create order in database
    orderID := s.orderRepo.CreateOrder(payload, tx)
    
    // Publish event asynchronously
    kafka.Produce("order.created", OrderCreatedEvent{
        OrderID: orderID,
        UserID: payload.UserId,
        TotalPrice: payload.TotalPrice,
    })
}

// Payment Service consumes event
func (c *PaymentConsumer) ConsumeOrderCreated() {
    for msg := range kafkaMessages {
        order := unmarshal(msg)
        payment := createPaymentFromOrder(order)
        // Async - non-blocking for order service
    }
}

// Payment webhook updates order status via gRPC
func (s *PaymentService) HandleWebhook(webhook) {
    if webhook.Status == "settlement" {
        orderClient.UpdateOrderStatus(orderID, "paid")
    }
}
```

### 2. Cache-Aside Pattern with Smart Invalidation
```go
// 3-Layer Caching Strategy
const (
    orderListCacheTTL = 5 * time.Minute   // High churn
    singleOrderCacheTTL = 10 * time.Minute // Medium churn
    productCacheTTL = 10 * time.Minute     // Low churn
)

// Read: Check cache first
func (s *OrderService) GetOrderById(orderID, userID int) (*Order, error) {
    cacheKey := fmt.Sprintf("order:id:%d", orderID)
    
    // Try cache first
    if cached, err := redis.Get(cacheKey); err == nil {
        return unmarshal(cached), nil
    }
    
    // Cache miss - fetch from DB
    order, err := s.orderRepo.GetOrderById(orderID, userID, db)
    if err != nil {
        return nil, err
    }
    
    // Update cache with TTL
    redis.SetEx(cacheKey, marshal(order), singleOrderCacheTTL)
    return order, nil
}

// Write: Invalidate affected caches
func (s *OrderService) UpdateOrderStatus(orderID int, status string) error {
    // Update database
    s.orderRepo.UpdateOrderStatus(status, orderID, tx)
    
    // Invalidate single order cache
    redis.Del(fmt.Sprintf("order:id:%d", orderID))
    
    // Invalidate all user order lists
    redis.DelPattern("orders:user:*")
    
    return nil
}
```

### 3. Token Bucket Rate Limiter (Lua Script)
```go
// Rate limiter middleware
func RateLimiterMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Value(UserKey).(int)
        
        allowed, err := rateLimiter.Allow(userID)
        if !allowed {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// Redis Lua script for atomic token bucket
const luaScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local tokens = redis.call('GET', key)
if not tokens then
    tokens = capacity
else
    tokens = tonumber(tokens)
end

if tokens > 0 then
    redis.call('DECR', key)
    redis.call('EXPIRE', key, 60)
    return 1  -- Allow request
else
    return 0  -- Rate limited
end
`
```

### 4. Secure User Context Propagation
```go
// JWT Middleware extracts user info
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        claims := validateToken(token)
        
        // Store userID in context
        ctx := context.WithValue(c.Request.Context(), UserKey, claims.UserID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}

// Handler retrieves userID from context
func (h *OrderHandler) GetOrderById(c *gin.Context) {
    orderID := c.Param("id")
    userID := c.Request.Context().Value(UserKey).(int)
    
    // Validate order ownership
    order, err := h.orderRepo.GetOrderById(orderID, userID)
}
```

---

##  Technologies Used

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.24 | Primary programming language |
| **gRPC** | - | High-performance inter-service communication |
| **Protocol Buffers** | v3 | Service definition & serialization |
| **PostgreSQL** | 13 | Primary relational database (ACID) |
| **Redis** | 7-alpine | Distributed cache + rate limiting |
| **Apache Kafka** | 7.4.0 | Event streaming platform (order events) |
| **Zookeeper** | latest | Kafka cluster coordination |
| **Docker** | 24+ | Containerization & multi-stage builds |
| **Kubernetes** | 1.28+ | Container orchestration (Minikube) |
| **JWT** | HS256 | Stateless authentication tokens |
| **OAuth2** | - | Social login (Google, Facebook, GitHub) |
| **Midtrans** | Snap API | Payment gateway integration |
| **Next.js** | 16 | Frontend framework (App Router + RSC) |
| **TypeScript** | 5+ | Type-safe frontend development |
| **React Query** | 5.90 | Server state management & caching |
| **TailwindCSS** | 3+ | Utility-first CSS framework |
| **Logrus** | - | Structured logging for Go services |

##  Features

### 🏗️ **Architecture & Design**
-  **Microservices Architecture** - Loosely coupled, independently deployable services
-  **Clean Architecture** - Layered design: Transport → Handler → Service → Repository
-  **gRPC Communication** - High-performance RPC with Protocol Buffers
-  **REST API Gateway** - Single entry point for external clients
-  **Event-Driven Design** - Kafka for async order-payment workflow

### 🔐 **Security**
-  **JWT Authentication** - Access (15min) + Refresh (7 days) tokens
-  **OAuth2 Integration** - Social login with PKCE support
-  **Token Bucket Rate Limiting** - Per-user limits using Redis + Lua scripts
-  **User Context Propagation** - Secure userID passing across services
-  **Webhook Signature Verification** - SHA512 for Midtrans webhooks

### ⚡ **Performance & Scalability**
-  **Redis Caching Strategy**:
   - Product list: 10min TTL
   - Order list (per user): 5min TTL
   - Single order: 10min TTL
-  **Smart Cache Invalidation** - Delete on write operations
-  **Connection Pooling** - PostgreSQL max connections
-  **Horizontal Scaling** - Kubernetes replica sets
-  **Resource Management** - CPU/Memory limits per pod

### 🚀 **Deployment & Operations**
-  **Docker Support** - Multi-stage builds for small images
-  **Kubernetes Ready** - Deployments, Services, ConfigMaps, Secrets
-  **Health Checks** - Liveness & Readiness probes
-  **Database Migrations** - Version-controlled schema changes
-  **Graceful Shutdown** - Signal handling for zero downtime
-  **Structured Logging** - JSON-formatted logs with context

### 💳 **Payment Integration**
-  **Midtrans Snap API** - Multiple payment channels support
-  **Async Payment Creation** - Kafka consumer for order events
-  **Payment Status Sync** - Webhook → gRPC → Order status update
-  **Idempotency** - Reuse existing tokens for pending payments
-  **Expiry Handling** - 24-hour payment window

## 🧪 Testing

### Unit Tests
```bash
# Test specific service
cd ../user && go test ./... -v -cover

cd ../product && go test ./... -v -cover

cd ../order && go test ./... -v -cover

cd ../payment && go test ./... -v -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests
```bash
# With Docker testcontainers
go test -tags=integration ./...

# Test with real dependencies
docker-compose -f docker-compose.test.yaml up --abort-on-container-exit
```

### API Testing
```bash
# Using curl
bash scripts/test_api.sh

# Using Postman/Insomnia
# Import collection from: docs/postman_collection.json
```

---

## 🎯 Production Readiness Checklist

- ✅ **Clean Architecture** - Proper layering and separation of concerns
- ✅ **Error Handling** - Consistent error responses across services
- ✅ **Logging** - Structured JSON logs with context
- ✅ **Health Checks** - Liveness and readiness probes
- ✅ **Resource Limits** - CPU and memory constraints
- ✅ **Connection Pooling** - Database connection management
- ✅ **Graceful Shutdown** - Signal handling for zero downtime
- ✅ **Rate Limiting** - Protection against abuse
- ✅ **Caching Strategy** - Redis with TTL and invalidation
- ✅ **Security** - JWT + OAuth + Rate limiting
- ✅ **Payment Security** - Webhook signature verification
- ✅ **Event-Driven** - Async communication via Kafka
- ✅ **Idempotency** - Prevent duplicate payment processing

---

##  API Examples

### 1. Register User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "message": "User registered successfully"
}
```

---

### 2. Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john@example.com"
  }
}
```

---

### 3. List Products (Cached)
```bash
curl http://localhost:8080/product?page=1&page_size=15
```

**Response:**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Laptop Dell XPS 13",
      "description": "Ultra-portable laptop",
      "price": 15000000,
      "stock": 10
    }
  ],
  "total": 50,
  "page": 1,
  "page_size": 15
}
```

---

### 4. Create Order (Protected + Rate Limited)
```bash
# Get token from login first
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X POST http://localhost:8080/order \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      },
      {
        "product_id": 3,
        "quantity": 1
      }
    ]
  }'
```

**Response:**
```json
{
  "order": {
    "id": 123,
    "user_id": 1,
    "total_price": 45000000,
    "status": "pending",
    "order_items": [
      {
        "id": 1,
        "product_id": 1,
        "quantity": 2,
        "price": 30000000
      },
      {
        "id": 2,
        "product_id": 3,
        "quantity": 1,
        "price": 15000000
      }
    ],
    "created_at": "2026-01-05T10:30:00Z"
  }
}
```

---

### 5. Get Order by ID (Protected + Cached)
```bash
curl http://localhost:8080/order/123 \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**
```json
{
  "order": {
    "id": 123,
    "user_id": 1,
    "total_price": 45000000,
    "status": "pending",
    "order_items": [...],
    "created_at": "2026-01-05T10:30:00Z"
  }
}
```

---

### 6. Initiate Payment (Midtrans)
```bash
curl -X POST http://localhost:8080/payment/initiate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": 123,
    "customer_name": "John Doe",
    "customer_email": "john@example.com"
  }'
```

**Response:**
```json
{
  "payment_id": 456,
  "gateway_token": "66e4fa55-fdac-4ef9-91b5-733b97d1b862",
  "gateway_redirect_url": "https://app.sandbox.midtrans.com/snap/v3/redirection/...",
  "status": "pending",
  "expired_at": "2026-01-06T10:30:00Z"
}
```

---

### 7. Test Rate Limiter
```bash
# Send 25 rapid requests (capacity: 20)
for i in {1..25}; do
  curl -X POST http://localhost:8080/product \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Product","price":10000,"stock":5}' &
done

# Expected: First 20 succeed (200), remaining 5 get rate limited (429)
```

**Rate Limited Response:**
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

---

## 🐛 **Troubleshooting**

### Issue: Rate limiter returns 429 immediately
**Cause:** Context key mismatch between JWT middleware and rate limiter

**Solution:**
```go
// Ensure UserKey constant matches in both files
const UserKey = "user_id"  // jwt.go
const UserKey = "user_id"  // rate_limiter.go
```

---

### Issue: Order not found when accessing page 2+ orders
**Cause:** Frontend only fetched page 1, then searched within those results

**Solution:** Use `GetOrderById` endpoint instead of filtering from list
```typescript
// ✅ Correct
const order = await orderService.getOrderById(orderId)

// ❌ Wrong
const orders = await orderService.getOrders(1)
const order = orders.find(o => o.id === orderId)  // Only works for page 1!
```

---

### Issue: Payment status not updating in frontend after payment
**Cause:** Redis cache not invalidated when order status changes

**Solution:** Invalidate both order list and single order caches
```go
func (s *OrderService) UpdateOrderStatus(orderID int, status string) {
    // Update DB
    s.orderRepo.UpdateOrderStatus(status, orderID, tx)
    
    // Invalidate caches
    redis.Del(fmt.Sprintf("order:id:%d", orderID))  // Single order
    redis.DelPattern("orders:user:*")                // All lists
}
```

---

### Issue: Kafka consumer not processing messages
**Check:**
```bash
# Verify Kafka topic exists
kubectl exec -it kafka -- kafka-topics --bootstrap-server localhost:9092 --list

# Check consumer group lag
kubectl exec -it kafka -- kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe --group payment-service-group
```

---

### Issue: gRPC connection refused
**Check:**
```bash
# Verify service is running
kubectl get pods -l app=order-service

# Check service endpoints
kubectl get endpoints order-service

# Test gRPC port
kubectl exec -it broker-service -- nc -zv order-service 30001
```

---

### Issue: Database connection pool exhausted
**Solution:** Adjust max connections
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

---

## 📊 **Monitoring**

### Check Service Health
```bash
# Kubernetes pods status
kubectl get pods -A

# Service logs
kubectl logs -f deployment/broker-service
kubectl logs -f deployment/order-service --tail=100

# Describe pod for events
kubectl describe pod broker-service-xxx
```

### Redis Monitoring
```bash
# Connect to Redis CLI
kubectl exec -it deploy/redis -- redis-cli

# Check cache statistics
INFO stats

# View all keys
KEYS *

# Check specific cache
GET "order:id:123"
TTL "order:id:123"

# Check rate limiter keys
KEYS "rate_limit:user:*"

# Monitor commands in real-time
MONITOR
```

### Kafka Monitoring
```bash
# List topics
kubectl exec -it kafka -- kafka-topics --bootstrap-server localhost:9092 --list

# Describe topic
kubectl exec -it kafka -- kafka-topics \
  --bootstrap-server localhost:9092 \
  --describe --topic order.created

# Consumer group status
kubectl exec -it kafka -- kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --describe --group payment-service-group

# Consume messages (debug)
kubectl exec -it kafka -- kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic order.created \
  --from-beginning
```

### Database Monitoring
```bash
# Connect to PostgreSQL
kubectl exec -it postgres -- psql -U postgres -d microservice_db

# Check active connections
SELECT count(*) FROM pg_stat_activity;

# View table sizes
SELECT 
  schemaname,
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

# Check slow queries
SELECT pid, now() - query_start as duration, query 
FROM pg_stat_activity 
WHERE state = 'active' AND now() - query_start > interval '2 seconds';
```

---

## 🎓 **Learning Outcomes**

After studying this project, you'll understand:

- ✅ How to design and implement **event-driven microservices**
- ✅ **gRPC vs REST** trade-offs and when to use each
- ✅ **Distributed transaction patterns** (Saga, 2PC)
- ✅ **Cache invalidation strategies** (write-through, cache-aside)
- ✅ **Rate limiting algorithms** (token bucket, leaky bucket)
- ✅ **Kubernetes resource management** (limits, requests, probes)
- ✅ **Payment gateway integration** (webhooks, idempotency)
- ✅ **Clean architecture** in Go microservices
- ✅ **Inter-service communication** patterns
- ✅ **Security best practices** (JWT, OAuth, rate limiting)

---

## 🛣️ **Roadmap**

- [ ] **Observability**: Prometheus + Grafana metrics
- [ ] **Distributed Tracing**: Jaeger/Zipkin integration
- [ ] **CI/CD Pipeline**: GitHub Actions for automated builds
- [ ] **Service Mesh**: Istio for traffic management
- [ ] **GraphQL Gateway**: Alternative to REST
- [ ] **WebSocket Support**: Real-time order updates
- [ ] **Admin Dashboard**: Order management UI
- [ ] **Multi-tenancy**: Support multiple merchants
- [ ] **CQRS Pattern**: Separate read/write models
- [ ] **Event Replay**: Kafka event sourcing

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Coding Standards:**
- Follow Go conventions (gofmt, golint)
- Write tests for new features
- Update documentation
- Keep PRs focused and small

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

---

## 👤 Author

**emot1con**
- GitHub: [@numpyh](https://github.com/emot1con)
- Docker Hub: [numpyh](https://hub.docker.com/u/numpyh)

---

## 🙏 Acknowledgments

- **Midtrans** for payment gateway integration
- **Confluent** for Kafka Docker images
- **Kubernetes Community** for excellent orchestration tools
- **Go Community** for amazing libraries and tools

---

## 📚 References

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Midtrans API Documentation](https://docs.midtrans.com/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**⭐ If you find this project helpful, please consider giving it a star!**

**💬 Questions? Open an issue or start a discussion!**

---

*Last updated: January 5, 2026*
