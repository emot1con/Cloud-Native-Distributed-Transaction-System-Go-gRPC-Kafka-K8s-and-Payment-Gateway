import React from 'react';
import { cn, formatPrice } from '@/lib/utils';
import { Button } from './button';
import { ShoppingCart } from 'lucide-react';
import type { Product } from '@/types';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product) => void;
  isInCart?: boolean;
}

export function ProductCard({ product, onAddToCart, isInCart }: ProductCardProps) {
  const isOutOfStock = product.stock <= 0;
  
  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
      {/* Product Image Placeholder */}
      <div className="h-48 bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
        <span className="text-4xl">📦</span>
      </div>
      
      {/* Product Info */}
      <div className="p-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-1 line-clamp-1">
          {product.name}
        </h3>
        <p className="text-sm text-gray-500 mb-3 line-clamp-2 h-10">
          {product.description || 'No description available'}
        </p>
        
        <div className="flex items-center justify-between mb-3">
          <span className="text-xl font-bold text-primary-600">
            {formatPrice(product.price)}
          </span>
          <span className={cn(
            'text-sm px-2 py-1 rounded-full',
            isOutOfStock
              ? 'bg-red-100 text-red-600'
              : product.stock < 10
                ? 'bg-yellow-100 text-yellow-600'
                : 'bg-green-100 text-green-600'
          )}>
            {isOutOfStock ? 'Out of Stock' : `Stock: ${product.stock}`}
          </span>
        </div>
        
        <Button
          variant={isInCart ? 'secondary' : 'primary'}
          className="w-full"
          disabled={isOutOfStock}
          onClick={() => onAddToCart?.(product)}
        >
          <ShoppingCart className="w-4 h-4 mr-2" />
          {isInCart ? 'In Cart' : 'Add to Cart'}
        </Button>
      </div>
    </div>
  );
}
