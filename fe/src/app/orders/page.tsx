'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Package, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useAuthStore } from '@/store';

export default function OrdersPage() {
  const router = useRouter();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const [orderId, setOrderId] = useState('');
  
  if (!isAuthenticated) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center">
        <Package className="w-16 h-16 text-gray-300 mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Please login first</h2>
        <p className="text-gray-600 mb-4">You need to be logged in to view your orders</p>
        <Link href="/login">
          <Button>Login</Button>
        </Link>
      </div>
    );
  }

  const handleSearchOrder = (e: React.FormEvent) => {
    e.preventDefault();
    if (orderId.trim()) {
      router.push(`/orders/${orderId}`);
    }
  };
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">My Orders</h1>
      
      <div className="bg-white rounded-lg shadow p-8">
        <div className="max-w-md mx-auto">
          <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-gray-900 mb-2 text-center">Find Your Order</h2>
          <p className="text-gray-600 mb-6 text-center">
            Enter your order ID to view order details and status.
          </p>
          
          <form onSubmit={handleSearchOrder} className="space-y-4">
            <Input
              label="Order ID"
              type="number"
              placeholder="Enter your order ID (e.g., 1, 2, 3...)"
              value={orderId}
              onChange={(e) => setOrderId(e.target.value)}
            />
            <Button type="submit" className="w-full">
              <Search className="w-4 h-4 mr-2" />
              View Order
            </Button>
          </form>
          
          <div className="mt-6 pt-6 border-t text-center">
            <p className="text-sm text-gray-500 mb-4">
              Don&apos;t have an order yet?
            </p>
            <Link href="/products">
              <Button variant="outline">Start Shopping</Button>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
