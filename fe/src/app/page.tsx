import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { config } from '@/lib/config';

export default function HomePage() {
  return (
    <div>
      {/* Hero Section */}
      <section className="bg-gradient-to-br from-primary-600 to-primary-800 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold mb-6">
              Welcome to {config.appName}
            </h1>
            <p className="text-xl md:text-2xl text-primary-100 mb-8 max-w-2xl mx-auto">
              {config.appDescription}
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link href="/products">
                <Button size="lg" className="bg-white text-primary-600 hover:bg-gray-100">
                  Browse Products
                </Button>
              </Link>
              <Link href="/register">
                <Button size="lg" variant="outline" className="border-white text-white hover:bg-white/10">
                  Get Started
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>
      
      {/* Features Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-3xl font-bold text-center mb-12">
            Why Choose Us?
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* Feature 1 */}
            <div className="text-center p-6">
              <div className="text-4xl mb-4">🚀</div>
              <h3 className="text-xl font-semibold mb-2">Fast & Reliable</h3>
              <p className="text-gray-600">
                Built with microservices architecture for maximum scalability and reliability.
              </p>
            </div>
            
            {/* Feature 2 */}
            <div className="text-center p-6">
              <div className="text-4xl mb-4">🔒</div>
              <h3 className="text-xl font-semibold mb-2">Secure Payments</h3>
              <p className="text-gray-600">
                Your transactions are protected with industry-standard security measures.
              </p>
            </div>
            
            {/* Feature 3 */}
            <div className="text-center p-6">
              <div className="text-4xl mb-4">📦</div>
              <h3 className="text-xl font-semibold mb-2">Wide Selection</h3>
              <p className="text-gray-600">
                Browse through our extensive catalog of quality products.
              </p>
            </div>
          </div>
        </div>
      </section>
      
      {/* CTA Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl font-bold mb-4">Ready to Start Shopping?</h2>
          <p className="text-gray-600 mb-8 max-w-xl mx-auto">
            Create an account today and get access to exclusive deals and offers.
          </p>
          <Link href="/register">
            <Button size="lg">Create Account</Button>
          </Link>
        </div>
      </section>
    </div>
  );
}
