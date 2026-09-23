'use client';

import React, { useState, useEffect, useRef, use } from 'react';
import { useRouter } from 'next/navigation';
import { Upload, Sparkles, Wand2 } from 'lucide-react';
import { useApp } from '../../../src/context/AppContext';
import AICreativeStudioModal from '../../../src/components/ai/AICreativeStudioModal';

export default function EditProductPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const { products, updateProduct } = useApp();

  const targetProduct = products.find((p) => p.id === id);

  const [rawImage, setRawImage] = useState<string>('');
  const [studioImage, setStudioImage] = useState<string>('');
  const [title, setTitle] = useState('');
  const [category, setCategory] = useState<'Makanan' | 'Minuman' | 'Pakaian' | 'Kriya' | 'Jasa'>('Makanan');
  const [price, setPrice] = useState<string>('');
  const [description, setDescription] = useState('');

  const [aiStudioOpen, setAiStudioOpen] = useState(false);
  const [isImprovingCaption, setIsImprovingCaption] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (targetProduct) {
      setRawImage(targetProduct.rawImageUrl || targetProduct.imageUrl);
      setStudioImage(targetProduct.imageUrl);
      setTitle(targetProduct.title);
      setCategory(targetProduct.category);
      setPrice(targetProduct.price.toString());
      setDescription(targetProduct.caption);
    }
  }, [targetProduct]);

  if (!targetProduct) {
    return (
      <div className="text-center py-20">
        <h2 className="text-xl font-bold text-gray-800">Produk tidak ditemukan.</h2>
        <button
          onClick={() => router.push('/')}
          className="mt-4 px-4 py-2 bg-[#0d6e42] text-white rounded-xl text-sm font-semibold"
        >
          Kembali ke Beranda
        </button>
      </div>
    );
  }

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
      setDescription(
        `${title} - Diolah istimewa dengan bahan alami khas kuliner lokal. Tekstur lezat, aroma wangi, dan dikemas secara higienis untuk menjaga mutu terbaik!`
      );
    }, 1000);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    updateProduct(id, {
      title,
      category,
      price: parseInt(price, 10) || 0,
      imageUrl: studioImage || rawImage,
      caption: description,
    });
    router.push('/');
  };

  return (
    <div className="space-y-6 max-w-2xl mx-auto">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight">
          Edit
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          Perbarui informasi produk kamu.
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
            className="bg-[#f6f8f5] border-2 border-dashed border-gray-200 hover:border-[#0d6e42] rounded-2xl p-4 text-center cursor-pointer transition-colors group relative overflow-hidden flex flex-col items-center justify-center min-h-[200px]"
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
                <Upload className="w-8 h-8 text-[#0d6e42] mx-auto" />
                <h4 className="text-sm font-bold text-gray-900">Upload Foto</h4>
                <p className="text-xs text-gray-400">JPG atau PNG · Maks. 10 MB</p>
              </div>
            )}
          </div>

          <button
            type="button"
            onClick={() => setAiStudioOpen(true)}
            className="inline-flex items-center gap-2 px-4 py-2.5 bg-white border border-gray-200 hover:border-[#0d6e42] text-gray-800 rounded-xl text-xs font-bold transition-all hover:bg-gray-50 shadow-2xs"
          >
            <Sparkles className="w-4 h-4 text-[#0d6e42]" />
            Perbarui dengan AI
          </button>
        </div>

        {/* Informasi Produk Section */}
        <div className="space-y-5 pt-2 border-t border-gray-100">
          <h3 className="text-sm font-bold text-gray-900">Informasi Produk</h3>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-700 block">
              Nama produk
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
              required
            />
          </div>

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

          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-700 block">
              Harga (Rp)
            </label>
            <input
              type="number"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              className="w-full bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all"
              required
            />
          </div>

          <div className="space-y-1.5">
            <label className="text-xs font-semibold text-gray-700 block">
              Deskripsi
            </label>
            <textarea
              rows={4}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full bg-white border border-gray-200 rounded-xl p-4 text-sm text-gray-900 focus:outline-none focus:border-[#0d6e42] focus:ring-1 focus:ring-[#0d6e42] transition-all leading-relaxed"
            />

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

        <div className="pt-4">
          <button
            type="submit"
            className="w-full py-4 bg-[#0d6e42] hover:bg-[#095432] text-white rounded-2xl font-bold transition-colors shadow-xs text-center text-sm"
          >
            Perbarui Produk
          </button>
        </div>
      </form>

      <AICreativeStudioModal
        isOpen={aiStudioOpen}
        onClose={() => setAiStudioOpen(false)}
        rawImage={rawImage}
        onApplyImage={(img) => setStudioImage(img)}
      />
    </div>
  );
}
