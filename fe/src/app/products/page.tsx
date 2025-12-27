'use client';

import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { ProductCard } from '@/components/ui/product-card';
import { Loading } from '@/components/ui/loading';
import { Button } from '@/components/ui/button';
import { productService } from '@/services/product-service';
import { useCartStore, useAuthStore } from '@/store';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import type { Product } from '@/types';

export default function ProductsPage() {
  const [page, setPage] = useState(1);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const addItem = useCartStore((state) => state.addItem);
  const getItem = useCartStore((state) => state.getItem);
  
  const { data, isLoading, error } = useQuery({
    queryKey: ['products', page],
    queryFn: () => productService.getProducts(page),
  });
  
  const handleAddToCart = (product: Product) => {
    if (!isAuthenticated) {
      toast.error('Please login to add items to cart');
      return;
    }
    addItem(product);
    toast.success(`${product.name} added to cart`);
  };
  
  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <Loading size="lg" text="Loading products..." />
      </div>
    );
  }
  
  if (error) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">Failed to load products</p>
          <Button onClick={() => window.location.reload()}>
            Try Again
          </Button>
        </div>
      </div>
    );
  }
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Products</h1>
        <p className="text-gray-600 mt-2">
          Browse our collection of {data?.total || 0} products
        </p>
      </div>
      
      {data?.products && data.products.length > 0 ? (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {data.products.map((product: Product) => (
              <ProductCard
                key={product.id}
                product={product}
                onAddToCart={handleAddToCart}
                isInCart={!!getItem(product.id)}
              />
            ))}
          </div>
          
          {/* Pagination */}
          {data.total_page > 1 && (
            <div className="mt-8 flex items-center justify-center gap-4">
              <Button
                variant="outline"
                disabled={page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
              >
                <ChevronLeft className="w-4 h-4 mr-1" />
                Previous
              </Button>
              
              <span className="text-gray-600">
                Page {data.page} of {data.total_page}
              </span>
              
              <Button
                variant="outline"
                disabled={page >= data.total_page}
                onClick={() => setPage((p) => p + 1)}
              >
                Next
                <ChevronRight className="w-4 h-4 ml-1" />
              </Button>
            </div>
          )}
        </>
      ) : (
        <div className="text-center py-12">
          <p className="text-gray-600">No products available</p>
        </div>
      )}
    </div>
  );
}
