'use client';

import React, { useState, useMemo } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { CreditCard, CheckCircle2, ArrowLeft, Clock, XCircle, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Loading } from '@/components/ui/loading';
import { useAuthStore } from '@/store';
import { orderService } from '@/services/order-service';
import { paymentService } from '@/services/payment-service';
import { formatPrice, formatDate } from '@/lib/utils';
import type { OrderItem } from '@/types';

export default function OrderPaymentPage() {
  const params = useParams();
  const router = useRouter();
  const orderId = Number(params.id);
  const [isPaymentLoading, setIsPaymentLoading] = useState(false);
  const [isSnapOpen, setIsSnapOpen] = useState(false);
  
  // Get user info from auth store
  const user = useAuthStore((state) => state.user);
  
  const { data: ordersData, isLoading, error, refetch } = useQuery({
    queryKey: ['orders', 0],
    queryFn: () => orderService.getOrders(0),
    enabled: !isNaN(orderId),
    refetchOnMount: 'always',
    staleTime: 0,
  });

  const { data: paymentData, isLoading: isPaymentDataLoading, refetch: refetchPayment, isError: isPaymentError } = useQuery({
    queryKey: ['payment', orderId],
    queryFn: () => paymentService.getPaymentByOrderId(orderId),
    enabled: !isNaN(orderId),
    refetchOnMount: 'always',
    staleTime: 0,
    retry: 3, // Retry up to 3 times if payment not found (race condition with Kafka)
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 5000), // Exponential backoff: 1s, 2s, 4s
  });
  
  const order = useMemo(() => {
    if (!ordersData?.orders) return null;
    return ordersData.orders.find(o => o.id === orderId);
  }, [ordersData, orderId]);

  const handlePayWithMidtrans = async () => {
    // Prevent double-click and opening popup when already open
    if (isPaymentLoading || isSnapOpen) {
      return;
    }

    // If payment not loaded yet, try to refetch
    if (!paymentData?.id) {
      toast('Fetching payment info...', { icon: '⏳' });
      const result = await refetchPayment();
      if (!result.data?.id) {
        toast.error('Payment not found. Please wait a moment and try again.');
        return;
      }
    }

    // Get user info - use stored user or fallback
    const customerName = user?.full_name || user?.email?.split('@')[0] || 'Customer';
    const customerEmail = user?.email || '';
    
    if (!customerEmail) {
      toast.error('User information not available. Please login again.');
      return;
    }
    
    setIsPaymentLoading(true);
    
    // Get current payment data (might have been refetched)
    const currentPayment = paymentData || (await refetchPayment()).data;
    
    try {
      let token = currentPayment?.gateway_token;
      
      if (!token) {
        const response = await paymentService.initiatePayment({ 
          order_id: orderId,
          customer_name: customerName,
          customer_email: customerEmail,
        });
        token = response.token;
        
        if (!token) {
          throw new Error('Failed to get payment token');
        }
      }
      
      // Mark snap as open before calling
      setIsSnapOpen(true);
      setIsPaymentLoading(false);
      
      paymentService.openSnapPayment(token, {
        onSuccess: () => {
          setIsSnapOpen(false);
          toast.success('Payment successful! 🎉');
          refetch();
          refetchPayment();
        },
        onPending: () => {
          setIsSnapOpen(false);
          toast('Payment pending. Complete your payment.', { icon: '⏳' });
          refetchPayment();
        },
        onError: () => {
          setIsSnapOpen(false);
          toast.error('Payment failed. Please try again.');
        },
        onClose: () => {
          setIsSnapOpen(false);
          toast('Payment cancelled', { icon: '❌' });
          refetchPayment();
        },
      });
    } catch (err) {
      setIsSnapOpen(false);
      toast.error(err instanceof Error ? err.message : 'Failed to initiate payment');
    } finally {
      setIsPaymentLoading(false);
    }
  };
  
  if (isLoading || isPaymentDataLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <Loading size="lg" text="Loading order..." />
      </div>
    );
  }
  
  if (error || !order) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center">
        <p className="text-red-600 mb-4">Order not found</p>
        <Button onClick={() => router.push('/orders')}>
          Back to Orders
        </Button>
      </div>
    );
  }

  const renderStatusContent = () => {
    const status = order.status.toLowerCase();
    
    if (status === 'paid' || status === 'success') {
      return (
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <CheckCircle2 className="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Payment Successful!</h1>
          <p className="text-gray-600 mb-2">Order #{order.id} has been paid</p>
          <p className="text-2xl font-bold text-green-600 mb-2">{formatPrice(order.total_price)}</p>
          <p className="text-sm text-gray-500 mb-6">{formatDate(order.created_at)}</p>
          <Link href="/orders">
            <Button className="w-full">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back to My Orders
            </Button>
          </Link>
        </div>
      );
    }
    
    if (status === 'expired' || status === 'failed' || status === 'cancelled') {
      return (
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <XCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            Payment {status.charAt(0).toUpperCase() + status.slice(1)}
          </h1>
          <p className="text-gray-600 mb-6">Order #{order.id} - {formatPrice(order.total_price)}</p>
          <Link href="/products">
            <Button className="w-full">Continue Shopping</Button>
          </Link>
        </div>
      );
    }
    
    return (
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="bg-primary-600 text-white p-6">
          <div className="flex items-center gap-3 mb-2">
            <CreditCard className="w-6 h-6" />
            <h1 className="text-xl font-bold">Complete Payment</h1>
          </div>
          <p className="text-primary-100 text-sm">Order #{order.id}</p>
        </div>

        <div className="p-6 border-b">
          <h2 className="text-sm font-medium text-gray-500 uppercase mb-3">Order Summary</h2>
          
          {order.order_items && order.order_items.length > 0 ? (
            <div className="space-y-2 mb-4">
              {order.order_items.map((item: OrderItem) => (
                <div key={item.id} className="flex justify-between text-sm">
                  <span className="text-gray-600">Product #{item.product_id} × {item.quantity}</span>
                  <span className="font-medium">{formatPrice(item.price)}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-500 mb-4">Order items not available</p>
          )}

          <div className="flex justify-between pt-3 border-t">
            <span className="font-semibold text-gray-900">Total</span>
            <span className="font-bold text-xl text-primary-600">{formatPrice(order.total_price)}</span>
          </div>
        </div>

        <div className="p-6 border-b bg-gray-50">
          <div className="flex items-center gap-2 text-amber-600">
            <Clock className="w-5 h-5" />
            <span className="font-medium">Awaiting Payment</span>
          </div>
          {paymentData ? (
            <div className="mt-2 text-sm text-gray-600">
              <p>Payment ID: {paymentData.id}</p>
              {paymentData.expired_at && <p>Expires: {formatDate(paymentData.expired_at)}</p>}
            </div>
          ) : isPaymentDataLoading ? (
            <div className="mt-2 text-sm text-gray-500">
              <p>Loading payment info...</p>
            </div>
          ) : (
            <div className="mt-2 text-sm text-amber-600">
              <p>Payment is being prepared. Click button below to continue.</p>
            </div>
          )}
        </div>

        <div className="p-6">
          <Button className="w-full" size="lg" onClick={handlePayWithMidtrans} disabled={isPaymentLoading || isPaymentDataLoading || isSnapOpen}>
            {isPaymentLoading ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Processing...
              </>
            ) : isPaymentDataLoading ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Loading...
              </>
            ) : isSnapOpen ? (
              <>
                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                Payment in progress...
              </>
            ) : (
              <>
                <CreditCard className="w-5 h-5 mr-2" />
                Pay Now - {formatPrice(order.total_price)}
              </>
            )}
          </Button>
          <p className="text-center text-sm text-gray-500 mt-4">Secure payment powered by Midtrans</p>
        </div>
      </div>
    );
  };
  
  return (
    <div className="max-w-lg mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <Link href="/orders" className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900 mb-6">
        <ArrowLeft className="w-4 h-4 mr-1" />
        Back to My Orders
      </Link>
      {renderStatusContent()}
    </div>
  );
}
