'use client';

import React from 'react';
import { usePathname } from 'next/navigation';
import Sidebar from './Sidebar';
import CommentsModal from '../feed/CommentsModal';
import DeleteModal from '../modals/DeleteModal';
import WarningModal from '../modals/WarningModal';

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const isAuthPage = ['/login', '/signup', '/onboarding'].includes(pathname);

  if (isAuthPage) {
    return <main className="min-h-screen bg-[#f6f8f5]">{children}</main>;
  }

  return (
    <div className="flex min-h-screen bg-[#f6f8f5] text-gray-900 font-sans">
      <Sidebar />
      <main className="flex-1 min-w-0 p-6 md:p-10 max-w-7xl mx-auto">
        {children}
      </main>

      {/* Global Modals */}
      <CommentsModal />
      <DeleteModal />
      <WarningModal />
    </div>
  );
}
