'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import toast from 'react-hot-toast';
import { Button } from '@/components/ui/button';
import { useCartStore, useAuthStore } from '@/store';
import { orderService } from '@/services/order-service';
import { formatPrice } from '@/lib/utils';
import type { CreateOrderRequest } from '@/types';

export default function CheckoutPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const items = useCartStore((state) => state.items);
  const totalPrice = useCartStore((state) => state.totalPrice);
  const clearCart = useCartStore((state) => state.clearCart);
  
  if (!isAuthenticated) {
    router.push('/login');
    return null;
  }
  
  if (items.length === 0) {
    router.push('/cart');
    return null;
  }
  
  const handlePlaceOrder = async () => {
    setIsLoading(true);
    try {
      const orderPayload: CreateOrderRequest = {
        items: items.map((item) => ({
          product_id: item.product.id,
          quantity: item.quantity,
        })),
      };
      
      const response = await orderService.createOrder(orderPayload);
      
      clearCart();
      toast.success('Order placed! Redirecting to payment...');
      
      // Redirect to order detail page where user can pay
      router.push(`/orders/${response.order.id}`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to place order');
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Checkout</h1>
      
      {/* Order Summary */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Order Summary</h2>
        
        <div className="divide-y">
          {items.map((item) => (
            <div key={item.product.id} className="py-4 flex justify-between">
              <div>
                <p className="font-medium text-gray-900">{item.product.name}</p>
                <p className="text-sm text-gray-500">
                  {formatPrice(item.product.price)} × {item.quantity}
                </p>
              </div>
              <p className="font-medium text-gray-900">
                {formatPrice(item.product.price * item.quantity)}
              </p>
            </div>
          ))}
        </div>
        
        <hr className="my-4" />
        
        <div className="flex justify-between text-lg font-semibold">
          <span>Total</span>
          <span>{formatPrice(totalPrice())}</span>
        </div>
      </div>
      
      {/* Payment Info */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Payment Information</h2>
        <p className="text-gray-600 text-sm">
          After placing your order, you will be redirected to complete payment 
          using our secure payment gateway (Midtrans). You can pay using various methods 
          including Credit Card, Bank Transfer, GoPay, OVO, and more.
        </p>
      </div>
      
      {/* Actions */}
      <div className="flex gap-4">
        <Button
          variant="outline"
          className="flex-1"
          onClick={() => router.back()}
        >
          Back to Cart
        </Button>
        <Button
          className="flex-1"
          onClick={handlePlaceOrder}
          isLoading={isLoading}
        >
          Place Order
        </Button>
      </div>
    </div>
  );
}
