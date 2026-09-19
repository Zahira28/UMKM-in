'use client';

import React, { useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { Upload, Sparkles, Wand2, ArrowLeft } from 'lucide-react';
import { useApp } from '../../src/context/AppContext';
import AICreativeStudioModal from '../../src/components/ai/AICreativeStudioModal';

export default function CreatePage() {
  const router = useRouter();
  const { addProduct, currentUser } = useApp();

  const [rawImage, setRawImage] = useState<string | null>(null);
  const [studioImage, setStudioImage] = useState<string | null>(null);
  const [title, setTitle] = useState('');
  const [category, setCategory] = useState<'Makanan' | 'Minuman' | 'Pakaian' | 'Kriya' | 'Jasa'>('Makanan');
  const [price, setPrice] = useState<string>('35000');
  const [description, setDescription] = useState('');
  const [hashtags, setHashtags] = useState<string[]>(['#umkmlokal', '#produkestetik']);

  const [aiStudioOpen, setAiStudioOpen] = useState(false);
  const [isImprovingCaption, setIsImprovingCaption] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        const result = event.target?.result as string;
        setRawImage(result);
        setStudioImage(result);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleImproveWithAI = () => {
    setIsImprovingCaption(true);
    setTimeout(() => {
      setIsImprovingCaption(false);
      const generatedCaptions: Record<string, string> = {
        Makanan: `${title || 'Kue Gula Aren'} buatan rumahan dengan bahan pilihan segar setiap hari. Dipanggang pas, manisnya gula aren alami tidak bikin enek! Pas untuk camilan santai keluarga 🍪`,
        Minuman: `${title || 'Kopi Susu Aren'} rahasia rasa otentik dari biji kopi lokal Indonesia. Segar, creamy, dan gula arennya legit banget ☕`,
        Pakaian: `${title || 'Batik Tulis Handmade'} karya perajin lokal pilihan. Kain adem, bahan jatuh, dan motif unik yang memancarkan keanggunan nusantara ✨`,
        Kriya: `${title || 'Produk Kriya Artisan'} buatan tangan terlatih. Detail rapi, estetik, dan tahan lama untuk mempercantik rumahmu 🌿`,
        Jasa: `Layanan ${title || 'Kreatif Lokal'} profesional dengan pengerjaan teliti & cepat. Siap membantu kebutuhan bisnis UMKM Anda!`,
      };
      setDescription(generatedCaptions[category] || `Produk unggulan ${title} berkualitas tinggi dari UMKM lokal.`);
      setHashtags([`#${category.toLowerCase()}lokal`, '#umkmjuara', '#karyanusantara']);
    }, 1200);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!title || (!rawImage && !studioImage)) {
      alert('Harap isi nama produk dan unggah foto produk.');
      return;
    }

    addProduct({
      userId: currentUser.id,
      username: currentUser.username,
      userAvatar: currentUser.avatarUrl,
      userLocation: currentUser.location || 'Bandung, Jawa Barat',
      title,
      category,
      price: parseInt(price, 10) || 0,
      imageUrl: studioImage || rawImage || 'https://images.unsplash.com/photo-1558961363-fa8fdf82db35?auto=format&fit=crop&w=1000&q=80',
      rawImageUrl: rawImage || undefined,
      caption: description || `Produk ${title} istimewa buatan rumahan.`,
      hashtags,
    });

    router.push('/');
  };

  return (
    <div className="space-y-6 max-w-2xl mx-auto">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight">
          Create
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          Bagikan produk lokalmu kepada komunitas.
        </p>
      </div>

      {/* Main Form Container */}
      <form onSubmit={handleSubmit} className="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm space-y-6">
        {/* Gambar Produk Section */}
        <div className="space-y-3">
          <label className="text-sm font-bold text-gray-900 block">
            Gambar Produk
          </label>

          <input
            type="file"
            ref={fileInputRef}
            onChange={handleImageUpload}
            accept="image/*"
            className="hidden"
          />

          <div
            onClick={() => fileInputRef.current?.click()}
            className="bg-[#f6f8f5] border-2 border-dashed border-gray-200 hover:border-[#0d6e42] rounded-2xl p-8 text-center cursor-pointer transition-colors group relative overflow-hidden min-h-[220px] flex flex-col items-center justify-center"
          >
            {studioImage ? (
              <div className="relative w-full h-56 rounded-xl overflow-hidden group">
                <img
                  src={studioImage}
                  alt="Preview"
                  className="w-full h-full object-cover"
                />
                <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center text-white text-xs font-bold">
                  Klik untuk mengganti foto
                </div>
              </div>
            ) : (
              <div className="space-y-2">
                <div className="w-12 h-12 rounded-full bg-[#e6f4ea] text-[#0d6e42] flex items-center justify-center mx-auto group-hover:scale-110 transition-transform">
                  <Upload className="w-6 h-6" />
                </div>
                <h4 className="text-sm font-bold text-gray-900">Upload Foto</h4>
                <p className="text-xs text-gray-400">JPG atau PNG · Maks. 10 MB</p>
              </div>
            )}
          </div>

          {/* AI Background Studio Button */}
          <button
            type="button"
            onClick={() => {
              if (!rawImage) {
                // If no image, set a sample default for demo
                const sample = 'https://images.unsplash.com/photo-1558961363-fa8fdf82db35?auto=format&fit=crop&w=800&q=80';
                setRawImage(sample);
                setStudioImage(sample);
              }
              setAiStudioOpen(true);
            }}
            className="inline-flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 hover:border-[#0d6e42] text-gray-800 rounded-xl text-xs font-bold transition-all hover:bg-gray-50 shadow-2xs"
          >
            <Sparkles className="w-4 h-4 text-[#0d6e42]" />
            Perbarui dengan AI
          </button>
        </div>

        {/* Informasi Produk Section */}
        <div className="space-y-5 pt-2 border-t border-gray-100">
          <h3 className="text-sm font-bold text-gray-900">Informasi Produk</h3>

          {/* Nama Produk */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-700 block">
              Nama produk
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Contoh: Cookies gula aren"
              className="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
              required
            />
          </div>

          {/* Kategori */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-700 block">
              Kategori
            </label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value as any)}
              className="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
            >
              <option value="Makanan">Makanan</option>
              <option value="Minuman">Minuman</option>
              <option value="Pakaian">Pakaian</option>
              <option value="Kriya">Kriya</option>
              <option value="Jasa">Jasa</option>
            </select>
          </div>

          {/* Harga */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-700 block">
              Harga (Rp)
            </label>
            <input
              type="number"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              placeholder="35000"
              className="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
              required
            />
          </div>

          {/* Deskripsi */}
          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-700 block">
              Deskripsi
            </label>
            <textarea
              rows={4}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Ceritakan produk, bahan, dan keunikannya..."
              className="w-full bg-white border border-gray-200 rounded-xl p-4 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all leading-relaxed"
            />

            {/* AI Copywriting Button */}
            <button
              type="button"
              onClick={handleImproveWithAI}
              disabled={isImprovingCaption}
              className="inline-flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 hover:border-[#0d6e42] text-gray-800 rounded-xl text-xs font-bold transition-all hover:bg-gray-50 shadow-2xs mt-1"
            >
              <Wand2 className="w-4 h-4 text-[#0d6e42]" />
              {isImprovingCaption ? 'Menulis Copywriting AI...' : 'Improve dengan AI'}
            </button>
          </div>
        </div>

        {/* Submit Button */}
        <div className="pt-4">
          <button
            type="submit"
            className="w-full py-4 bg-[#0d6e42] hover:bg-[#095432] text-white rounded-2xl font-bold transition-colors shadow-xs text-center text-sm"
          >
            Buat Produk
          </button>
        </div>
      </form>

      {/* AI Studio Modal */}
      <AICreativeStudioModal
        isOpen={aiStudioOpen}
        onClose={() => setAiStudioOpen(false)}
        rawImage={rawImage || 'https://images.unsplash.com/photo-1558961363-fa8fdf82db35?auto=format&fit=crop&w=800&q=80'}
        onApplyImage={(img) => setStudioImage(img)}
      />
    </div>
  );
}
