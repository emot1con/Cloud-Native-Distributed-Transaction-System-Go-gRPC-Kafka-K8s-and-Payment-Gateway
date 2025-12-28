'use client';

import React, { useState, useMemo } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { CreditCard, CheckCircle2, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loading } from '@/components/ui/loading';
import { orderService } from '@/services/order-service';
import { paymentService } from '@/services/payment-service';
import { formatPrice, formatDate } from '@/lib/utils';
import type { OrderItem } from '@/types';

export default function OrderPaymentPage() {
  const params = useParams();
  const router = useRouter();
  const orderId = Number(params.id);
  const [money, setMoney] = useState('');
  const [isPaymentLoading, setIsPaymentLoading] = useState(false);
  
  // Fetch all orders and find the specific one
  const { data: ordersData, isLoading, error, refetch } = useQuery({
    queryKey: ['orders', 0],
    queryFn: () => orderService.getOrders(0),
    enabled: !isNaN(orderId),
    refetchOnMount: 'always',
    staleTime: 0,
  });

  // Fetch payment info for this order
  const { data: paymentData, isLoading: isPaymentDataLoading } = useQuery({
    queryKey: ['payment', orderId],
    queryFn: () => paymentService.getPaymentByOrderId(orderId),
    enabled: !isNaN(orderId),
    refetchOnMount: 'always',
    staleTime: 0,
  });
  
  // Find the specific order from the list
  const order = useMemo(() => {
    if (!ordersData?.orders) return null;
    return ordersData.orders.find(o => o.id === orderId);
  }, [ordersData, orderId]);

  const handlePayment = async () => {
    if (!money) {
      toast.error('Please enter the payment amount');
      return;
    }

    if (!paymentData?.id) {
      toast.error('Payment not found for this order');
      return;
    }
    
    setIsPaymentLoading(true);
    try {
      await paymentService.processTransaction({
        payment_id: paymentData.id,
        money: Number(money),
      });
      toast.success('Payment successful! Your order is being processed.');
      refetch();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Payment failed');
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

  // If order is already paid, show success message
  if (order.status.toLowerCase() !== 'pending') {
    return (
      <div className="max-w-lg mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <CheckCircle2 className="w-16 h-16 text-green-500 mx-auto mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            Order #{order.id} - {order.status}
          </h1>
          <p className="text-gray-600 mb-2">
            Total: {formatPrice(order.total_price)}
          </p>
          <p className="text-sm text-gray-500 mb-6">
            {formatDate(order.created_at)}
          </p>
          <Link href="/orders">
            <Button className="w-full">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back to My Orders
            </Button>
          </Link>
        </div>
      </div>
    );
  }
  
  return (
    <div className="max-w-lg mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <Link href="/orders" className="inline-flex items-center text-sm text-gray-600 hover:text-gray-900 mb-6">
        <ArrowLeft className="w-4 h-4 mr-1" />
        Back to My Orders
      </Link>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        {/* Header */}
        <div className="bg-primary-600 text-white p-6">
          <div className="flex items-center gap-3 mb-2">
            <CreditCard className="w-6 h-6" />
            <h1 className="text-xl font-bold">Complete Payment</h1>
          </div>
          <p className="text-primary-100 text-sm">Order #{order.id}</p>
        </div>

        {/* Order Summary */}
        <div className="p-6 border-b">
          <h2 className="text-sm font-medium text-gray-500 uppercase mb-3">Order Summary</h2>
          
          {order.order_items && order.order_items.length > 0 ? (
            <div className="space-y-2 mb-4">
              {order.order_items.map((item: OrderItem) => (
                <div key={item.id} className="flex justify-between text-sm">
                  <span className="text-gray-600">
                    Product #{item.product_id} × {item.quantity}
                  </span>
                  <span className="font-medium">{formatPrice(item.price)}</span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-500 mb-4">Order items not available</p>
          )}

          <div className="flex justify-between pt-3 border-t">
            <span className="font-semibold text-gray-900">Total</span>
            <span className="font-bold text-xl text-primary-600">
              {formatPrice(order.total_price)}
            </span>
          </div>
        </div>

        {/* Payment Form */}
        <div className="p-6">
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
            <p className="text-sm text-blue-800">
              <strong>Payment ID:</strong> {paymentData?.id || 'Loading...'}<br />
              <strong>Status:</strong> {paymentData?.status || 'Loading...'}
            </p>
          </div>

          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
            <p className="text-sm text-yellow-800">
              <strong>💡 Tip:</strong> Enter the exact amount shown above ({formatPrice(order.total_price)}) to complete payment.
            </p>
          </div>

          <div className="space-y-4">
            <Input
              label="Payment Amount (IDR)"
              type="number"
              placeholder={String(order.total_price)}
              value={money}
              onChange={(e) => setMoney(e.target.value)}
            />
            <Button
              className="w-full"
              size="lg"
              onClick={handlePayment}
              isLoading={isPaymentLoading}
              disabled={!paymentData?.id}
            >
              <CreditCard className="w-5 h-5 mr-2" />
              Pay {formatPrice(order.total_price)}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
