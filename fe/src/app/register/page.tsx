'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import toast from 'react-hot-toast';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { authService } from '@/services/auth-service';

const registerSchema = z.object({
  full_name: z.string().min(3, 'Name must be at least 3 characters').max(100),
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters').max(100),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ['confirmPassword'],
});

type RegisterFormData = z.infer<typeof registerSchema>;

export default function RegisterPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  });
  
  const onSubmit = async (data: RegisterFormData) => {
    setIsLoading(true);
    try {
      await authService.register({
        full_name: data.full_name,
        email: data.email,
        password: data.password,
      });
      toast.success('Registration successful! Please login.');
      router.push('/login');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Registration failed');
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleOAuth = (provider: 'google' | 'facebook' | 'github') => {
    const urls = {
      google: authService.getGoogleOAuthUrl(),
      facebook: authService.getFacebookOAuthUrl(),
      github: authService.getGithubOAuthUrl(),
    };
    window.location.href = urls[provider];
  };
  
  return (
    <div className="min-h-[80vh] flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900">
            Create an account
          </h2>
          <p className="mt-2 text-gray-600">
            Join us today and start shopping
          </p>
        </div>
        
        <form onSubmit={handleSubmit(onSubmit)} className="mt-8 space-y-6">
          <div className="space-y-4">
            <Input
              label="Full Name"
              type="text"
              placeholder="John Doe"
              error={errors.full_name?.message}
              {...register('full_name')}
            />
            
            <Input
              label="Email"
              type="email"
              placeholder="you@example.com"
              error={errors.email?.message}
              {...register('email')}
            />
            
            <Input
              label="Password"
              type="password"
              placeholder="••••••••"
              error={errors.password?.message}
              {...register('password')}
            />
            
            <Input
              label="Confirm Password"
              type="password"
              placeholder="••••••••"
              error={errors.confirmPassword?.message}
              {...register('confirmPassword')}
            />
          </div>
          
          <Button
            type="submit"
            className="w-full"
            isLoading={isLoading}
          >
            Create Account
          </Button>
        </form>
        
        <div className="mt-6">
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300" />
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-2 bg-gray-50 text-gray-500">
                Or continue with
              </span>
            </div>
          </div>
          
          <div className="mt-6 grid grid-cols-3 gap-3">
            <Button
              type="button"
              variant="outline"
              onClick={() => handleOAuth('google')}
            >
              <span className="text-xl">🔵</span>
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => handleOAuth('facebook')}
            >
              <span className="text-xl">📘</span>
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => handleOAuth('github')}
            >
              <span className="text-xl">⚫</span>
            </Button>
          </div>
        </div>
        
        <p className="text-center text-sm text-gray-600">
          Already have an account?{' '}
          <Link href="/login" className="font-medium text-primary-600 hover:text-primary-500">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  );
}
