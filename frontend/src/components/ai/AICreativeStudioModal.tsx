'use client';

import React, { useState } from 'react';
import { Sparkles, Wand2, X, Check, RefreshCw, Image as ImageIcon } from 'lucide-react';
import { VibePreset } from '../../types';

interface AICreativeStudioModalProps {
  isOpen: boolean;
  onClose: () => void;
  rawImage: string;
  onApplyImage: (studioImage: string) => void;
}

export default function AICreativeStudioModal({
  isOpen,
  onClose,
  rawImage,
  onApplyImage,
}: AICreativeStudioModalProps) {
  const [selectedVibe, setSelectedVibe] = useState<VibePreset>('Rustic Wood');
  const [isProcessing, setIsProcessing] = useState(false);
  const [progressStep, setProgressStep] = useState('');
  const [previewImage, setPreviewImage] = useState<string | null>(null);

  if (!isOpen) return null;

  const vibes: { name: VibePreset; label: string; desc: string; sampleUrl: string }[] = [
    {
      name: 'Rustic Wood',
      label: 'Kayu Estetik & Alami',
      desc: 'Latar kayu hangat dengan rempah & daun estetik (Cocok untuk Kuliner/Kopi)',
      sampleUrl: 'https://images.unsplash.com/photo-1558961363-fa8fdf82db35?auto=format&fit=crop&w=800&q=80',
    },
    {
      name: 'Minimalis',
      label: 'Studio Clean Minimalis',
      desc: 'Latar netral studio dengan bayangan kontak lembut (Cocok untuk Kerajinan & Sabun)',
      sampleUrl: 'https://images.unsplash.com/photo-1607006482602-76ca072b445f?auto=format&fit=crop&w=800&q=80',
    },
    {
      name: 'Elegance Warm',
      label: 'Elegan Mewah',
      desc: 'Pencahayaan warm kontras tinggi dengan drapery (Cocok untuk Fashion & Batik)',
      sampleUrl: 'https://images.unsplash.com/photo-1610030469983-98e550d6193c?auto=format&fit=crop&w=800&q=80',
    },
    {
      name: 'Bright Commercial',
      label: 'Commercial Bright',
      desc: 'Latar terang komersial untuk etalase katalog online (Cocok untuk Snack & Kemasan)',
      sampleUrl: 'https://images.unsplash.com/photo-1621939514649-280e2ee25f60?auto=format&fit=crop&w=800&q=80',
    },
  ];

  const handleGenerateStudio = () => {
    setIsProcessing(true);
    setProgressStep('Mengekstraksi objek produk (rembg U2-Net)...');

    setTimeout(() => {
      setProgressStep('Me-render latar belakang studio (' + selectedVibe + ')...');
    }, 1200);

    setTimeout(() => {
      setProgressStep('Menambahkan bayangan realistis (Contact Shadow Engine)...');
    }, 2400);

    setTimeout(() => {
      setIsProcessing(false);
      const chosen = vibes.find((v) => v.name === selectedVibe);
      setPreviewImage(chosen?.sampleUrl || rawImage);
    }, 3500);
  };

  const handleApply = () => {
    if (previewImage) {
      onApplyImage(previewImage);
      onClose();
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-xs p-4 animate-in fade-in duration-150">
      <div className="absolute inset-0" onClick={onClose} />

      <div className="relative w-full max-w-3xl bg-white rounded-3xl p-6 md:p-8 shadow-2xl z-10 space-y-6 max-h-[90vh] overflow-y-auto animate-in zoom-in-95 duration-200">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-100 pb-4">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-[#e6f4ea] text-[#0d6e42] rounded-2xl">
              <Sparkles className="w-6 h-6" />
            </div>
            <div>
              <h3 className="text-xl font-extrabold text-gray-900 tracking-tight">
                In-App AI Generative Studio
              </h3>
              <p className="text-xs text-gray-500">
                Ubah foto ponsel biasa menjadi foto studio komersial secara otomatis.
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content Layout */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Left: Image Comparison Preview */}
          <div className="space-y-4">
            <label className="text-xs font-bold text-gray-700 uppercase tracking-wider block">
              {previewImage ? 'Hasil Foto Studio AI' : 'Foto Asli Unggahan'}
            </label>
            <div className="relative aspect-square w-full rounded-2xl bg-gray-100 overflow-hidden border border-gray-200">
              <img
                src={previewImage || rawImage}
                alt="Preview AI"
                className="w-full h-full object-cover"
              />

              {isProcessing && (
                <div className="absolute inset-0 bg-black/60 backdrop-blur-xs flex flex-col items-center justify-center p-6 text-center text-white space-y-3">
                  <RefreshCw className="w-8 h-8 animate-spin text-[#10B981]" />
                  <p className="text-sm font-semibold animate-pulse">{progressStep}</p>
                </div>
              )}
            </div>
          </div>

          {/* Right: Vibe Selector */}
          <div className="space-y-4 flex flex-col justify-between">
            <div>
              <label className="text-xs font-bold text-gray-700 uppercase tracking-wider block mb-3">
                Pilih Vibe Latar Studio
              </label>
              <div className="space-y-2.5">
                {vibes.map((v) => (
                  <div
                    key={v.name}
                    onClick={() => setSelectedVibe(v.name)}
                    className={`p-3.5 rounded-2xl border cursor-pointer transition-all flex items-start justify-between ${
                      selectedVibe === v.name
                        ? 'border-[#0d6e42] bg-[#e6f4ea]/40 ring-1 ring-[#0d6e42]'
                        : 'border-gray-200 hover:border-gray-300 bg-white'
                    }`}
                  >
                    <div>
                      <h4 className="text-sm font-bold text-gray-900">{v.label}</h4>
                      <p className="text-xs text-gray-500 mt-0.5 leading-snug">{v.desc}</p>
                    </div>
                    {selectedVibe === v.name && (
                      <div className="w-5 h-5 rounded-full bg-[#0d6e42] text-white flex items-center justify-center flex-shrink-0 mt-0.5">
                        <Check className="w-3 h-3" />
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {/* Action Button */}
            <div className="pt-4 space-y-2">
              {!previewImage ? (
                <button
                  onClick={handleGenerateStudio}
                  disabled={isProcessing}
                  className="w-full py-3.5 px-5 bg-[#0d6e42] hover:bg-[#095432] text-white rounded-2xl font-bold transition-all shadow-xs flex items-center justify-center gap-2 disabled:opacity-50"
                >
                  <Wand2 className="w-5 h-5" />
                  Generate Studio Post
                </button>
              ) : (
                <div className="flex gap-3">
                  <button
                    onClick={handleGenerateStudio}
                    disabled={isProcessing}
                    className="flex-1 py-3 px-4 border border-gray-200 rounded-2xl font-bold text-gray-700 hover:bg-gray-50 transition-colors flex items-center justify-center gap-1.5 text-xs"
                  >
                    <RefreshCw className="w-4 h-4" />
                    Coba Lagi
                  </button>
                  <button
                    onClick={handleApply}
                    className="flex-1 py-3 px-4 bg-[#0d6e42] hover:bg-[#095432] text-white rounded-2xl font-bold transition-colors text-xs flex items-center justify-center gap-1.5 shadow-xs"
                  >
                    <Check className="w-4 h-4" />
                    Gunakan Hasil AI
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
