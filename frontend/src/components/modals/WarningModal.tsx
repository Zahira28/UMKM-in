'use client';

import React from 'react';
import { AlertTriangle } from 'lucide-react';
import { useApp } from '../../context/AppContext';

export default function WarningModal() {
  const { warningModalConfig, closeWarningModal } = useApp();

  if (!warningModalConfig?.isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-in fade-in duration-150">
      <div className="absolute inset-0" onClick={closeWarningModal} />

      <div className="relative w-full max-w-md bg-white rounded-3xl p-8 shadow-2xl z-10 space-y-6 animate-in zoom-in-95 duration-200">
        {/* Warning Icon & Title */}
        <div className="space-y-4">
          <div className="w-12 h-12 rounded-full bg-red-50 text-red-500 flex items-center justify-center">
            <AlertTriangle className="w-6 h-6" />
          </div>
          <h3 className="text-2xl font-extrabold text-gray-900 tracking-tight">
            Peringatan
          </h3>
        </div>

        {/* Text */}
        <p className="text-sm text-gray-600 leading-relaxed">
          Keluar dari halaman Create akan menghapus seluruh progres yang ada !
        </p>

        {/* Stacked Buttons */}
        <div className="space-y-3 pt-2">
          <button
            onClick={() => {
              if (warningModalConfig.onConfirm) warningModalConfig.onConfirm();
              closeWarningModal();
            }}
            className="w-full py-3.5 px-5 bg-red-600 hover:bg-red-700 text-white rounded-2xl font-bold transition-colors shadow-xs"
          >
            Hapus
          </button>
          <button
            onClick={closeWarningModal}
            className="w-full py-3.5 px-5 bg-[#0d6e42] hover:bg-[#095432] text-white rounded-2xl font-bold transition-colors shadow-xs"
          >
            Lanjutkan Mengedit
          </button>
        </div>
      </div>
    </div>
  );
}
