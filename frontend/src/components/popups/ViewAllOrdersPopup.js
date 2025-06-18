import React, { useState, useEffect } from 'react';

const ViewAllOrdersPopup = ({ isOpen, onClose }) => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  const fetchAllOrders = async () => {
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
      setOrders(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error('Error fetching orders:', err);
      setError(err.message || 'Failed to fetch orders');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isOpen) {
      fetchAllOrders();
    }
  }, [isOpen]);

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    try {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return 'N/A';
    }
  };

  const formatCurrency = (value) => {
    if (!value) return 'N/A';
    const numValue = parseFloat(value);
    return isNaN(numValue) ? 'N/A' : `$${numValue.toFixed(2)}`;
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'filled':
        return 'bg-green-100 text-green-800';
      case 'pending':
      case 'submitted':
      case 'accepted':
        return 'bg-yellow-100 text-yellow-800';
      case 'cancelled':
      case 'rejected':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-6xl max-h-[90vh] m-4">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-800">All Orders</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 text-2xl font-bold"
          >
            ×
          </button>
        </div>

        {/* Content */}
        <div className="p-6 overflow-auto max-h-[70vh]">
          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-gray-500">Loading orders...</div>
            </div>
          ) : error ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-center">
                <p className="text-red-500 mb-2">Error loading orders</p>
                <p className="text-gray-500 text-sm">{error}</p>
                <button 
                  onClick={fetchAllOrders}
                  className="mt-3 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                >
                  Retry
                </button>
              </div>
            </div>
          ) : orders.length === 0 ? (
            <div className="flex items-center justify-center py-12">
              <div className="text-gray-500">No orders found</div>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="bg-gray-50">
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border border-gray-200">Symbol</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border border-gray-200">Side</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border border-gray-200">Quantity</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border border-gray-200">Notional</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border border-gray-200">Status</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border border-gray-200">Date</th>
                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border border-gray-200">Order ID</th>
                  </tr>
                </thead>
                <tbody>
                  {orders.map((order, index) => (
                    <tr key={order.id || index} className="hover:bg-gray-50">
                      <td className="px-4 py-3 text-sm border border-gray-200 font-medium">
                        {order.symbol || 'N/A'}
                      </td>
                      <td className="px-4 py-3 text-sm border border-gray-200">
                        <span className={`px-2 py-1 rounded text-xs font-medium ${
                          order.side?.toLowerCase() === 'buy' 
                            ? 'bg-green-100 text-green-800' 
                            : 'bg-red-100 text-red-800'
                        }`}>
                          {order.side?.toUpperCase() || 'N/A'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-sm border border-gray-200">
                        {order.qty || 'N/A'}
                      </td>
                      <td className="px-4 py-3 text-sm border border-gray-200">
                        {formatCurrency(order.notional)}
                      </td>
                      <td className="px-4 py-3 text-sm border border-gray-200">
                        <span className={`px-2 py-1 rounded text-xs font-medium ${getStatusColor(order.status)}`}>
                          {order.status || 'N/A'}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-sm border border-gray-200">
                        {formatDate(order.created_at)}
                      </td>
                      <td className="px-4 py-3 text-sm border border-gray-200 font-mono text-xs">
                        {order.id || 'N/A'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex justify-end p-6 border-t border-gray-200">
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

export default ViewAllOrdersPopup; 