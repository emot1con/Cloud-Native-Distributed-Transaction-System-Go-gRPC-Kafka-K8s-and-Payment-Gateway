'use client';

import React, { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Loading } from '@/components/ui/loading';
import { orderService } from '@/services/order-service';
import { paymentService } from '@/services/payment-service';
import { formatPrice, formatDate, getStatusColor, cn } from '@/lib/utils';
import type { OrderItem } from '@/types';

export default function OrderDetailPage() {
  const params = useParams();
  const router = useRouter();
  const orderId = Number(params.id);
  const [paymentId, setPaymentId] = useState('');
  const [money, setMoney] = useState('');
  const [isPaymentLoading, setIsPaymentLoading] = useState(false);
  
  const { data: order, isLoading, error, refetch } = useQuery({
    queryKey: ['order', orderId],
    queryFn: () => orderService.getOrder(orderId),
    enabled: !isNaN(orderId),
  });
  
  const handlePayment = async () => {
    if (!paymentId || !money) {
      toast.error('Please fill in payment details');
      return;
    }
    
    setIsPaymentLoading(true);
    try {
      await paymentService.processTransaction({
        payment_id: Number(paymentId),
        money: Number(money),
      });
      toast.success('Payment successful!');
      refetch();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Payment failed');
    } finally {
      setIsPaymentLoading(false);
    }
  };
  
  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <Loading size="lg" text="Loading order..." />
      </div>
    );
  }
  
  if (error || !order) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center">
        <p className="text-red-600 mb-4">Failed to load order</p>
        <Button onClick={() => router.push('/orders')}>
          Back to Orders
        </Button>
      </div>
    );
  }
  
  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold text-gray-900">
          Order #{order.id}
        </h1>
        <span className={cn(
          'px-3 py-1 rounded-full text-sm font-medium',
          getStatusColor(order.status)
        )}>
          {order.status}
        </span>
      </div>
      
      {/* Order Info */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Order Details</h2>
        
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p className="text-gray-500">Order ID</p>
            <p className="font-medium">#{order.id}</p>
          </div>
          <div>
            <p className="text-gray-500">Status</p>
            <p className="font-medium capitalize">{order.status}</p>
          </div>
          <div>
            <p className="text-gray-500">Total Price</p>
            <p className="font-medium text-primary-600">{formatPrice(order.total_price)}</p>
          </div>
          <div>
            <p className="text-gray-500">Created At</p>
            <p className="font-medium">{formatDate(order.created_at)}</p>
          </div>
        </div>
      </div>
      
      {/* Order Items */}
      {order.order_items && order.order_items.length > 0 && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Items</h2>
          
          <div className="divide-y">
            {order.order_items.map((item: OrderItem) => (
              <div key={item.id} className="py-3 flex justify-between">
                <div>
                  <p className="font-medium">Product #{item.product_id}</p>
                  <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
                </div>
                <p className="font-medium">{formatPrice(item.price)}</p>
              </div>
            ))}
          </div>
        </div>
      )}
      
      {/* Payment Form - Show only for pending orders */}
      {order.status === 'Pending' && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Make Payment</h2>
          <p className="text-sm text-gray-600 mb-4">
            Enter your payment ID and amount to complete the order.
          </p>
          
          <div className="space-y-4">
            <Input
              label="Payment ID"
              type="number"
              placeholder="Enter payment ID"
              value={paymentId}
              onChange={(e) => setPaymentId(e.target.value)}
            />
            <Input
              label="Amount"
              type="number"
              placeholder="Enter payment amount"
              value={money}
              onChange={(e) => setMoney(e.target.value)}
              helperText={`Order total: ${formatPrice(order.total_price)}`}
            />
            <Button
              className="w-full"
              onClick={handlePayment}
              isLoading={isPaymentLoading}
            >
              Complete Payment
            </Button>
          </div>
        </div>
      )}
      
      <Link href="/orders">
        <Button variant="outline" className="w-full">
          Back to Orders
        </Button>
      </Link>
    </div>
  );
}
