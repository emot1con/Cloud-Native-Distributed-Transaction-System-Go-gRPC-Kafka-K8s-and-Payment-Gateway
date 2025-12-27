import React from 'react';
import Link from 'next/link';
import { config } from '@/lib/config';

export function Footer() {
  return (
    <footer className="bg-gray-50 border-t">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center space-x-2 mb-4">
              <span className="text-2xl">🛒</span>
              <span className="font-bold text-xl text-gray-900">
                {config.appName}
              </span>
            </div>
            <p className="text-gray-600 text-sm max-w-md">
              {config.appDescription}. Built with Go Microservices, gRPC, Kafka, and Next.js.
            </p>
          </div>
          
          {/* Quick Links */}
          <div>
            <h3 className="font-semibold text-gray-900 mb-4">Quick Links</h3>
            <ul className="space-y-2">
              <li>
                <Link href="/" className="text-gray-600 hover:text-gray-900 text-sm">
                  Home
                </Link>
              </li>
              <li>
                <Link href="/products" className="text-gray-600 hover:text-gray-900 text-sm">
                  Products
                </Link>
              </li>
              <li>
                <Link href="/cart" className="text-gray-600 hover:text-gray-900 text-sm">
                  Cart
                </Link>
              </li>
            </ul>
          </div>
          
          {/* Account */}
          <div>
            <h3 className="font-semibold text-gray-900 mb-4">Account</h3>
            <ul className="space-y-2">
              <li>
                <Link href="/login" className="text-gray-600 hover:text-gray-900 text-sm">
                  Login
                </Link>
              </li>
              <li>
                <Link href="/register" className="text-gray-600 hover:text-gray-900 text-sm">
                  Register
                </Link>
              </li>
              <li>
                <Link href="/orders" className="text-gray-600 hover:text-gray-900 text-sm">
                  My Orders
                </Link>
              </li>
            </ul>
          </div>
        </div>
        
        <div className="border-t mt-8 pt-8 text-center text-sm text-gray-500">
          <p>&copy; {new Date().getFullYear()} {config.appName}. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}
