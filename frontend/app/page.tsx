'use client';

import React from 'react';
import { useApp } from '../src/context/AppContext';
import PostCard from '../src/components/feed/PostCard';

export default function HomePage() {
  const { products } = useApp();

  return (
    <div className="space-y-6">
      {/* Page Header matching Beranda design */}
      <div className="max-w-2xl mx-auto">
        <span className="text-xs font-bold uppercase tracking-wider text-[#0d6e42]">
          KOMUNITAS HARI INI
        </span>
        <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight mt-0.5">
          Cerita produk lokal
        </h1>
      </div>

      {/* Feed List */}
      <div className="space-y-6">
        {products.map((product) => (
          <PostCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  );
}
