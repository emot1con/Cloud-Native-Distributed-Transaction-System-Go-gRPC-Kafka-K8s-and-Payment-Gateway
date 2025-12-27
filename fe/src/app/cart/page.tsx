'use client';

import React from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Trash2, Plus, Minus, ShoppingBag } from 'lucide-react';
import toast from 'react-hot-toast';
import { Button } from '@/components/ui/button';
import { useCartStore, useAuthStore } from '@/store';
import { formatPrice } from '@/lib/utils';

export default function CartPage() {
  const router = useRouter();
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const items = useCartStore((state) => state.items);
  const removeItem = useCartStore((state) => state.removeItem);
  const updateQuantity = useCartStore((state) => state.updateQuantity);
  const totalPrice = useCartStore((state) => state.totalPrice);
  
  if (!isAuthenticated) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center">
        <ShoppingBag className="w-16 h-16 text-gray-300 mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Please login first</h2>
        <p className="text-gray-600 mb-4">You need to be logged in to view your cart</p>
        <Link href="/login">
          <Button>Login</Button>
        </Link>
      </div>
    );
  }
  
  if (items.length === 0) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center">
        <ShoppingBag className="w-16 h-16 text-gray-300 mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Your cart is empty</h2>
        <p className="text-gray-600 mb-4">Add some products to get started</p>
        <Link href="/products">
          <Button>Browse Products</Button>
        </Link>
      </div>
    );
  }
  
  const handleCheckout = () => {
    router.push('/checkout');
  };
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Shopping Cart</h1>
      
      <div className="lg:grid lg:grid-cols-12 lg:gap-8">
        {/* Cart Items */}
        <div className="lg:col-span-8">
          <div className="bg-white rounded-lg shadow divide-y">
            {items.map((item) => (
              <div key={item.product.id} className="p-6 flex items-center gap-4">
                {/* Product Image */}
                <div className="w-20 h-20 bg-gray-100 rounded-lg flex items-center justify-center flex-shrink-0">
                  <span className="text-2xl">📦</span>
                </div>
                
                {/* Product Info */}
                <div className="flex-1 min-w-0">
                  <h3 className="text-lg font-medium text-gray-900 truncate">
                    {item.product.name}
                  </h3>
                  <p className="text-sm text-gray-500">
                    {formatPrice(item.product.price)} each
                  </p>
                </div>
                
                {/* Quantity Controls */}
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => updateQuantity(item.product.id, item.quantity - 1)}
                    className="p-1 rounded-md hover:bg-gray-100"
                  >
                    <Minus className="w-4 h-4" />
                  </button>
                  <span className="w-8 text-center font-medium">{item.quantity}</span>
                  <button
                    onClick={() => {
                      if (item.quantity < item.product.stock) {
                        updateQuantity(item.product.id, item.quantity + 1);
                      } else {
                        toast.error('Maximum stock reached');
                      }
                    }}
                    className="p-1 rounded-md hover:bg-gray-100"
                  >
                    <Plus className="w-4 h-4" />
                  </button>
                </div>
                
                {/* Item Total */}
                <div className="text-right">
                  <p className="font-semibold text-gray-900">
                    {formatPrice(item.product.price * item.quantity)}
                  </p>
                </div>
                
                {/* Remove Button */}
                <button
                  onClick={() => {
                    removeItem(item.product.id);
                    toast.success('Item removed from cart');
                  }}
                  className="p-2 text-red-500 hover:bg-red-50 rounded-md"
                >
                  <Trash2 className="w-5 h-5" />
                </button>
              </div>
            ))}
          </div>
        </div>
        
        {/* Order Summary */}
        <div className="lg:col-span-4 mt-8 lg:mt-0">
          <div className="bg-white rounded-lg shadow p-6 sticky top-24">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Order Summary</h2>
            
            <div className="space-y-3 mb-6">
              <div className="flex justify-between text-gray-600">
                <span>Subtotal ({items.reduce((s, i) => s + i.quantity, 0)} items)</span>
                <span>{formatPrice(totalPrice())}</span>
              </div>
              <div className="flex justify-between text-gray-600">
                <span>Shipping</span>
                <span className="text-green-600">Free</span>
              </div>
              <hr />
              <div className="flex justify-between text-lg font-semibold">
                <span>Total</span>
                <span>{formatPrice(totalPrice())}</span>
              </div>
            </div>
            
            <Button className="w-full" onClick={handleCheckout}>
              Proceed to Checkout
            </Button>
            
            <Link href="/products" className="block mt-4">
              <Button variant="outline" className="w-full">
                Continue Shopping
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
