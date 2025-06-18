import React, { useState } from 'react';
import Sidebar from './components/layout/Sidebar';
import Header from './components/layout/Header';
import Dashboard from './components/dashboard/Dashboard';
import Chatbot from './components/chatbot/Chatbot';
import PandoraPopup from './components/popups/PandoraPopup';
import AuthScreen from './components/auth/AuthScreen';

// Main App Component
const App = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [currentPage, setCurrentPage] = useState('dashboard');
  const [showPandoraPopup, setShowPandoraPopup] = useState(false);

  const handleAuthSuccess = () => {
    setIsAuthenticated(true);
  };

  const handleLogout = () => {
    // Clear any auth tokens or user data from localStorage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setIsAuthenticated(false);
  };

  if (!isAuthenticated) {
    return <AuthScreen onLogin={() => setIsAuthenticated(true)} />;
  }

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Sidebar */}
      <Sidebar onLogout={handleLogout} currentPage={currentPage} setCurrentPage={setCurrentPage} />
      
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

export default App;