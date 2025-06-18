import React, { useState } from 'react';
import Sidebar from './components/layout/Sidebar';
import Header from './components/layout/Header';
import Dashboard from './components/dashboard/Dashboard';
import Chatbot from './components/chatbot/Chatbot';
import PandoraPopup from './components/popups/PandoraPopup';
import AuthScreen from './components/auth/AuthScreen';
import { AuthProvider, useAuth } from './contexts/AuthContext';

// Main App Component Content
const AppContent = () => {
  const { isAuthenticated, loading, logout, user } = useAuth();
  const [currentPage, setCurrentPage] = useState('dashboard');
  const [showPandoraPopup, setShowPandoraPopup] = useState(false);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <AuthScreen />;
  }

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Sidebar */}
      <Sidebar onLogout={logout} currentPage={currentPage} setCurrentPage={setCurrentPage} user={user} />
      
      {/* Main Content */}
      <div className="flex-1 flex flex-col h-screen overflow-hidden">
        <Header onOpenPandora={() => setShowPandoraPopup(true)} />
        
        {/* Page Content */}
        <div className="flex-1 overflow-hidden">
          {currentPage === 'dashboard' ? <Dashboard /> : <Chatbot />}
        </div>
      </div>

      {/* Pandora Popup */}
      {showPandoraPopup && (
        <PandoraPopup onClose={() => setShowPandoraPopup(false)} />
      )}
    </div>
  );
};

// Main App Component with AuthProvider
const App = () => {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
};

export default App;