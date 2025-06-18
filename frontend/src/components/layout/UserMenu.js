import React from 'react';

const UserMenu = ({ isOpen, onClose, onLogout }) => {
  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div 
        className="fixed inset-0 bg-black bg-opacity-50 z-40 flex items-center justify-center"
        onClick={onClose}
      >
        {/* Menu */}
        <div 
          className="w-[400px] rounded-2xl bg-white shadow-xl z-50"
          onClick={e => e.stopPropagation()}
        >
          {/* User Info */}
          <div className="p-6 border-b border-gray-100">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-full bg-gradient-to-r from-purple-500 to-blue-500 flex items-center justify-center text-white font-semibold">
                JD
              </div>
              <div>
                <div className="font-medium text-xl text-gray-900">John Doe</div>
                <div className="text-sm text-gray-500">john@example.com</div>
              </div>
            </div>
          </div>

          {/* Menu Items */}
          <div className="p-4">
            <button
              className="w-full text-left px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 rounded-lg flex items-center gap-3 mb-2"
              onClick={() => {}}
            >
              <svg className="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              Profile Settings
            </button>

            <button
              className="w-full text-left px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 rounded-lg flex items-center gap-3 mb-2"
              onClick={() => {}}
            >
              <svg className="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              Account Settings
            </button>

            <hr className="my-3 border-gray-200" />

            <button
              className="w-full text-left px-4 py-3 text-sm text-red-600 hover:bg-red-50 rounded-lg flex items-center gap-3"
              onClick={onLogout}
            >
              <svg className="h-5 w-5 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Sign out
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default UserMenu; 