import { api } from '@/lib/api';
import { config } from '@/lib/config';
import type {
  Product,
  ProductPayload,
  ProductListResponse,
} from '@/types';

export const productService = {
  /**
   * Get paginated list of products
   */
  getProducts: async (page: number = 1): Promise<ProductListResponse> => {
    const response = await api.get(`${config.endpoints.products}?page=${page}`);
    return response.data;
  },

  /**
   * Create a new product (admin only)
   */
  createProduct: async (payload: ProductPayload): Promise<{ message: string }> => {
    const response = await api.post(config.endpoints.products, payload);
    return response.data;
  },

  /**
   * Update a product (admin only)
   */
  updateProduct: async (id: number, payload: Product): Promise<Product> => {
    const response = await api.put(`${config.endpoints.products}?id=${id}`, payload);
    return response.data;
  },

  /**
   * Delete a product (admin only)
   */
  deleteProduct: async (id: number): Promise<{ message: string }> => {
    const response = await api.delete(`${config.endpoints.products}?id=${id}`);
    return response.data;
  },
};
