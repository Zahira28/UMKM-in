export interface User {
  id: string;
  username: string;
  fullName: string;
  avatarUrl?: string;
  bio?: string;
  location?: string;
  whatsappNumber: string;
  productsCount: number;
  followersCount: number;
}

export interface Comment {
  id: string;
  productId: string;
  username: string;
  userAvatar?: string;
  content: string;
  createdAt: string;
  replies?: Comment[];
}

export interface Product {
  id: string;
  userId: string;
  username: string;
  userAvatar?: string;
  userLocation: string;
  title: string;
  category: 'Makanan' | 'Minuman' | 'Pakaian' | 'Kriya' | 'Jasa';
  price: number;
  imageUrl: string;
  rawImageUrl?: string;
  caption: string;
  hashtags: string[];
  likesCount: number;
  commentsCount: number;
  createdAt: string;
}

export type VibePreset = 'Minimalis' | 'Rustic Wood' | 'Elegance Warm' | 'Bright Commercial';
