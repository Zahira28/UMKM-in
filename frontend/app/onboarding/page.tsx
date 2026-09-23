'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { User } from 'lucide-react';

export default function OnboardingPage() {
  const router = useRouter();
  const [name, setName] = useState('Ayu Lestari');
  const [username, setUsername] = useState('ayu_lestari');
  const [location, setLocation] = useState('Bandung, Jawa Barat');
  const [preference, setPreference] = useState('Makanan & Minuman');
  const [whatsapp, setWhatsapp] = useState('+62 812 3456 7890');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    router.push('/');
  };

  return (
    <div className="min-h-screen flex flex-col md:flex-row bg-[#f6f8f5]">
      {/* Left Dark Green Banner Panel */}
      <div className="w-full md:w-1/2 bg-[#0d6e42] text-white p-8 md:p-16 flex flex-col justify-between min-h-[300px] md:min-h-screen">
        <div>
          <span className="text-3xl font-black tracking-tight">UMKM-in</span>
        </div>

        <div className="max-w-md my-auto space-y-4 py-8">
          <h1 className="text-3xl md:text-5xl font-extrabold leading-tight">
            Temukan karya lokal, dukung cerita di baliknya.
          </h1>
          <p className="text-sm md:text-base text-emerald-100/90 leading-relaxed">
            Ruang untuk produk UMKM tumbuh, ditemukan, dan terhubung dengan pelanggan baru.
          </p>
        </div>

        <div className="text-xs text-emerald-200/60 hidden md:block">
          &copy; {new Date().getFullYear()} UMKM-in. Social Commerce & AI Creative Studio.
        </div>
      </div>

      {/* Right Onboarding Card Panel */}
      <div className="w-full md:w-1/2 p-6 md:p-12 flex items-center justify-center">
        <div className="w-full max-w-md bg-white rounded-3xl p-8 md:p-10 shadow-sm border border-gray-100 space-y-5">
          <div>
            <h2 className="text-2xl font-extrabold text-gray-900 tracking-tight">
              Lengkapi profilmu
            </h2>
            <p className="text-xs text-gray-500 mt-1">
              Bantu komunitas mengenalmu lebih dekat.
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Profile Avatar Upload Dropzone */}
            <div className="flex items-center gap-4 py-1">
              <div className="w-16 h-16 rounded-full bg-gray-100 border border-gray-200 flex items-center justify-center text-gray-400 flex-shrink-0">
                <User className="w-8 h-8 stroke-1" />
              </div>
              <div>
                <button
                  type="button"
                  className="text-xs font-bold text-[#138a53] hover:underline"
                >
                  Upload foto profil
                </button>
                <p className="text-[11px] text-gray-400 mt-0.5">
                  JPG atau PNG, maksimal 2 MB
                </p>
              </div>
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Nama
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Nama lengkap"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Username
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="contoh: ayu_lestari"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Lokasi
              </label>
              <input
                type="text"
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                placeholder="Kota, provinsi"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Preferensi
              </label>
              <input
                type="text"
                value={preference}
                onChange={(e) => setPreference(e.target.value)}
                placeholder="Pilih kategori favorit"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                No. WhatsApp
              </label>
              <input
                type="text"
                value={whatsapp}
                onChange={(e) => setWhatsapp(e.target.value)}
                placeholder="+62 812 3456 7890"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <button
              type="submit"
              className="w-full py-3.5 bg-[#138a53] hover:bg-[#0d6e42] text-white font-bold rounded-2xl transition-colors shadow-xs text-sm mt-3"
            >
              Simpan dan lanjutkan
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
