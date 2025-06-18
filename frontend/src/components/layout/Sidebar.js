import React, { useState } from 'react';
import UserMenu from './UserMenu';

// Navigation Item Component
const NavItem = ({ icon, label, isActive, onClick }) => (
  <button
    onClick={onClick}
    className={`w-12 h-12 rounded-xl flex items-center justify-center ${
      isActive 
        ? 'bg-gradient-to-br from-purple-500 via-blue-500 to-indigo-500 text-white shadow-lg' 
        : 'bg-white text-gray-400 hover:bg-gray-50 hover:text-purple-600'
    }`}
  >
    {icon}
  </button>
);

// Sidebar Component
const Sidebar = ({ currentPage, setCurrentPage, onLogout }) => {
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  return (
    <div className="w-20 bg-gray-50 border-r border-gray-200 flex flex-col items-center py-5">
      {/* Logo */}
      
      
      {/* Navigation */}
      <nav className="flex flex-col items-center space-y-4">
        <NavItem
          icon="📊"
          label="Dashboard"
          isActive={currentPage === 'dashboard'}
          onClick={() => setCurrentPage('dashboard')}
        />
        <NavItem
          icon={
            <div className="text-2xl font-bold bg-gradient-to-br from-yellow-400 via-purple-500 to-blue-500 bg-clip-text text-transparent transform hover:scale-110 transition-transform duration-200">
              Z
            </div>
          }
          label="Zeus AI Assistant"
          isActive={currentPage === 'ai'}
          onClick={() => setCurrentPage('ai')}
        />
      </nav>
      
      {/* User Menu */}
      <div className="mt-auto relative">
        <button
          onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
          className="w-12 h-12 bg-white border border-gray-200 rounded-full flex items-center justify-center text-sm font-semibold cursor-pointer hover:bg-gray-50"
        >
          JD
        </button>

        <UserMenu
          isOpen={isUserMenuOpen}
          onClose={() => setIsUserMenuOpen(false)}
          onLogout={() => {
            setIsUserMenuOpen(false);
            onLogout();
          }}
        />
      </div>
    </div>
  );
};

export default Sidebar; 