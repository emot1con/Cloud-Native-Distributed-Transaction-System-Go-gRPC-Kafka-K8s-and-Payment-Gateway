import { api } from '@/lib/api';
import { config } from '@/lib/config';
import type { PaymentTransaction } from '@/types';

export const paymentService = {
  /**
   * Process payment transaction
   */
  processTransaction: async (payload: PaymentTransaction): Promise<{ message: string }> => {
    const response = await api.post(config.endpoints.paymentTransaction, payload);
    return response.data;
  },
};
