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
  order_id: number;
  amount: number;
  currency: string;
  status: string;
  payment_method: string;
  gateway_order_id: string;
  gateway_token: string;
  gateway_redirect_url: string;
  gateway_transaction_id: string;
  gateway_status: string;
  va_number: string;
  qr_code_url: string;
  created_at: string;
  paid_at: string;
  expired_at: string;
}

export interface InitiatePaymentRequest {
  order_id: number;
  payment_method?: string;
  payment_channel?: string;
  customer_name: string;
  customer_email: string;
  customer_phone?: string;
}

export interface InitiatePaymentResponse {
  payment_id: number;
  token: string;
  redirect_url: string;
  va_number: string;
  qr_code_url: string;
  expired_at: string;
  status: string;
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
