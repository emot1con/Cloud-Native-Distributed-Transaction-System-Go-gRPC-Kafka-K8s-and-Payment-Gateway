import { api } from '@/lib/api';
import type { Payment, InitiatePaymentRequest, InitiatePaymentResponse } from '@/types';

export const paymentService = {
  /**
   * Get payment by order ID
   */
  getPaymentByOrderId: async (orderId: number): Promise<Payment> => {
    const response = await api.get(`/payment/order/${orderId}`);
    return response.data;
  },

  /**
   * Initiate payment - Get Snap token from backend
   */
  initiatePayment: async (payload: InitiatePaymentRequest): Promise<InitiatePaymentResponse> => {
    const response = await api.post('/payment/initiate', payload);
    return response.data;
  },

  /**
   * Open Midtrans Snap payment popup
   */
  openSnapPayment: (
    token: string,
    callbacks?: {
      onSuccess?: (result: unknown) => void;
      onPending?: (result: unknown) => void;
      onError?: (result: unknown) => void;
      onClose?: () => void;
    }
  ): void => {
    // Type assertion for Snap global object
    const snap = (window as unknown as { snap?: { pay: (token: string, options: object) => void } }).snap;
    
    if (!snap) {
      console.error('Midtrans Snap.js not loaded');
      throw new Error('Payment system not available. Please refresh the page.');
    }

    snap.pay(token, {
      onSuccess: (result: unknown) => {
        console.log('Payment success:', result);
        callbacks?.onSuccess?.(result);
      },
      onPending: (result: unknown) => {
        console.log('Payment pending:', result);
        callbacks?.onPending?.(result);
      },
      onError: (result: unknown) => {
        console.error('Payment error:', result);
        callbacks?.onError?.(result);
      },
      onClose: () => {
        console.log('Payment popup closed');
        callbacks?.onClose?.();
      },
    });
  },
};
