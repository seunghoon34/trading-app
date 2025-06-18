import React, { useState } from 'react';
import DepositPopup from '../popups/DepositPopup';

const Header = ({ onOpenPandora }) => {
  const [showDepositPopup, setShowDepositPopup] = useState(false);

  return (
    <>
      <header className="bg-white border-b border-gray-200 px-8 py-5 flex justify-between items-center">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-r from-purple-500 via-blue-500 to-indigo-500 flex items-center justify-center text-white font-bold text-xl">
            P
          </div>
          <div className="text-xl font-bold bg-gradient-to-r from-purple-500 via-blue-500 to-indigo-500 bg-clip-text text-transparent">
            Pandora Wealth
          </div>
        </div>
        
        <div className="flex items-center gap-4">
          {/* Cash Balance */}
          <div className="bg-white border border-gray-200 rounded-lg px-4 py-3 text-center">
            <div className="text-xs text-gray-500 mb-1">Cash Balance</div>
            <div className="text-base font-bold text-green-600">$45,230.50</div>
          </div>
          
          {/* Deposit Button */}
          <button
            onClick={() => setShowDepositPopup(true)}
            className="px-5 py-2.5 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-lg font-semibold hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-200"
          >
            Deposit
          </button>
          
          {/* Pandora Button */}
          <button
            onClick={onOpenPandora}
            className="px-5 py-2.5 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-lg font-semibold hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-200"
          >
            Pandora
          </button>
        </div>
      </header>

      {/* Deposit Popup */}
      {showDepositPopup && (
        <DepositPopup onClose={() => setShowDepositPopup(false)} />
      )}
    </>
  );
};

export default Header; 