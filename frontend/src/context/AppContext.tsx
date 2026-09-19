'use client';

import React, { createContext, useContext, useState } from 'react';
import { Product, User, Comment } from '../types';
import { currentUserMock, initialProductsMock, initialCommentsMock } from '../lib/mock-data';

interface AppContextType {
  products: Product[];
  currentUser: User;
  followingMap: Record<string, boolean>;
  likesMap: Record<string, boolean>;
  commentsMap: Record<string, Comment[]>;
  activeCommentProductId: string | null;
  deleteModalProduct: Product | null;
  warningModalConfig: { isOpen: boolean; onConfirm?: () => void } | null;
  
  // Actions
  toggleLike: (productId: string) => void;
  toggleFollow: (username: string) => void;
  addProduct: (product: Omit<Product, 'id' | 'likesCount' | 'commentsCount' | 'createdAt'>) => void;
  updateProduct: (id: string, updated: Partial<Product>) => void;
  deleteProduct: (id: string) => void;
  addComment: (productId: string, content: string) => void;
  
  // Modal handlers
  openCommentModal: (productId: string) => void;
  closeCommentModal: () => void;
  openDeleteModal: (product: Product) => void;
  closeDeleteModal: () => void;
  confirmDeleteProduct: () => void;
  openWarningModal: (onConfirm: () => void) => void;
  closeWarningModal: () => void;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

export function AppProvider({ children }: { children: React.ReactNode }) {
  const [products, setProducts] = useState<Product[]>(initialProductsMock);
  const [currentUser] = useState<User>(currentUserMock);
  const [followingMap, setFollowingMap] = useState<Record<string, boolean>>({
    jowi_lokal: false,
    dapur_nisa: true,
  });
  const [likesMap, setLikesMap] = useState<Record<string, boolean>>({});
  const [commentsMap, setCommentsMap] = useState<Record<string, Comment[]>>(initialCommentsMock);
  
  // Modals state
  const [activeCommentProductId, setActiveCommentProductId] = useState<string | null>(null);
  const [deleteModalProduct, setDeleteModalProduct] = useState<Product | null>(null);
  const [warningModalConfig, setWarningModalConfig] = useState<{ isOpen: boolean; onConfirm?: () => void } | null>(null);

  const toggleLike = (productId: string) => {
    setLikesMap((prev) => {
      const isLiked = !!prev[productId];
      const nextState = !isLiked;
      
      setProducts((currentProducts) =>
        currentProducts.map((p) =>
          p.id === productId
            ? { ...p, likesCount: isLiked ? p.likesCount - 1 : p.likesCount + 1 }
            : p
        )
      );
      
      return { ...prev, [productId]: nextState };
    });
  };

  const toggleFollow = (username: string) => {
    setFollowingMap((prev) => ({
      ...prev,
      [username]: !prev[username],
    }));
  };

  const addProduct = (newProdData: Omit<Product, 'id' | 'likesCount' | 'commentsCount' | 'createdAt'>) => {
    const newProduct: Product = {
      ...newProdData,
      id: 'p_' + Date.now(),
      likesCount: 0,
      commentsCount: 0,
      createdAt: 'Baru saja',
    };
    setProducts((prev) => [newProduct, ...prev]);
  };

  const updateProduct = (id: string, updated: Partial<Product>) => {
    setProducts((prev) =>
      prev.map((p) => (p.id === id ? { ...p, ...updated } : p))
    );
  };

  const deleteProduct = (id: string) => {
    setProducts((prev) => prev.filter((p) => p.id !== id));
  };

  const addComment = (productId: string, content: string) => {
    if (!content.trim()) return;
    
    const newComment: Comment = {
      id: 'c_' + Date.now(),
      productId,
      username: currentUser.username,
      userAvatar: currentUser.avatarUrl,
      content: content.trim(),
      createdAt: 'Baru saja',
    };

    setCommentsMap((prev) => ({
      ...prev,
      [productId]: [...(prev[productId] || []), newComment],
    }));

    setProducts((prev) =>
      prev.map((p) =>
        p.id === productId ? { ...p, commentsCount: p.commentsCount + 1 } : p
      )
    );
  };

  const openCommentModal = (productId: string) => setActiveCommentProductId(productId);
  const closeCommentModal = () => setActiveCommentProductId(null);
  
  const openDeleteModal = (product: Product) => setDeleteModalProduct(product);
  const closeDeleteModal = () => setDeleteModalProduct(null);
  const confirmDeleteProduct = () => {
    if (deleteModalProduct) {
      deleteProduct(deleteModalProduct.id);
      closeDeleteModal();
    }
  };

  const openWarningModal = (onConfirm: () => void) => {
    setWarningModalConfig({ isOpen: true, onConfirm });
  };
  const closeWarningModal = () => setWarningModalConfig(null);

  return (
    <AppContext.Provider
      value={{
        products,
        currentUser,
        followingMap,
        likesMap,
        commentsMap,
        activeCommentProductId,
        deleteModalProduct,
        warningModalConfig,
        toggleLike,
        toggleFollow,
        addProduct,
        updateProduct,
        deleteProduct,
        addComment,
        openCommentModal,
        closeCommentModal,
        openDeleteModal,
        closeDeleteModal,
        confirmDeleteProduct,
        openWarningModal,
        closeWarningModal,
      }}
    >
      {children}
    </AppContext.Provider>
  );
}

export function useApp() {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error('useApp must be used within an AppProvider');
  }
  return context;
}
