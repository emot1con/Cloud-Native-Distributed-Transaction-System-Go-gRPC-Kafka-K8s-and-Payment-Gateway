import React from 'react';
import { Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';

interface LoadingProps {
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  text?: string;
}

export function Loading({ size = 'md', className, text }: LoadingProps) {
  const sizes = {
    sm: 'w-4 h-4',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  };
  
  return (
    <div className={cn('flex flex-col items-center justify-center', className)}>
      <Loader2 className={cn(sizes[size], 'animate-spin text-primary-600')} />
      {text && <p className="mt-2 text-sm text-gray-500">{text}</p>}
    </div>
  );
}

interface LoadingPageProps {
  text?: string;
}

export function LoadingPage({ text = 'Loading...' }: LoadingPageProps) {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <Loading size="lg" text={text} />
    </div>
  );
}
