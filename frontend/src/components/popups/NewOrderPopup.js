import React, { useState } from 'react';

const NewOrderPopup = ({ isOpen, onClose, onOrderPlaced }) => {
  const [formData, setFormData] = useState({
    symbol: '',
    side: 'buy',
    orderType: 'qty', // 'qty' or 'notional'
    qty: '',
    notional: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    setError(null);
  };

  const handleOrderTypeChange = (type) => {
    setFormData(prev => ({
      ...prev,
      orderType: type,
      qty: type === 'qty' ? prev.qty : '',
      notional: type === 'notional' ? prev.notional : ''
    }));
  };

  const validateForm = () => {
    if (!formData.symbol.trim()) {
      setError('Symbol is required');
      return false;
    }

    if (formData.orderType === 'qty') {
      if (!formData.qty || parseFloat(formData.qty) <= 0) {
        setError('Valid quantity is required');
        return false;
      }
    } else {
      if (!formData.notional || parseFloat(formData.notional) <= 0) {
        setError('Valid notional amount is required');
        return false;
      }
    }

    return true;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const orderData = {
        symbol: formData.symbol.toUpperCase().trim(),
        side: formData.side
      };

      // Add either qty or notional, not both
      if (formData.orderType === 'qty') {
        orderData.qty = formData.qty;
      } else {
        orderData.notional = formData.notional;
      }

      const response = await fetch(`${API_BASE_URL}/api/v1/trading/orders`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(orderData)
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      setSuccess(true);
      
      // Call callback to refresh orders list
      if (onOrderPlaced) {
        onOrderPlaced(result);
      }

      // Auto-close after 2 seconds
      setTimeout(() => {
        handleClose();
      }, 2000);

    } catch (err) {
      console.error('Error placing order:', err);
      setError(err.message || 'Failed to place order');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      symbol: '',
      side: 'buy',
      orderType: 'qty',
      qty: '',
      notional: ''
    });
    setError(null);
    setSuccess(false);
    setLoading(false);
    onClose();
  };

  if (!isOpen) return null;

  if (success) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div className="bg-white rounded-lg shadow-xl w-full max-w-md m-4 p-6">
          <div className="text-center">
            <div className="w-16 h-16 mx-auto mb-4 bg-green-100 rounded-full flex items-center justify-center">
              <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7"></path>
              </svg>
            </div>
            <h3 className="text-lg font-semibold text-gray-800 mb-2">Order Placed Successfully!</h3>
            <p className="text-gray-600 mb-4">Your order has been submitted to the market.</p>
            <button
              onClick={handleClose}
              className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md m-4">
        {/* Header */}
        <div className="flex justify-between items-center p-6 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-800">Place New Order</h2>
          <button
            onClick={handleClose}
            className="text-gray-400 hover:text-gray-600 text-2xl font-bold"
          >
            ×
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6">
          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded text-red-700 text-sm">
              {error}
            </div>
          )}

          {/* Symbol */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Symbol *
            </label>
            <input
              type="text"
              name="symbol"
              value={formData.symbol}
              onChange={handleInputChange}
              placeholder="e.g., AAPL"
              className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent uppercase"
              required
            />
          </div>

          {/* Side */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Side *
            </label>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => setFormData(prev => ({ ...prev, side: 'buy' }))}
                className={`flex-1 px-4 py-2 rounded font-medium ${
                  formData.side === 'buy'
                    ? 'bg-green-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                Buy
              </button>
              <button
                type="button"
                onClick={() => setFormData(prev => ({ ...prev, side: 'sell' }))}
                className={`flex-1 px-4 py-2 rounded font-medium ${
                  formData.side === 'sell'
                    ? 'bg-red-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                Sell
              </button>
            </div>
          </div>

          {/* Order Type */}
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Order Type *
            </label>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => handleOrderTypeChange('qty')}
                className={`flex-1 px-4 py-2 rounded font-medium ${
                  formData.orderType === 'qty'
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                Quantity
              </button>
              <button
                type="button"
                onClick={() => handleOrderTypeChange('notional')}
                className={`flex-1 px-4 py-2 rounded font-medium ${
                  formData.orderType === 'notional'
                    ? 'bg-blue-500 text-white'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                Dollar Amount
              </button>
            </div>
          </div>

          {/* Quantity or Notional */}
          {formData.orderType === 'qty' ? (
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Quantity *
              </label>
              <input
                type="number"
                name="qty"
                value={formData.qty}
                onChange={handleInputChange}
                placeholder="Number of shares"
                min="1"
                step="1"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>
          ) : (
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Dollar Amount *
              </label>
              <div className="relative">
                <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
                <input
                  type="number"
                  name="notional"
                  value={formData.notional}
                  onChange={handleInputChange}
                  placeholder="0.00"
                  min="1"
                  step="0.01"
                  className="w-full pl-8 pr-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>
            </div>
          )}

          {/* Buttons */}
          <div className="flex gap-3">
            <button
              type="button"
              onClick={handleClose}
              className="flex-1 px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
              disabled={loading}
            >
              {loading ? 'Placing...' : `${formData.side === 'buy' ? 'Buy' : 'Sell'} ${formData.symbol || 'Stock'}`}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default NewOrderPopup; 