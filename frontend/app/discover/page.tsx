'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Search, Heart, MessageCircle } from 'lucide-react';
import { useApp } from '../../src/context/AppContext';

export default function DiscoverPage() {
  const { products, openCommentModal } = useApp();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('Semua');

  const categories = ['Semua', 'Makanan', 'Minuman', 'Kriya', 'Pakaian'];

  const filteredProducts = products.filter((product) => {
    const matchesCategory =
      selectedCategory === 'Semua' || product.category === selectedCategory;
    const matchesSearch =
      product.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      product.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
      product.caption.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesCategory && matchesSearch;
  });

  return (
    <div className="space-y-8 max-w-6xl mx-auto">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight">
          Temukan yang lokal
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          Produk unik dan cerita baru dari pelaku UMKM di sekitarmu.
        </p>
      </div>

      {/* Search Bar & Category Filters */}
      <div className="space-y-4">
        <div className="relative">
          <Search className="w-5 h-5 text-gray-400 absolute left-4 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Cari produk atau orang..."
            className="w-full bg-white border border-gray-200 rounded-2xl pl-12 pr-4 py-3.5 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all shadow-xs"
          />
        </div>

        {/* Category Pills */}
        <div className="flex flex-wrap gap-2 pt-1">
          {categories.map((cat) => (
            <button
              key={cat}
              onClick={() => setSelectedCategory(cat)}
              className={`px-4 py-2 rounded-xl text-xs font-semibold transition-all ${
                selectedCategory === cat
                  ? 'bg-[#0d6e42] text-white shadow-xs'
                  : 'bg-white border border-gray-200 text-gray-700 hover:bg-gray-50'
              }`}
            >
              {cat}
            </button>
          ))}
        </div>
      </div>

      {/* Product Image Grid (3 Columns) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-5">
        {filteredProducts.map((product) => (
          <div
            key={product.id}
            className="group relative bg-white rounded-2xl overflow-hidden border border-gray-100 shadow-xs aspect-square cursor-pointer transition-all duration-300 hover:shadow-lg"
          >
            <img
              src={product.imageUrl}
              alt={product.title}
              className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            />

            {/* Hover Gradient Overlay */}
            <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-200 p-4 flex flex-col justify-between text-white">
              <div className="flex justify-between items-start">
                <span className="bg-white/20 backdrop-blur-xs text-white text-[11px] font-semibold px-2.5 py-1 rounded-full">
                  {product.category}
                </span>
              </div>

              <div className="space-y-2">
                <div>
                  <h4 className="font-bold text-sm line-clamp-1">{product.title}</h4>
                  <p className="text-xs text-white/80">@{product.username}</p>
                </div>

                <div className="flex items-center justify-between pt-1">
                  <div className="flex items-center gap-3 text-xs font-semibold">
                    <span className="flex items-center gap-1">
                      <Heart className="w-4 h-4 fill-white text-white" />
                      {product.likesCount}
                    </span>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        openCommentModal(product.id);
                      }}
                      className="flex items-center gap-1 hover:text-emerald-300 transition-colors"
                    >
                      <MessageCircle className="w-4 h-4" />
                      {product.commentsCount}
                    </button>
                  </div>

                  <Link
                    href={`/profile/${product.username}`}
                    className="text-xs font-bold text-[#10B981] bg-white px-3 py-1.5 rounded-xl hover:bg-emerald-50 transition-colors"
                  >
                    Lihat Produk
                  </Link>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
