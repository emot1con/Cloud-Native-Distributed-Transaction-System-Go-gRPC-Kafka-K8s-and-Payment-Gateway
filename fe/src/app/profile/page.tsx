'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { User, Mail, LogOut, Package, Calendar } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/store';
import { formatDate } from '@/lib/utils';
import Link from 'next/link';

export default function ProfilePage() {
  const router = useRouter();
  const { user, isAuthenticated, logout, role } = useAuthStore();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  if (!isAuthenticated || !user) {
    return (
      <div className="min-h-[60vh] flex flex-col items-center justify-center">
        <User className="w-16 h-16 text-gray-300 mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 mb-2">Please login first</h2>
        <p className="text-gray-600 mb-4">You need to be logged in to view your profile</p>
        <Link href="/login">
          <Button>Login</Button>
        </Link>
      </div>
    );
  }

  const handleLogout = async () => {
    setIsLoggingOut(true);
    try {
      logout();
      router.push('/');
    } finally {
      setIsLoggingOut(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">My Profile</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {/* Profile Card */}
        <div className="md:col-span-1">
          <div className="bg-white rounded-lg shadow p-6 text-center">
            <div className="w-24 h-24 bg-blue-100 rounded-full mx-auto mb-4 flex items-center justify-center">
              <User className="w-12 h-12 text-blue-600" />
            </div>
            <h2 className="text-xl font-semibold text-gray-900 mb-1">
              {user.full_name}
            </h2>
            <p className="text-gray-500 text-sm mb-2">{user.email}</p>
            {role && (
              <span className="inline-block px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded-full mb-4">
                {role}
              </span>
            )}
            <div className="space-y-2">
              <Link href="/orders" className="block">
                <Button variant="outline" className="w-full">
                  <Package className="w-4 h-4 mr-2" />
                  My Orders
                </Button>
              </Link>
              <Button
                variant="outline"
                className="w-full text-red-600 hover:bg-red-50"
                onClick={handleLogout}
                isLoading={isLoggingOut}
              >
                <LogOut className="w-4 h-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>
        </div>

        {/* Profile Details */}
        <div className="md:col-span-2">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-6">Account Information</h3>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Full Name
                </label>
                <div className="flex items-center space-x-2 p-3 bg-gray-50 rounded-md">
                  <User className="w-4 h-4 text-gray-400" />
                  <span>{user.full_name}</span>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <div className="flex items-center space-x-2 p-3 bg-gray-50 rounded-md">
                  <Mail className="w-4 h-4 text-gray-400" />
                  <span>{user.email}</span>
                </div>
              </div>

              {user.provider && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Login Provider
                  </label>
                  <div className="flex items-center space-x-2 p-3 bg-gray-50 rounded-md">
                    <span className="capitalize">{user.provider}</span>
                  </div>
                </div>
              )}

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Member Since
                </label>
                <div className="flex items-center space-x-2 p-3 bg-gray-50 rounded-md">
                  <Calendar className="w-4 h-4 text-gray-400" />
                  <span>{formatDate(user.created_at)}</span>
                </div>
              </div>
            </div>

            <div className="mt-6 pt-6 border-t">
              <p className="text-sm text-gray-500">
                Note: Profile editing functionality is not available yet. 
                Contact support to update your information.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
