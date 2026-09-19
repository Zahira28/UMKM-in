'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Compass, PlusSquare, User as UserIcon } from 'lucide-react';
import { useApp } from '../../context/AppContext';

export default function Sidebar() {
  const pathname = usePathname();
  const { currentUser } = useApp();

  const navItems = [
    { name: 'Beranda', href: '/', icon: Home },
    { name: 'Discover', href: '/discover', icon: Compass },
    { name: 'Create', href: '/create', icon: PlusSquare },
    { name: 'Profil', href: '/profile', icon: UserIcon },
  ];

  return (
    <aside className="w-64 bg-white border-r border-[#e5e8e3] min-h-screen flex flex-col justify-between p-6 sticky top-0 h-screen z-30">
      <div>
        {/* Brand Logo */}
        <div className="mb-8 px-2">
          <Link href="/" className="text-2xl font-black text-[#0d6e42] tracking-tight">
            UMKM-in
          </Link>
        </div>

        {/* Navigation Items */}
        <nav className="space-y-1.5">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive =
              item.href === '/'
                ? pathname === '/'
                : pathname.startsWith(item.href);

            return (
              <Link
                key={item.name}
                href={item.href}
                className={`flex items-center gap-3.5 px-4 py-3 rounded-xl transition-all duration-150 ${
                  isActive
                    ? 'bg-[#e6f4ea] text-[#0d6e42] font-semibold shadow-xs'
                    : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900 font-medium'
                }`}
              >
                <Icon className={`w-5 h-5 ${isActive ? 'text-[#0d6e42]' : 'text-gray-500'}`} />
                <span className="text-sm">{item.name}</span>
              </Link>
            );
          })}
        </nav>
      </div>

      {/* Logged in User Profile Footer */}
      <div className="pt-4 border-t border-gray-100">
        <Link
          href="/profile"
          className="flex items-center gap-3 px-2 py-2 rounded-xl hover:bg-gray-50 transition-colors group"
        >
          <div className="w-10 h-10 rounded-full overflow-hidden bg-gray-200 border border-gray-200 flex-shrink-0">
            {currentUser.avatarUrl ? (
              <img
                src={currentUser.avatarUrl}
                alt={currentUser.fullName}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center bg-[#e6f4ea] text-[#0d6e42]">
                <UserIcon className="w-5 h-5" />
              </div>
            )}
          </div>
          <div className="overflow-hidden">
            <h4 className="text-xs font-bold text-gray-900 group-hover:text-[#0d6e42] truncate">
              {currentUser.username}
            </h4>
            <p className="text-[11px] text-gray-400 truncate">Akun saya</p>
          </div>
        </Link>
      </div>
    </aside>
  );
}
