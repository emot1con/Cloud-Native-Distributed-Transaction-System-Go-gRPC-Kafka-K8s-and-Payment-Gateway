import { api } from '@/lib/api';
import { config } from '@/lib/config';
import type {
  CreateOrderRequest,
  OrderResponse,
  OrdersResponse,
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
   * Get orders for current user (uses JWT token for user identification)
   * @param offset - Pagination offset (default: 0)
   */
  getOrders: async (offset: number = 0): Promise<OrdersResponse> => {
    const response = await api.get(`${config.endpoints.orders}?offset=${offset}`);
    return response.data;
  },
};
