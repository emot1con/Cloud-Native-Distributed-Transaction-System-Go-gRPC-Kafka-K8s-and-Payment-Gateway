# Go Microservices E-Commerce Platform

A modern, scalable microservices architecture built with Go, gRPC, Kafka, and PostgreSQL. This project implements a complete e-commerce platform with user authentication, product management, order processing, and payment handling.

##  Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Frontend (Next.js)                              │
│                            http://localhost:3000                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Broker Service (API Gateway)                         │
│                              Port: 8080                                       │
│                    • REST API Endpoints                                       │
│                    • JWT Authentication                                       │
│                    • Rate Limiting                                            │
│                    • CORS Handling                                            │
└─────────────────────────────────────────────────────────────────────────────┘
                    │           │           │           │
                    ▼           ▼           ▼           ▼
         ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
         │User Service  │ │Product Service│ │Order Service │ │Payment Service│
         │  Port: 50001 │ │  Port: 40001  │ │  Port: 30001 │ │  Port: 60001 │
         │    gRPC      │ │    gRPC       │ │    gRPC      │ │    gRPC      │
         └──────────────┘ └───────────────┘ └──────────────┘ └──────────────┘
                    │           │           │           │
                    └───────────┴─────┬─────┴───────────┘
                                      ▼
                    ┌─────────────────────────────────────┐
                    │        PostgreSQL Database          │
                    │            Port: 5432               │
                    └─────────────────────────────────────┘
                                      │
                    ┌─────────────────────────────────────┐
                    │         Apache Kafka                │
                    │            Port: 9092               │
                    │    (Event-Driven Communication)     │
                    └─────────────────────────────────────┘
```

##  Project Structure

```
go_microservice/
├── broker/                 # API Gateway Service
│   ├── auth/              # JWT middleware
│   ├── cmd/               # Application entry points
│   ├── handler/           # HTTP handlers
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
│   ├── cmd/               # Application entry
│   ├── helper/            # Utility functions
│   ├── migrations/        # Database migrations
│   ├── proto/             # gRPC definitions
│   ├── repository/        # Database operations
│   ├── service/           # Business logic
│   └── transport/         # gRPC handlers
│
├── order/                  # Order Management Service
│   ├── cmd/               # Application entry
│   ├── helper/            # Utility functions
│   ├── proto/             # gRPC definitions
│   ├── repository/        # Database operations
│   ├── service/           # Business logic
│   └── transport/         # gRPC & Kafka handlers
│
├── payment/                # Payment Processing Service
│   ├── cmd/               # Application entry
│   ├── helper/            # Utility functions
│   ├── proto/             # gRPC definitions
│   ├── repository/        # Database operations
│   ├── service/           # Business logic
│   └── transport/         # gRPC & Kafka handlers
│
└── project/                # Infrastructure Configuration
    ├── docker-compose.yaml
    ├── Makefile
    ├── migrations/         # Database schema
    ├── kafka.yaml          # Kubernetes Kafka config
    ├── postgres.yaml       # Kubernetes PostgreSQL config
    └── zookeeper.yaml      # Kubernetes Zookeeper config
```

##  Services Overview

### 1. Broker Service (API Gateway)
The central entry point for all client requests.

**Features:**
- RESTful API endpoints
- JWT-based authentication
- Rate limiting per user
- CORS configuration
- Request routing to microservices via gRPC

**Port:** `8080`

### 2. User Service
Handles user authentication and management.

**Features:**
- User registration & login
- JWT token generation & refresh
- OAuth2 integration (Google, Facebook, GitHub)
- Password hashing with bcrypt

**Port:** `50001` (gRPC)

### 3. Product Service
Manages the product catalog.

**Features:**
- CRUD operations for products
- Paginated product listing
- Stock management

**Port:** `40001` (gRPC)

### 4. Order Service
Handles order creation and management.

**Features:**
- Create orders with multiple items
- Order status tracking
- Order history per user

**Port:** `30001` (gRPC)

### 5. Payment Service
Processes payments for orders.

**Features:**
- Payment creation when order is placed
- Transaction processing
- Payment status management
- Get payment by order ID

**Port:** `60001` (gRPC)

##  Database Schema

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
    user_id INTEGER NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Order Items Table
```sql
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Payments Table
```sql
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔌 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register new user | ❌ |
| POST | `/api/v1/auth/login` | Login user | ❌ |
| GET | `/api/v1/auth/refresh-token` | Refresh JWT token | 🔑 |
| GET | `/api/v1/auth/google` | Google OAuth URL | ❌ |
| GET | `/api/v1/auth/facebook` | Facebook OAuth URL | ❌ |
| GET | `/api/v1/auth/github` | GitHub OAuth URL | ❌ |
| GET | `/api/v1/oauth/google/callback` | Google OAuth callback | ❌ |
| GET | `/api/v1/oauth/facebook/callback` | Facebook OAuth callback | ❌ |
| GET | `/api/v1/oauth/github/callback` | GitHub OAuth callback | ❌ |

### Products
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/product/?page={n}` | List products (paginated) | ❌ |
| POST | `/product/` | Create product | ✅ |
| PUT | `/product/?id={id}` | Update product | ✅ |
| DELETE | `/product/?id={id}` | Delete product | ✅ |

### Orders
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/order/` | Create order | ✅ |
| GET | `/order/?offset={n}` | Get user orders | ✅ |

### Payments
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/payment/transaction` | Process payment | ✅ |
| GET | `/payment/order/{order_id}` | Get payment by order ID | ✅ |

##  gRPC Services

### AuthService (User)
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
  rpc GetOrder(GetOrderRequest) returns (OrdersResponse);
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (EmptyOrder);
}
```

### PaymentService
```protobuf
service PaymentService {
  rpc PayOrder(CreatePaymentRequest) returns (OrderPayment);
  rpc Transaction(PaymentTransaction) returns (EmptyPayment);
  rpc GetPaymentByOrderId(GetPaymentByOrderIdRequest) returns (OrderPayment);
}
```

##  Getting Started

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Kubernetes (Minikube) for K8s deployment
- Protocol Buffers compiler (`protoc`)

### Local Development with Docker Compose

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/go_microservice.git
   cd go_microservice/project
   ```

2. **Set up environment variables**
   ```bash
   # Copy and configure .env files for each service
   cp ../broker/.env.example ../broker/.env
   cp ../user/.env.example ../user/.env
   cp ../product/.env.example ../product/.env
   cp ../order/.env.example ../order/.env
   cp ../payment/.env.example ../payment/.env
   ```

3. **Build and run with Docker Compose**
   ```bash
   make up_build
   ```

4. **Access the services**
   - API Gateway: `http://localhost:8080`
   - PostgreSQL: `localhost:5432`
   - Kafka: `localhost:9092`

### Building Individual Services

```bash
# Build all services
make up_build

# Build individual service
make build_broker
make build_user
make build_product
make build_order
make build_payment

# Stop all services
make down
```

### Kubernetes Deployment

1. **Start Minikube**
   ```bash
   minikube start
   ```

2. **Deploy infrastructure**
   ```bash
   kubectl apply -f postgres.yaml
   kubectl apply -f zookeeper.yaml
   kubectl apply -f kafka.yaml
   ```

3. **Deploy services**
   ```bash
   kubectl apply -f ../broker/deployment.yaml
   kubectl apply -f ../user/deployment.yaml
   kubectl apply -f ../product/deployment.yaml
   kubectl apply -f ../order/deployment.yaml
   kubectl apply -f ../payment/deployment.yaml
   ```

4. **Access via NodePort**
   ```bash
   minikube service broker-service --url
   ```

##  Environment Variables

### Broker Service
```env
JWT_SECRET=your_jwt_secret
USER_SERVICE_ADDR=user-service:50001
PRODUCT_SERVICE_ADDR=product-service:40001
ORDER_SERVICE_ADDR=order-service:30001
PAYMENT_SERVICE_ADDR=payment-service:60001
```

### User Service
```env
DATABASE_URL=postgres://user:password@postgres:5432/microservice?sslmode=disable
JWT_SECRET=your_jwt_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/oauth/google/callback
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URL=http://localhost:8080/api/v1/oauth/github/callback
FACEBOOK_CLIENT_ID=your_facebook_client_id
FACEBOOK_CLIENT_SECRET=your_facebook_client_secret
FACEBOOK_REDIRECT_URL=http://localhost:8080/api/v1/oauth/facebook/callback
```

### Product/Order/Payment Service
```env
DATABASE_URL=postgres://user:password@postgres:5432/microservice?sslmode=disable
KAFKA_BROKER=kafka:9092
```

##  Technologies Used

| Technology | Purpose |
|------------|---------|
| **Go** | Primary programming language |
| **Gin** | HTTP web framework (Broker) |
| **gRPC** | Inter-service communication |
| **Protocol Buffers** | Service definition & serialization |
| **PostgreSQL** | Primary database |
| **Apache Kafka** | Event streaming (order events) |
| **Docker** | Containerization |
| **Kubernetes** | Container orchestration |
| **JWT** | Authentication tokens |
| **OAuth2** | Social login (Google, Facebook, GitHub) |
| **Logrus** | Structured logging |

##  Features

- ✅ **Microservices Architecture** - Loosely coupled, independently deployable services
- ✅ **gRPC Communication** - High-performance RPC between services
- ✅ **REST API Gateway** - Single entry point for clients
- ✅ **JWT Authentication** - Secure token-based auth
- ✅ **OAuth2 Integration** - Social login support
- ✅ **Rate Limiting** - Protection against abuse
- ✅ **Database Migrations** - Version-controlled schema changes
- ✅ **Docker Support** - Containerized deployment
- ✅ **Kubernetes Ready** - Cloud-native deployment
- ✅ **Event-Driven** - Kafka for async communication
- ✅ **Structured Logging** - JSON-formatted logs

## 🧪 Testing

```bash
# Run tests for a specific service
cd user && go test ./...
cd product && go test ./...
cd order && go test ./...
cd payment && go test ./...
```

##  API Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### List Products
```bash
curl http://localhost:8080/product/?page=0
```

### Create Order (Protected)
```bash
curl -X POST http://localhost:8080/order/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "total_price": 99.99,
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 3, "quantity": 1}
    ]
  }'
```

### Process Payment (Protected)
```bash
curl -X POST http://localhost:8080/payment/transaction \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "payment_id": 1,
    "money": 100.00
  }'
```

##  Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**NumPyH**

---

⭐ Star this repo if you find it helpful!
