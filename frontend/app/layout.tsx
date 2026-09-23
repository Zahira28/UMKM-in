import type { Metadata } from 'next';
import './globals.css';
import { AppProvider } from '../src/context/AppContext';
import AppLayout from '../src/components/layout/AppLayout';

export const metadata: Metadata = {
  title: 'UMKM-in | Social Commerce & AI Creative Studio',
  description: 'Official Web App UMKM-in - Social Commerce Feed & Generative AI Studio khusus Produk UMKM Lokal',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="id" className="h-full antialiased">
      <body className="min-h-full flex flex-col bg-[#f6f8f5]">
        <AppProvider>
          <AppLayout>{children}</AppLayout>
        </AppProvider>
      </body>
    </html>
  );
}
