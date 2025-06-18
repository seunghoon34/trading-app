import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import ViewAllOrdersPopup from '../popups/ViewAllOrdersPopup';
import NewOrderPopup from '../popups/NewOrderPopup';

const RecentOrders = () => {
  const { isAuthenticated } = useAuth();
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showAllOrdersPopup, setShowAllOrdersPopup] = useState(false);
  const [showNewOrderPopup, setShowNewOrderPopup] = useState(false);

  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  const fetchRecentOrders = async () => {
    if (!isAuthenticated) return;

    try {
      setLoading(true);
      setError(null);

      const response = await fetch(`${API_BASE_URL}/api/v1/trading/orders`, {
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
      // Take only the 3 most recent orders
      const recentOrders = Array.isArray(data) ? data.slice(0, 3) : [];
      setOrders(recentOrders);
    } catch (err) {
      console.error('Error fetching orders:', err);
      setError(err.message || 'Failed to fetch orders');
      // Fallback to static data on error
      setOrders([
        { id: '1', symbol: 'AAPL', side: 'buy', qty: '100', status: 'filled', created_at: new Date().toISOString() },
        { id: '2', symbol: 'TSLA', side: 'sell', qty: '50', status: 'pending', created_at: new Date().toISOString() },
        { id: '3', symbol: 'MSFT', side: 'buy', qty: '75', status: 'filled', created_at: new Date().toISOString() },
      ]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRecentOrders();
  }, [isAuthenticated]);

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    try {
      return new Date(dateString).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return 'N/A';
    }
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'filled':
        return 'text-green-600';
      case 'pending':
      case 'submitted':
      case 'accepted':
        return 'text-yellow-600';
      case 'cancelled':
      case 'rejected':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const handleOrderPlaced = () => {
    // Refresh the orders list after a new order is placed
    fetchRecentOrders();
  };

  if (!isAuthenticated) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
        <div className="flex justify-between items-center mb-5">
          <h2 className="text-lg font-semibold">Recent Orders</h2>
          <div className="flex gap-2">
            <button className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50 opacity-50 cursor-not-allowed">View All</button>
            <button className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50 opacity-50 cursor-not-allowed">+ New Order</button>
          </div>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <p className="text-gray-500">Please log in to view orders</p>
        </div>
      </div>
    );
  }

  return (
    <>
      <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
        <div className="flex justify-between items-center mb-5">
          <h2 className="text-lg font-semibold">Recent Orders</h2>
          <div className="flex gap-2">
            <button 
              onClick={() => setShowAllOrdersPopup(true)}
              className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50"
            >
              View All
            </button>
            <button 
              onClick={() => setShowNewOrderPopup(true)}
              className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50"
            >
              + New Order
            </button>
          </div>
        </div>

        {error && (
          <div className="mb-3 p-2 bg-yellow-50 border border-yellow-200 rounded text-yellow-700 text-sm">
            {error} - Showing fallback data
          </div>
        )}

        <div className="border border-gray-200 rounded flex-1 min-h-0 overflow-auto">
          <div className="bg-gray-50 px-4 py-3 grid grid-cols-4 gap-4 font-semibold text-sm">
            <div>Symbol</div>
            <div>Side</div>
            <div>Quantity</div>
            <div>Status</div>
          </div>
          
          {loading ? (
            <div className="flex items-center justify-center py-8">
              <div className="text-gray-500">Loading orders...</div>
            </div>
          ) : orders.length === 0 ? (
            <div className="flex items-center justify-center py-8">
              <div className="text-gray-500">No orders found</div>
            </div>
          ) : (
            orders.map((order, index) => (
              <div key={order.id || index} className="px-4 py-3 grid grid-cols-4 gap-4 border-t border-gray-100 hover:bg-gray-50">
                <div className="font-medium">{order.symbol || 'N/A'}</div>
                <div>
                  <span className={`px-2 py-1 rounded text-xs font-medium ${
                    order.side?.toLowerCase() === 'buy' 
                      ? 'bg-green-100 text-green-800' 
                      : 'bg-red-100 text-red-800'
                  }`}>
                    {order.side?.toUpperCase() || 'N/A'}
                  </span>
                </div>
                <div>{order.qty || order.notional || 'N/A'}</div>
                <div className={`font-medium ${getStatusColor(order.status)}`}>
                  {order.status || 'N/A'}
                </div>
              </div>
            ))
          )}
        </div>

        {loading && orders.length > 0 && (
          <div className="mt-2 text-xs text-gray-500 text-center">
            Refreshing...
          </div>
        )}
      </div>

      {/* Popups */}
      <ViewAllOrdersPopup 
        isOpen={showAllOrdersPopup} 
        onClose={() => setShowAllOrdersPopup(false)} 
      />
      <NewOrderPopup 
        isOpen={showNewOrderPopup} 
        onClose={() => setShowNewOrderPopup(false)}
        onOrderPlaced={handleOrderPlaced}
      />
    </>
  );
};

export default RecentOrders; 