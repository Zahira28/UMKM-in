'use client';

import React from 'react';
import { useApp } from '../../src/context/AppContext';
import Link from 'next/link';
import { Heart, MessageCircle } from 'lucide-react';

export default function ProfilePage() {
  const { currentUser, products, openCommentModal } = useApp();

  const userProducts = products.filter(
    (p) => p.username === currentUser.username || p.userId === currentUser.id
  );

  return (
    <div className="space-y-8 max-w-5xl mx-auto">
      {/* Profile Header */}
      <div className="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
        <div className="flex flex-col md:flex-row items-start md:items-center gap-6">
          {/* Avatar */}
          <div className="w-24 h-24 md:w-28 md:h-28 rounded-full bg-gray-100 overflow-hidden border-2 border-gray-100 flex-shrink-0">
            <img
              src={currentUser.avatarUrl}
              alt={currentUser.fullName}
              className="w-full h-full object-cover"
            />
          </div>

          {/* Info */}
          <div className="space-y-2">
            <div>
              <h1 className="text-2xl md:text-3xl font-extrabold text-gray-900 tracking-tight">
                {currentUser.username}
              </h1>
              <p className="text-sm font-bold text-gray-800 mt-0.5">
                {currentUser.fullName}
              </p>
            </div>

            <p className="text-sm text-gray-600 max-w-lg leading-relaxed">
              {currentUser.bio}
            </p>

            <div className="flex items-center gap-4 text-xs font-semibold text-gray-600 pt-1">
              <span>{userProducts.length} Produk</span>
              <span>·</span>
              <span>{currentUser.followersCount.toLocaleString('id-ID')} Pengikut</span>
            </div>
          </div>
        </div>

        {/* Edit Profile Button */}
        <button className="px-5 py-2.5 bg-white border border-gray-200 rounded-xl text-xs font-bold text-gray-800 hover:bg-gray-50 transition-colors shadow-2xs self-stretch md:self-start text-center">
          Edit profil
        </button>
      </div>

      {/* Section Title */}
      <div className="space-y-4">
        <h2 className="text-xl font-extrabold text-gray-900 tracking-tight">
          Produk
        </h2>

        {/* Grid of Products (3 columns) */}
        {userProducts.length === 0 ? (
          <div className="bg-white rounded-3xl p-12 text-center text-gray-400 border border-gray-100">
            Belum ada produk yang diunggah. Klik menu <Link href="/create" className="text-[#0d6e42] underline font-bold">Create</Link> untuk memposting produk pertamamu!
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-5">
            {userProducts.map((product) => (
              <div
                key={product.id}
                className="group relative bg-white rounded-2xl overflow-hidden border border-gray-100 shadow-xs aspect-square cursor-pointer transition-all duration-300 hover:shadow-lg"
              >
                <img
                  src={product.imageUrl}
                  alt={product.title}
                  className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                />

                <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity duration-200 p-4 flex flex-col justify-between text-white">
                  <span className="bg-white/20 backdrop-blur-xs text-white text-[11px] font-semibold px-2.5 py-1 rounded-full w-fit">
                    {product.category}
                  </span>

                  <div>
                    <h4 className="font-bold text-sm line-clamp-1">{product.title}</h4>
                    <p className="text-xs text-white/80">Rp {product.price.toLocaleString('id-ID')}</p>

                    <div className="flex items-center gap-4 text-xs font-semibold pt-2">
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
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
