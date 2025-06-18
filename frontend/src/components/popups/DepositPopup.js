import React, { useState } from 'react';

const DepositPopup = ({ onClose }) => {
  const [depositAmount, setDepositAmount] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async () => {
    if (!depositAmount || parseFloat(depositAmount) <= 0) {
      alert('Please enter a valid deposit amount');
      return;
    }

    setIsSubmitting(true);
    
    // Simulate API call
    setTimeout(() => {
      const amount = parseFloat(depositAmount);
      console.log('Deposit Amount:', amount);
      alert(`Successfully deposited $${amount.toLocaleString()}!`);
      setIsSubmitting(false);
      onClose();
    }, 1000);
  };

  const handleAmountChange = (e) => {
    const value = e.target.value;
    // Only allow numbers and decimal point
    if (value === '' || /^\d*\.?\d*$/.test(value)) {
      setDepositAmount(value);
    }
  };

  const quickAmounts = [100, 500, 1000, 5000, 10000];

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl max-w-md w-full">
        {/* Header */}
        <div className="relative p-6 text-center border-b border-gray-100">
          <button
            onClick={onClose}
            className="absolute top-4 right-4 w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 text-gray-500 hover:text-gray-700"
          >
            ×
          </button>
          <h2 className="text-2xl font-bold bg-gradient-to-r from-green-500 to-emerald-500 bg-clip-text text-transparent mb-2">
            Deposit Funds
          </h2>
          <p className="text-sm text-gray-600">
            Add funds to your investment account
          </p>
        </div>

        {/* Content */}
        <div className="p-6">
          {/* Amount Input */}
          <div className="mb-6">
            <label className="block text-sm font-semibold text-gray-700 mb-2">
              Deposit Amount
            </label>
            <div className="relative">
              <span className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500 text-lg">
                $
              </span>
              <input
                type="text"
                value={depositAmount}
                onChange={handleAmountChange}
                placeholder="0.00"
                className="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-lg text-lg focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent"
              />
            </div>
          </div>

          {/* Quick Amount Buttons */}
          <div className="mb-6">
            <label className="block text-sm font-semibold text-gray-700 mb-3">
              Quick Select
            </label>
            <div className="grid grid-cols-3 gap-2">
              {quickAmounts.map((amount) => (
                <button
                  key={amount}
                  onClick={() => setDepositAmount(amount.toString())}
                  className="px-3 py-2 border border-gray-300 rounded-lg text-sm hover:border-green-500 hover:bg-green-50 transition-colors"
                >
                  ${amount.toLocaleString()}
                </button>
              ))}
            </div>
          </div>

          {/* Current Balance Info */}
          <div className="bg-gray-50 rounded-lg p-4 mb-6">
            <div className="flex justify-between items-center text-sm">
              <span className="text-gray-600">Current Balance:</span>
              <span className="font-semibold text-gray-800">$45,230.50</span>
            </div>
            {depositAmount && parseFloat(depositAmount) > 0 && (
              <div className="flex justify-between items-center text-sm mt-2 pt-2 border-t border-gray-200">
                <span className="text-gray-600">New Balance:</span>
                <span className="font-bold text-green-600">
                  ${(45230.50 + parseFloat(depositAmount)).toLocaleString()}
                </span>
              </div>
            )}
          </div>
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
          <button
            onClick={onClose}
            disabled={isSubmitting}
            className="px-5 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={!depositAmount || parseFloat(depositAmount) <= 0 || isSubmitting}
            className="px-5 py-2 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-lg font-semibold hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-200 disabled:opacity-50 disabled:transform-none disabled:shadow-none"
          >
            {isSubmitting ? 'Processing...' : 'Deposit Funds'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default DepositPopup; 