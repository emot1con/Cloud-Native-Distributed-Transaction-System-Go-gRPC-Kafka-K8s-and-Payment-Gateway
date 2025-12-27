import { api } from '@/lib/api';
import { config } from '@/lib/config';
import type {
  CreateOrderRequest,
  OrderResponse,
  Order,
} from '@/types';

export const orderService = {
  /**
   * Create a new order
   */
  createOrder: async (payload: CreateOrderRequest): Promise<OrderResponse> => {
    const response = await api.post(config.endpoints.orders, payload);
    return response.data;
  },

  /**
   * Get order by ID
   */
  getOrder: async (orderId: number): Promise<Order> => {
    const response = await api.get(`${config.endpoints.orders}?order_id=${orderId}`);
    return response.data;
  },
};
