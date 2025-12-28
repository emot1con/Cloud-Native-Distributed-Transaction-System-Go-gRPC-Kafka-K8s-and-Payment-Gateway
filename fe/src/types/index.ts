// API Response Types - matches Go proto definitions

// ============ User / Auth ============
export interface User {
  id: number;
  full_name: string;
  email: string;
  provider?: string;
  provider_id?: string;
  created_at: string;
  updated_at: string;
}

export interface RegisterPayload {
  full_name: string;
  email: string;
  password: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface TokenResponse {
  message: string;
  token: string;
  expired_at: string;
  refresh_token: string;
  refresh_token_expired_at: string;
  role: string;
}

// ============ Product ============
export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  stock: number;
  created_at: string;
  updated_at: string;
}

export interface ProductPayload {
  name: string;
  description: string;
  price: number;
  stock: number;
}

export interface ProductListResponse {
  total: number;
  page: number;
  total_page: number;
  products: Product[];
}

// ============ Order ============
export interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
  quantity: number;
  price: number;
}

export interface Order {
  id: number;
  user_id: number;
  total_price: number;
  status: string; // 'pending' | 'paid' | 'failed' - case insensitive
  created_at: string;
  updated_at: string;
  order_items?: OrderItem[];
}

export interface OrderItemRequest {
  product_id: number;
  quantity: number;
}

export interface CreateOrderRequest {
  items: OrderItemRequest[];
}

export interface OrderResponse {
  order: Order;
}

export interface OrdersResponse {
  orders: Order[];
}

// ============ Payment ============
export interface Payment {
  id: number;
  user_id: number;
  order_id: number;
  status: 'pending' | 'paid';
  total_price: number;
  created_at: string;
  updated_at: string;
}

export interface PaymentTransaction {
  payment_id: number;
  money: number;
}

// ============ Cart (Frontend only) ============
export interface CartItem {
  product: Product;
  quantity: number;
}

// ============ API Error ============
export interface ApiError {
  error: string;
}
