'use client';

import React from 'react';
import { AlertCircle } from 'lucide-react';
import { useApp } from '../../context/AppContext';

export default function DeleteModal() {
  const { deleteModalProduct, closeDeleteModal, confirmDeleteProduct } = useApp();

  if (!deleteModalProduct) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-in fade-in duration-150">
      <div className="absolute inset-0" onClick={closeDeleteModal} />

      <div className="relative w-full max-w-md bg-white rounded-3xl p-8 shadow-2xl z-10 text-center space-y-6 animate-in zoom-in-95 duration-200">
        {/* Warning Icon */}
        <div className="flex justify-start items-center gap-4">
          <div className="w-16 h-16 rounded-full bg-red-50 text-red-500 flex items-center justify-center flex-shrink-0">
            <AlertCircle className="w-9 h-9" />
          </div>
          <h3 className="text-2xl font-extrabold text-gray-900 tracking-tight text-left">
            Hapus Produk?
          </h3>
        </div>

        {/* Text */}
        <p className="text-sm text-gray-600 text-left leading-relaxed">
          Produk ini akan dihapus secara permanen dan tidak dapat dikembalikan.
        </p>

        {/* Buttons */}
        <div className="flex items-center gap-4 pt-2">
          <button
            onClick={closeDeleteModal}
            className="flex-1 py-3.5 px-5 border border-gray-200 rounded-2xl font-bold text-gray-800 bg-white hover:bg-gray-50 transition-colors"
          >
            Batal
          </button>
          <button
            onClick={confirmDeleteProduct}
            className="flex-1 py-3.5 px-5 bg-red-600 hover:bg-red-700 text-white rounded-2xl font-bold transition-colors shadow-xs"
          >
            Hapus
          </button>
        </div>
      </div>
    </div>
  );
}
