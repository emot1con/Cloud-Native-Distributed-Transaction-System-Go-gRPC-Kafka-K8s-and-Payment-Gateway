'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ShoppingCart, User, LogOut, Menu, X } from 'lucide-react';
import { useState } from 'react';
import { cn } from '@/lib/utils';
import { useAuthStore, useCartStore } from '@/store';
import { Button } from '@/components/ui/button';
import { config } from '@/lib/config';

export function Navbar() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const pathname = usePathname();
  const { isAuthenticated, logout } = useAuthStore();
  const cartItems = useCartStore((state) => state.items);
  const totalItems = cartItems.reduce((sum, item) => sum + item.quantity, 0);
  
  const navLinks = [
    { href: '/', label: 'Home' },
    { href: '/products', label: 'Products' },
  ];
  
  const authLinks = isAuthenticated
    ? [
        { href: '/orders', label: 'My Orders' },
      ]
    : [];
  
  const allLinks = [...navLinks, ...authLinks];
  
  return (
    <nav className="bg-white shadow-sm border-b sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-2">
              <span className="text-2xl">🛒</span>
              <span className="font-bold text-xl text-gray-900">
                {config.appName}
              </span>
            </Link>
          </div>
          
          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            {allLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={cn(
                  'text-sm font-medium transition-colors',
                  pathname === link.href
                    ? 'text-primary-600'
                    : 'text-gray-600 hover:text-gray-900'
                )}
              >
                {link.label}
              </Link>
            ))}
          </div>
          
          {/* Right side actions */}
          <div className="hidden md:flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                {/* Cart */}
                <Link href="/cart" className="relative p-2 text-gray-600 hover:text-gray-900">
                  <ShoppingCart className="w-6 h-6" />
                  {totalItems > 0 && (
                    <span className="absolute -top-1 -right-1 bg-primary-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                      {totalItems > 99 ? '99+' : totalItems}
                    </span>
                  )}
                </Link>
                
                {/* Logout */}
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={logout}
                  className="text-gray-600"
                >
                  <LogOut className="w-4 h-4 mr-2" />
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Link href="/login">
                  <Button variant="ghost" size="sm">
                    Login
                  </Button>
                </Link>
                <Link href="/register">
                  <Button variant="primary" size="sm">
                    Register
                  </Button>
                </Link>
              </>
            )}
          </div>
          
          {/* Mobile menu button */}
          <div className="md:hidden flex items-center">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="p-2 text-gray-600"
            >
              {isMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>
          </div>
        </div>
      </div>
      
      {/* Mobile menu */}
      {isMenuOpen && (
        <div className="md:hidden bg-white border-t">
          <div className="px-4 py-4 space-y-3">
            {allLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={cn(
                  'block px-3 py-2 rounded-md text-base font-medium',
                  pathname === link.href
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-600 hover:bg-gray-50'
                )}
                onClick={() => setIsMenuOpen(false)}
              >
                {link.label}
              </Link>
            ))}
            
            <hr className="my-2" />
            
            {isAuthenticated ? (
              <>
                <Link
                  href="/cart"
                  className="flex items-center px-3 py-2 text-gray-600"
                  onClick={() => setIsMenuOpen(false)}
                >
                  <ShoppingCart className="w-5 h-5 mr-2" />
                  Cart ({totalItems})
                </Link>
                <button
                  onClick={() => {
                    logout();
                    setIsMenuOpen(false);
                  }}
                  className="flex items-center w-full px-3 py-2 text-gray-600"
                >
                  <LogOut className="w-5 h-5 mr-2" />
                  Logout
                </button>
              </>
            ) : (
              <div className="space-y-2">
                <Link href="/login" onClick={() => setIsMenuOpen(false)}>
                  <Button variant="outline" className="w-full">Login</Button>
                </Link>
                <Link href="/register" onClick={() => setIsMenuOpen(false)}>
                  <Button variant="primary" className="w-full">Register</Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      )}
    </nav>
  );
}
