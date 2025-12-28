import { api } from '@/lib/api';
import { config } from '@/lib/config';
import type { PaymentTransaction, Payment } from '@/types';

export const paymentService = {
  /**
   * Process payment transaction
   */
  processTransaction: async (payload: PaymentTransaction): Promise<{ message: string }> => {
    const response = await api.post(config.endpoints.paymentTransaction, payload);
    return response.data;
  },

  /**
   * Get payment by order ID
   */
  getPaymentByOrderId: async (orderId: number): Promise<Payment> => {
    const response = await api.get(`/payment/order/${orderId}`);
    return response.data;
  },
};
