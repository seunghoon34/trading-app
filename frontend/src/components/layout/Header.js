import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import DepositPopup from '../popups/DepositPopup';

const Header = ({ onOpenPandora }) => {
  const { isAuthenticated } = useAuth();
  const [showDepositPopup, setShowDepositPopup] = useState(false);
  const [cashBalance, setCashBalance] = useState(null);
  const [loading, setLoading] = useState(false);

  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  const fetchCashBalance = async () => {
    if (!isAuthenticated) return;

    try {
      setLoading(true);
      const response = await fetch(`${API_BASE_URL}/api/v1/portfolio/positions`, {
        method: 'GET',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setCashBalance(data.Cash || '0');
    } catch (err) {
      console.error('Error fetching cash balance:', err);
      setCashBalance('0');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCashBalance();
  }, [isAuthenticated]);

  const formatCurrency = (value) => {
    if (!value) return '$0.00';
    const numValue = parseFloat(value);
    return isNaN(numValue) ? '$0.00' : new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(numValue);
  };

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
            <div className="text-base font-bold text-green-600">
              {loading ? '...' : formatCurrency(cashBalance)}
            </div>
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