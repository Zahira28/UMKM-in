'use client';

import React, { use } from 'react';
import { useApp } from '../../../src/context/AppContext';
import { otherUsersMock } from '../../../src/lib/mock-data';
import { Heart, MessageCircle } from 'lucide-react';

export default function UserProfilePage({ params }: { params: Promise<{ username: string }> }) {
  const { username } = use(params);
  const normalizedUsername = decodeURIComponent(username).toLowerCase();
  const { products, followingMap, toggleFollow, openCommentModal, openDetailModal, toggleLike, likesMap } = useApp();

  const isFollowed = !!followingMap[normalizedUsername] || !!followingMap[username];
  const userObj = otherUsersMock[normalizedUsername] || otherUsersMock[username] || {
    id: 'u_' + username,
    username,
    fullName: username.replace('_', ' '),
    avatarUrl: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=300&q=80',
    bio: 'Merangkai produk lokal dengan penuh cerita ✨ Bandung · Terima pesanan melalui WhatsApp',
    location: 'Bandung, Jawa Barat',
    whatsappNumber: '+6281234567890',
    productsCount: 9,
    followersCount: 1200,
  };

  const userProducts = products.filter(
    (p) => p.username.toLowerCase() === normalizedUsername || p.username === username
  );

  const handleWhatsAppChat = () => {
    const message = encodeURIComponent(`Halo *${userObj.fullName}*, saya ingin bertanya mengenai produk UMKM-in Anda!`);
    window.open(`https://wa.me/${userObj.whatsappNumber.replace(/[^0-9]/g, '')}?text=${message}`, '_blank');
  };

  return (
    <div className="space-y-8 max-w-5xl mx-auto">
      {/* Profile Header Card */}
      <div className="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm flex flex-col md:flex-row items-start md:items-center justify-between gap-6">
        <div className="flex flex-col md:flex-row items-start md:items-center gap-6">
          <div className="w-24 h-24 md:w-28 md:h-28 rounded-full bg-gray-100 overflow-hidden border-2 border-gray-100 flex-shrink-0">
            <img
              src={userObj.avatarUrl}
              alt={userObj.fullName}
              className="w-full h-full object-cover"
            />
          </div>

          <div className="space-y-2">
            <div>
              <h1 className="text-2xl md:text-3xl font-extrabold text-gray-900 tracking-tight">
                {userObj.username}
              </h1>
              <p className="text-sm font-bold text-gray-800 mt-0.5">
                {userObj.fullName}
              </p>
            </div>

            <p className="text-sm text-gray-600 max-w-lg leading-relaxed">
              {userObj.bio}
            </p>

            <div className="flex items-center gap-4 text-xs font-semibold text-gray-600 pt-1">
              <span>{userProducts.length || userObj.productsCount} Produk</span>
              <span>·</span>
              <span>
                {(userObj.followersCount + (isFollowed ? 1 : 0)).toLocaleString('id-ID')} Pengikut
              </span>
            </div>
          </div>
        </div>

        {/* Buttons (Follow & WA) */}
        <div className="flex items-center gap-3 self-stretch md:self-start">
          <button
            onClick={() => toggleFollow(username)}
            className={`px-6 py-2.5 rounded-xl text-xs font-bold transition-all shadow-2xs ${
              isFollowed
                ? 'bg-white border border-gray-200 text-gray-800 hover:bg-gray-50'
                : 'bg-[#0d6e42] hover:bg-[#095432] text-white'
            }`}
          >
            {isFollowed ? 'Followed' : 'Follow'}
          </button>
          <button
            onClick={handleWhatsAppChat}
            className="px-4 py-2.5 bg-emerald-50 text-[#0d6e42] hover:bg-emerald-100 rounded-xl text-xs font-bold transition-colors"
          >
            Chat WA
          </button>
        </div>
      </div>

      {/* Products Section */}
      <div className="space-y-4">
        <h2 className="text-xl font-extrabold text-gray-900 tracking-tight">
          Produk
        </h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-5">
          {userProducts.map((product) => {
            const isLiked = !!likesMap[product.id];
            return (
              <div
                key={product.id}
                onClick={() => openDetailModal(product.id)}
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
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          toggleLike(product.id);
                        }}
                        className="flex items-center gap-1 hover:text-red-300 transition-colors"
                      >
                        <Heart
                          fill={isLiked ? '#ef4444' : 'white'}
                          className={`w-4 h-4 ${isLiked ? 'text-red-500' : 'text-white'}`}
                        />
                        <span>{product.likesCount}</span>
                      </button>
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
            );
          })}
        </div>
      </div>
    </div>
  );
}
