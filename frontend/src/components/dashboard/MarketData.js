import React from 'react';

const MarketData = () => {
  const marketData = [
    { name: 'S&P 500', symbol: 'SPY', price: '$4,185.44', change: '+0.8%', positive: true },
    { name: 'NASDAQ', symbol: 'QQQ', price: '$368.92', change: '+1.2%', positive: true },
    { name: 'Bitcoin', symbol: 'BTC-USD', price: '$67,432', change: '-2.1%', positive: false },
    { name: 'Gold', symbol: 'GLD', price: '$2,056.30', change: '+0.3%', positive: true },
  ];

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
      <div className="flex justify-between items-center mb-5">
        <h2 className="text-lg font-semibold">Market Data</h2>
        <button className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50">+ Watch</button>
      </div>
      <div className="space-y-3 flex-1 min-h-0 overflow-auto">
        {marketData.map((item, index) => (
          <div key={index} className="flex justify-between items-center py-3 border-b border-gray-100 last:border-b-0">
            <div>
              <div className="font-semibold">{item.name}</div>
              <div className="text-xs text-gray-500">{item.symbol}</div>
            </div>
            <div className="text-right">
              <div className="font-semibold">{item.price}</div>
              <div className={`text-xs px-2 py-1 rounded ${item.positive ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>{item.change}</div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default MarketData; 