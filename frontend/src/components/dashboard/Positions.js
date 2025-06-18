import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';

// Position Details Popup Component
const PositionDetailsPopup = ({ position, onClose }) => {
  if (!position) return null;

  const formatCurrency = (value) => {
    if (!value) return '$0.00';
    const numValue = parseFloat(value);
    return isNaN(numValue) ? '$0.00' : new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(numValue);
  };

  const formatPercentage = (value) => {
    if (!value) return '0.00%';
    const numValue = parseFloat(value) * 100;
    const sign = numValue >= 0 ? '+' : '';
    return `${sign}${numValue.toFixed(2)}%`;
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold">{position.symbol} Position Details</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            ✕
          </button>
        </div>
        
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-gray-500">Quantity</div>
              <div className="font-medium">{parseFloat(position.qty).toFixed(6)}</div>
            </div>
            <div>
              <div className="text-sm text-gray-500">Side</div>
              <div className="font-medium capitalize">{position.side}</div>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-gray-500">Current Price</div>
              <div className="font-medium">{formatCurrency(position.current_price)}</div>
            </div>
            <div>
              <div className="text-sm text-gray-500">Avg Entry Price</div>
              <div className="font-medium">{formatCurrency(position.avg_entry_price)}</div>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-gray-500">Market Value</div>
              <div className="font-medium text-blue-600">{formatCurrency(position.market_value)}</div>
            </div>
            <div>
              <div className="text-sm text-gray-500">Cost Basis</div>
              <div className="font-medium">{formatCurrency(position.cost_basis)}</div>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-gray-500">Unrealized P&L</div>
              <div className={`font-medium ${parseFloat(position.unrealized_pl) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {formatCurrency(position.unrealized_pl)}
              </div>
            </div>
            <div>
              <div className="text-sm text-gray-500">Unrealized P&L %</div>
              <div className={`font-medium ${parseFloat(position.unrealized_plpc) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {formatPercentage(position.unrealized_plpc)}
              </div>
            </div>
          </div>
          
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-sm text-gray-500">Intraday P&L</div>
              <div className={`font-medium ${parseFloat(position.unrealized_intraday_pl) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {formatCurrency(position.unrealized_intraday_pl)}
              </div>
            </div>
            <div>
              <div className="text-sm text-gray-500">Intraday P&L %</div>
              <div className={`font-medium ${parseFloat(position.unrealized_intraday_plpc) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                {formatPercentage(position.unrealized_intraday_plpc)}
              </div>
            </div>
          </div>
        </div>
        
        <div className="mt-6 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

const Positions = () => {
  const { isAuthenticated } = useAuth();
  const [positions, setPositions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [selectedPosition, setSelectedPosition] = useState(null);
  const [showDetailsPopup, setShowDetailsPopup] = useState(false);

  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  const fetchPositions = async () => {
    if (!isAuthenticated) return;

    try {
      setLoading(true);
      setError(null);

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
      
      // Set positions (filter out positions with zero quantity)
      const filteredPositions = Array.isArray(data.positions) 
        ? data.positions.filter(pos => parseFloat(pos.qty || 0) !== 0)
        : [];
      setPositions(filteredPositions);
      
    } catch (err) {
      console.error('Error fetching positions:', err);
      setError(err.message || 'Failed to fetch positions');
      // Set fallback data on error
      setPositions([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPositions();
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

  const formatQuantity = (qty) => {
    if (!qty) return '0';
    const numValue = parseFloat(qty);
    return isNaN(numValue) ? '0' : numValue.toFixed(3);
  };

  const calculatePL = (position) => {
    const unrealizedPL = parseFloat(position.unrealized_pl || 0);
    return unrealizedPL;
  };

  const formatPL = (pl) => {
    const value = parseFloat(pl);
    if (isNaN(value)) return '$0.00';
    const sign = value >= 0 ? '+' : '';
    return `${sign}${formatCurrency(Math.abs(value))}`;
  };

  const handleViewMore = (position) => {
    setSelectedPosition(position);
    setShowDetailsPopup(true);
  };

  const handleCloseDetails = () => {
    setSelectedPosition(null);
    setShowDetailsPopup(false);
  };

  if (!isAuthenticated) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
        <div className="flex justify-between items-center mb-5">
          <h2 className="text-lg font-semibold">Positions</h2>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <p className="text-gray-500">Please log in to view positions</p>
        </div>
      </div>
    );
  }

  return (
    <>
      <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
        <div className="flex justify-between items-center mb-5">
          <h2 className="text-lg font-semibold">Positions</h2>
          <button 
            className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50"
            onClick={fetchPositions}
            disabled={loading}
          >
            {loading ? 'Loading...' : 'Refresh'}
          </button>
        </div>

        <div className="border border-gray-200 rounded flex-1 min-h-0 overflow-auto">
          <div className="bg-gray-50 px-4 py-3 grid grid-cols-[2fr_1.5fr_1.5fr_80px] gap-4 font-semibold text-sm">
            <div>Symbol</div>
            <div className="text-center">Qty</div>
            <div className="text-center">P&L</div>
            <div className="text-center">Details</div>
          </div>
          
          {loading ? (
            <div className="px-4 py-8 text-center text-gray-500">
              Loading positions...
            </div>
          ) : error ? (
            <div className="px-4 py-8 text-center">
              <p className="text-red-500 mb-2">Error loading positions</p>
              <p className="text-gray-500 text-sm">{error}</p>
              <button 
                onClick={fetchPositions}
                className="mt-2 px-3 py-1 bg-blue-500 text-white rounded text-sm hover:bg-blue-600"
              >
                Retry
              </button>
            </div>
          ) : positions.length === 0 ? (
            <div className="px-4 py-8 text-center text-gray-500">
              No positions found
            </div>
          ) : (
            positions.map((position, index) => {
              const pl = calculatePL(position);
              const isPositive = pl >= 0;
              
              return (
                <div key={index} className="px-4 py-3 grid grid-cols-[2fr_1.5fr_1.5fr_80px] gap-4 border-t border-gray-100 hover:bg-gray-50 items-center">
                  <div className="font-medium">{position.symbol}</div>
                  <div className="text-center">{formatQuantity(position.qty)}</div>
                  <div className={`text-center ${isPositive ? 'text-green-600' : 'text-red-600'}`}>
                    {formatPL(pl)}
                  </div>
                  <div className="flex justify-center">
                    <button
                      onClick={() => handleViewMore(position)}
                      className="w-8 h-8 flex items-center justify-center rounded-full bg-gray-100 hover:bg-blue-100 text-gray-600 hover:text-blue-600 transition-all duration-200 hover:shadow-md"
                      title="View position details"
                    >
                      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2"/>
                        <path d="M12 16v-4" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
                        <path d="M12 8h.01" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
                      </svg>
                    </button>
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>

      {/* Position Details Popup */}
      {showDetailsPopup && (
        <PositionDetailsPopup 
          position={selectedPosition} 
          onClose={handleCloseDetails} 
        />
      )}
    </>
  );
};

export default Positions; 