import React from 'react';

const RecentOrders = () => {
  const orders = [
    { symbol: 'AAPL', type: 'Buy', quantity: 100, status: 'Filled' },
    { symbol: 'TSLA', type: 'Sell', quantity: 50, status: 'Pending' },
    { symbol: 'MSFT', type: 'Buy', quantity: 75, status: 'Filled' },
  ];

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
      <div className="flex justify-between items-center mb-5">
        <h2 className="text-lg font-semibold">Recent Orders</h2>
        <div className="flex gap-2">
          <button className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50">View All</button>
          <button className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50">+ New Order</button>
        </div>
      </div>
      <div className="border border-gray-200 rounded flex-1 min-h-0 overflow-auto">
        <div className="bg-gray-50 px-4 py-3 grid grid-cols-4 gap-4 font-semibold text-sm">
          <div>Symbol</div>
          <div>Type</div>
          <div>Quantity</div>
          <div>Status</div>
        </div>
        {orders.map((order, index) => (
          <div key={index} className="px-4 py-3 grid grid-cols-4 gap-4 border-t border-gray-100">
            <div>{order.symbol}</div>
            <div>{order.type}</div>
            <div>{order.quantity}</div>
            <div>{order.status}</div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default RecentOrders; 