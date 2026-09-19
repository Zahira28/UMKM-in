'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

export default function SignupPage() {
  const router = useRouter();
  const [fullName, setFullName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');

  const handleSignup = (e: React.FormEvent) => {
    e.preventDefault();
    router.push('/onboarding');
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

      {/* Right Signup Card Panel */}
      <div className="w-full md:w-1/2 p-6 md:p-12 flex items-center justify-center">
        <div className="w-full max-w-md bg-white rounded-3xl p-8 md:p-10 shadow-sm border border-gray-100 space-y-5">
          <div>
            <h2 className="text-2xl font-extrabold text-gray-900 tracking-tight">
              Buat akun baru
            </h2>
            <p className="text-xs text-gray-500 mt-1">
              Mulai bagikan dan temukan produk UMKM favoritmu.
            </p>
          </div>

          <form onSubmit={handleSignup} className="space-y-3.5">
            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Nama lengkap
              </label>
              <input
                type="text"
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                placeholder="Nama lengkapmu"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Email
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="nama@email.com"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Password
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Minimal 8 karakter"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <div className="space-y-1">
              <label className="text-xs font-bold text-gray-800 block">
                Konfirmasi password
              </label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Ulangi password"
                className="w-full bg-white border border-gray-200 rounded-2xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
                required
              />
            </div>

            <button
              type="submit"
              className="w-full py-3.5 bg-[#138a53] hover:bg-[#0d6e42] text-white font-bold rounded-2xl transition-colors shadow-xs text-sm mt-2"
            >
              Daftar
            </button>
          </form>

          {/* Google Sign In Button */}
          <button
            onClick={() => router.push('/onboarding')}
            className="w-full py-3 px-4 border border-gray-200 rounded-2xl text-xs font-bold text-gray-800 hover:bg-gray-50 transition-colors flex items-center justify-center gap-2.5"
          >
            <svg className="w-4 h-4" viewBox="0 0 24 24">
              <path
                fill="#4285F4"
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
              />
              <path
                fill="#34A853"
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
              />
              <path
                fill="#FBBC05"
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
              />
              <path
                fill="#EA4335"
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
              />
            </svg>
            Daftar dengan Google
          </button>

          <p className="text-center text-xs text-gray-600 pt-1">
            Sudah punya akun?{' '}
            <Link href="/login" className="text-[#138a53] font-bold hover:underline">
              Masuk
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
