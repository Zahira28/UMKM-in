'use client';

import React from 'react';
import { X } from 'lucide-react';
import { useApp } from '../../context/AppContext';
import PostCard from '../feed/PostCard';

export default function ProductDetailModal() {
  const { activeDetailProductId, closeDetailModal, products } = useApp();

  if (!activeDetailProductId) return null;

  const product = products.find((p) => p.id === activeDetailProductId);
  if (!product) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 overflow-y-auto animate-in fade-in duration-150">
      {/* Backdrop click */}
      <div className="fixed inset-0" onClick={closeDetailModal} />

      {/* Modal Content Container */}
      <div className="relative w-full max-w-2xl bg-white rounded-3xl p-4 md:p-6 shadow-2xl z-10 space-y-4 max-h-[90vh] overflow-y-auto my-auto animate-in zoom-in-95 duration-200">
        {/* Header matching Detail Produk design */}
        <div className="flex items-center justify-between border-b border-gray-100 pb-3 px-1">
          <h2 className="text-xl font-bold text-gray-900 tracking-tight">
            Detail produk
          </h2>
          <button
            onClick={closeDetailModal}
            className="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100 transition-colors"
            aria-label="Tutup detail produk"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* PostCard component */}
        <div className="pt-1">
          <PostCard product={product} />
        </div>
      </div>
    </div>
  );
}
