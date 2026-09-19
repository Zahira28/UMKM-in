'use client';

import React, { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import { Heart, MessageCircle, MoreVertical, Edit3, Trash2 } from 'lucide-react';
import { Product } from '../../types';
import { useApp } from '../../context/AppContext';

interface PostCardProps {
  product: Product;
}

export default function PostCard({ product }: PostCardProps) {
  const { likesMap, toggleLike, openCommentModal, openDeleteModal, currentUser } = useApp();
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  const isLiked = !!likesMap[product.id];
  const isOwner = product.username === currentUser.username;

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setMenuOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleWhatsAppOrder = () => {
    const cleanNumber = (product.userId === 'u1' ? currentUser.whatsappNumber : '+6281234567890').replace(/[^0-9]/g, '');
    const message = encodeURIComponent(
      `Halo *${product.username}*, saya tertarik membeli *${product.title}* (Rp ${product.price.toLocaleString('id-ID')}) yang saya lihat di UMKM-in! Apakah stok masih tersedia?`
    );
    window.open(`https://wa.me/${cleanNumber}?text=${message}`, '_blank');
  };

  return (
    <article className="bg-white rounded-2xl border border-gray-100 shadow-xs max-w-2xl mx-auto overflow-hidden transition-all duration-200 hover:shadow-md">
      {/* Header */}
      <header className="flex items-center justify-between p-4 bg-white">
        <div className="flex items-center gap-3">
          <Link href={`/profile/${product.username}`} className="flex-shrink-0">
            <div className="w-10 h-10 rounded-full overflow-hidden bg-gray-100 border border-gray-200">
              <img
                src={product.userAvatar || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=150&q=80'}
                alt={product.username}
                className="w-full h-full object-cover"
              />
            </div>
          </Link>
          <div>
            <div className="flex items-center gap-2">
              <Link
                href={`/profile/${product.username}`}
                className="font-bold text-sm text-gray-900 hover:underline"
              >
                {product.username}
              </Link>
              <span className="bg-[#e6f4ea] text-[#0d6e42] text-[11px] font-semibold px-2.5 py-0.5 rounded-full">
                {product.category}
              </span>
            </div>
            <p className="text-xs text-gray-500">{product.userLocation}</p>
          </div>
        </div>

        {/* 3 Dots Menu Dropdown */}
        <div className="relative" ref={menuRef}>
          <button
            onClick={() => setMenuOpen(!menuOpen)}
            className="p-2 rounded-full text-gray-400 hover:text-gray-700 hover:bg-gray-100 transition-colors"
            aria-label="Opsi postingan"
          >
            <MoreVertical className="w-5 h-5" />
          </button>

          {menuOpen && (
            <div className="absolute right-0 mt-1 w-44 bg-white rounded-xl shadow-lg border border-gray-100 py-1.5 z-20 animate-in fade-in zoom-in-95 duration-100">
              {isOwner ? (
                <>
                  <Link
                    href={`/edit/${product.id}`}
                    onClick={() => setMenuOpen(false)}
                    className="flex items-center gap-2.5 px-4 py-2.5 text-xs font-semibold text-gray-700 hover:bg-gray-50 transition-colors"
                  >
                    <Edit3 className="w-4 h-4 text-gray-500" />
                    Edit Produk
                  </Link>
                  <button
                    onClick={() => {
                      setMenuOpen(false);
                      openDeleteModal(product);
                    }}
                    className="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs font-semibold text-red-600 hover:bg-red-50 transition-colors text-left"
                  >
                    <Trash2 className="w-4 h-4 text-red-600" />
                    Hapus Produk
                  </button>
                </>
              ) : (
                <button
                  onClick={handleWhatsAppOrder}
                  className="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs font-semibold text-[#0d6e42] hover:bg-[#e6f4ea] transition-colors text-left"
                >
                  Pesan via WhatsApp
                </button>
              )}
            </div>
          )}
        </div>
      </header>

      {/* Main Image */}
      <div className="relative aspect-square w-full bg-gray-100 overflow-hidden">
        <img
          src={product.imageUrl}
          alt={product.title}
          className="w-full h-full object-cover transition-transform duration-300 hover:scale-[1.02]"
        />
      </div>

      {/* Action Bar */}
      <div className="p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            {/* Like Button */}
            <button
              onClick={() => toggleLike(product.id)}
              className="group flex items-center gap-1.5 text-gray-700 hover:text-red-500 transition-colors"
            >
              <Heart
                className={`w-6 h-6 transition-transform group-active:scale-125 ${
                  isLiked ? 'fill-red-500 text-red-500' : 'text-gray-700'
                }`}
              />
              <span className="text-xs font-semibold text-gray-700">
                {product.likesCount}
              </span>
            </button>

            {/* Comment Button */}
            <button
              onClick={() => openCommentModal(product.id)}
              className="flex items-center gap-1.5 text-gray-700 hover:text-[#0d6e42] transition-colors"
            >
              <MessageCircle className="w-6 h-6 text-gray-700" />
              <span className="text-xs font-semibold text-gray-700">
                {product.commentsCount}
              </span>
            </button>

            {/* WhatsApp Instant Order Button */}
            <button
              onClick={handleWhatsAppOrder}
              className="p-1 rounded-full text-emerald-600 hover:bg-emerald-50 transition-colors"
              title="Pesan via WhatsApp"
            >
              <svg className="w-6 h-6" viewBox="0 0 24 24" fill="none">
                <path
                  d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.297-.347.446-.521.151-.172.2-.296.3-.495.099-.198.05-.372-.025-.521-.075-.148-.669-1.611-.916-2.206-.242-.579-.487-.501-.669-.51l-.57-.01c-.198 0-.52.074-.792.372s-1.04 1.016-1.04 2.479 1.065 2.876 1.213 3.074c.149.198 2.095 3.2 5.076 4.487.709.306 1.263.489 1.694.626.712.226 1.36.194 1.872.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z"
                  fill="#10B981"
                />
                <path
                  fillRule="evenodd"
                  clipRule="evenodd"
                  d="M12 2C6.477 2 2 6.477 2 12c0 1.95.558 3.771 1.524 5.313L2.05 21.95l4.743-1.442A9.957 9.957 0 0012 22c5.523 0 10-4.477 10-10S17.523 2 12 2zm-8 10c0-4.411 3.589-8 8-8s8 3.589 8 8-3.589 8-8 8a7.95 7.95 0 01-4.223-1.205l-.303-.18-2.82.857.87-2.753-.195-.312A7.953 7.953 0 014 12z"
                  fill="#10B981"
                />
              </svg>
            </button>
          </div>
        </div>

        {/* Caption */}
        <div className="text-sm text-gray-800 leading-relaxed">
          <span className="font-bold text-gray-900 mr-2">{product.username}</span>
          {product.caption}
        </div>

        {/* Hashtags */}
        {product.hashtags && product.hashtags.length > 0 && (
          <div className="flex flex-wrap gap-1.5 pt-1">
            {product.hashtags.map((tag, idx) => (
              <span key={idx} className="text-xs font-semibold text-[#0d6e42] hover:underline cursor-pointer">
                {tag}
              </span>
            ))}
          </div>
        )}

        {/* Timestamp */}
        <p className="text-[11px] text-gray-400 font-medium pt-1">
          {product.createdAt}
        </p>
      </div>
    </article>
  );
}
