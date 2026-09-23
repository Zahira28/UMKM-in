'use client';

import React, { useState } from 'react';
import { Send, X, User } from 'lucide-react';
import { useApp } from '../../context/AppContext';

export default function CommentsModal() {
  const { activeCommentProductId, closeCommentModal, commentsMap, addComment } = useApp();
  const [newCommentText, setNewCommentText] = useState('');

  if (!activeCommentProductId) return null;

  const currentComments = commentsMap[activeCommentProductId] || [];

  const handleSend = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCommentText.trim()) return;
    addComment(activeCommentProductId, newCommentText);
    setNewCommentText('');
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-in fade-in duration-150">
      {/* Backdrop click */}
      <div className="absolute inset-0" onClick={closeCommentModal} />

      {/* Modal Container */}
      <div className="relative w-full max-w-lg bg-white rounded-3xl p-6 shadow-2xl z-10 space-y-6 max-h-[85vh] flex flex-col animate-in zoom-in-95 duration-200">
        {/* Modal Header */}
        <div className="flex items-center justify-between border-b border-gray-100 pb-4">
          <h2 className="text-2xl font-extrabold text-gray-900 tracking-tight">
            Komentar
          </h2>
          <button
            onClick={closeCommentModal}
            className="p-2 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Comments Scrollable Area */}
        <div className="flex-1 overflow-y-auto space-y-5 pr-2">
          {currentComments.length === 0 ? (
            <div className="text-center py-10 text-gray-400 text-sm">
              Belum ada komentar. Jadilah yang pertama berkomentar!
            </div>
          ) : (
            currentComments.map((comment) => (
              <div key={comment.id} className="space-y-3">
                {/* Main Comment */}
                <div className="flex gap-3">
                  <div className="w-9 h-9 rounded-full bg-gray-100 overflow-hidden border border-gray-200 flex-shrink-0">
                    {comment.userAvatar ? (
                      <img
                        src={comment.userAvatar}
                        alt={comment.username}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-gray-400">
                        <User className="w-4 h-4" />
                      </div>
                    )}
                  </div>
                  <div className="flex-1">
                    <h4 className="text-sm font-bold text-gray-900">
                      {comment.username}
                    </h4>
                    <p className="text-sm text-gray-700 mt-0.5 leading-relaxed">
                      {comment.content}
                    </p>
                    <button className="text-xs font-semibold text-[#0d6e42] hover:underline mt-1">
                      Balas
                    </button>

                    {/* Replies */}
                    {comment.replies && comment.replies.length > 0 && (
                      <div className="mt-3 space-y-3 pl-4 border-l-2 border-gray-100">
                        {comment.replies.map((reply) => (
                          <div key={reply.id} className="flex gap-2.5">
                            <div className="w-7 h-7 rounded-full bg-gray-100 overflow-hidden border border-gray-200 flex-shrink-0">
                              {reply.userAvatar ? (
                                <img
                                  src={reply.userAvatar}
                                  alt={reply.username}
                                  className="w-full h-full object-cover"
                                />
                              ) : (
                                <div className="w-full h-full flex items-center justify-center text-gray-400">
                                  <User className="w-3 h-3" />
                                </div>
                              )}
                            </div>
                            <div>
                              <h5 className="text-xs font-bold text-gray-900">
                                {reply.username}
                              </h5>
                              <p className="text-xs text-gray-700 mt-0.5">
                                {reply.content}
                              </p>
                              <button className="text-[11px] font-semibold text-[#0d6e42] hover:underline mt-0.5">
                                Balas
                              </button>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Comment Input Footer */}
        <form onSubmit={handleSend} className="pt-2">
          <div className="relative flex items-center bg-[#f6f8f5] border border-gray-200 rounded-2xl px-4 py-2.5 focus-within:border-[#0d6e42] focus-within:ring-1 focus-within:ring-[#0d6e42] transition-all">
            <input
              type="text"
              value={newCommentText}
              onChange={(e) => setNewCommentText(e.target.value)}
              placeholder="Buat komentar..."
              className="w-full bg-transparent text-sm text-gray-900 focus:outline-none pr-10"
            />
            <button
              type="submit"
              disabled={!newCommentText.trim()}
              className="absolute right-3 p-1.5 text-[#0d6e42] hover:bg-[#e6f4ea] rounded-xl transition-colors disabled:opacity-40 disabled:hover:bg-transparent"
            >
              <Send className="w-5 h-5 transform rotate-45" />
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
