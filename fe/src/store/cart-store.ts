import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type { CartItem, Product } from '@/types';
import { config } from '@/lib/config';

interface CartState {
  items: CartItem[];
  
  // Computed
  totalItems: () => number;
  totalPrice: () => number;
  
  // Actions
  addItem: (product: Product, quantity?: number) => void;
  removeItem: (productId: number) => void;
  updateQuantity: (productId: number, quantity: number) => void;
  clearCart: () => void;
  getItem: (productId: number) => CartItem | undefined;
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      
      totalItems: () => {
        return get().items.reduce((total, item) => total + item.quantity, 0);
      },
      
      totalPrice: () => {
        return get().items.reduce(
          (total, item) => total + item.product.price * item.quantity,
          0
        );
      },
      
      addItem: (product: Product, quantity = 1) => {
        set((state) => {
          const existingItem = state.items.find(
            (item) => item.product.id === product.id
          );
          
          if (existingItem) {
            // Check stock limit
            const newQuantity = Math.min(
              existingItem.quantity + quantity,
              product.stock
            );
            
            return {
              items: state.items.map((item) =>
                item.product.id === product.id
                  ? { ...item, quantity: newQuantity }
                  : item
              ),
            };
          }
          
          // Add new item
          const addQuantity = Math.min(quantity, product.stock);
          return {
            items: [...state.items, { product, quantity: addQuantity }],
          };
        });
      },
      
      removeItem: (productId: number) => {
        set((state) => ({
          items: state.items.filter((item) => item.product.id !== productId),
        }));
      },
      
      updateQuantity: (productId: number, quantity: number) => {
        if (quantity <= 0) {
          get().removeItem(productId);
          return;
        }
        
        set((state) => ({
          items: state.items.map((item) => {
            if (item.product.id === productId) {
              const newQuantity = Math.min(quantity, item.product.stock);
              return { ...item, quantity: newQuantity };
            }
            return item;
          }),
        }));
      },
      
      clearCart: () => {
        set({ items: [] });
      },
      
      getItem: (productId: number) => {
        return get().items.find((item) => item.product.id === productId);
      },
    }),
    {
      name: config.storageKeys.cart,
      storage: createJSONStorage(() => localStorage),
    }
  )
);
