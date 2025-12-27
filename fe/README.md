# Frontend - E-Commerce Microservice

Frontend application untuk Go Microservice E-Commerce platform.

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Data Fetching**: React Query (TanStack Query)
- **Form Handling**: React Hook Form + Zod
- **HTTP Client**: Axios
- **Icons**: Lucide React

## Prerequisites

- Node.js 18+ 
- npm atau yarn
- Backend microservices running (broker service at port 8080)

## Installation

```bash
# Install dependencies
npm install

# Copy environment file
cp .env.example .env.local

# Edit .env.local dengan konfigurasi yang sesuai
```

## Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=Go Microservice Store
```

## Development

```bash
# Run development server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint
npm run lint
```

## Project Structure

```
src/
├── app/                    # Next.js App Router pages
│   ├── layout.tsx          # Root layout
│   ├── page.tsx            # Home page
│   ├── login/              # Login page
│   ├── register/           # Register page
│   ├── products/           # Product listing
│   ├── cart/               # Shopping cart
│   ├── checkout/           # Checkout page
│   ├── orders/             # Orders pages
│   │   ├── page.tsx        # Order history
│   │   └── [id]/           # Order detail
│   └── profile/            # User profile
├── components/
│   ├── ui/                 # Reusable UI components
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   ├── loading.tsx
│   │   └── product-card.tsx
│   └── layout/             # Layout components
│       ├── navbar.tsx
│       └── footer.tsx
├── lib/
│   ├── api.ts              # Axios instance with interceptors
│   ├── config.ts           # Environment config
│   └── utils.ts            # Utility functions
├── services/               # API service layers
│   ├── auth-service.ts
│   ├── product-service.ts
│   ├── order-service.ts
│   └── payment-service.ts
├── store/                  # Zustand stores
│   ├── auth-store.ts
│   ├── cart-store.ts
│   └── index.ts
└── types/                  # TypeScript types
    └── index.ts
```

## Features

### Authentication
- Login dengan email/password
- Register akun baru
- OAuth support (Google, Facebook, GitHub) - backend ready
- JWT token management dengan auto-refresh

### Products
- Product listing dengan pagination
- Add to cart functionality

### Shopping Cart
- Add/remove items
- Update quantities
- Persistent cart (localStorage)

### Orders
- Create order dari cart
- View order detail
- Order status tracking

### Payments
- Create payment untuk order
- Payment status tracking

## API Endpoints (Broker Service)

### Auth
- `POST /auth/login` - Login
- `POST /auth/register` - Register
- `POST /auth/logout` - Logout

### Products
- `GET /products` - Get all products
- `GET /products/:id` - Get product by ID

### Orders
- `POST /orders` - Create order
- `GET /orders/:id` - Get order by ID

### Payments
- `POST /payments` - Create payment
- `GET /payments/:id` - Get payment by ID

## Styling

Menggunakan Tailwind CSS dengan konfigurasi default. Warna utama:
- Primary: Blue (blue-600)
- Secondary: Gray
- Success: Green
- Warning: Yellow
- Error: Red

## State Management

### Auth Store
- `user` - Current user data
- `token` - JWT token
- `isAuthenticated` - Auth status
- `setAuth()` - Set auth data
- `logout()` - Clear auth data

### Cart Store
- `items` - Cart items
- `addItem()` - Add item to cart
- `removeItem()` - Remove item
- `updateQuantity()` - Update item quantity
- `clearCart()` - Clear all items
- `total` - Computed total price

## Notes

1. **Backend Connection**: Pastikan broker service berjalan di port 8080
2. **CORS**: Backend harus mengizinkan origin dari frontend
3. **Authentication**: Token disimpan di cookies dan localStorage
4. **Error Handling**: Global error handling via Axios interceptors

## Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## License

MIT
