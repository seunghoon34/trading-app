import React from 'react';

const Positions = () => {
  const positions = [
    { symbol: 'AAPL', qty: 200, pl: '+$2,450', positive: true },
    { symbol: 'TSLA', qty: 50, pl: '-$890', positive: false },
    { symbol: 'MSFT', qty: 150, pl: '+$1,230', positive: true },
    { symbol: 'GOOGL', qty: 25, pl: '+$567', positive: true },
  ];

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
      <div className="flex justify-between items-center mb-5">
        <h2 className="text-lg font-semibold">Positions</h2>
        <button className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50">Manage</button>
      </div>
      <div className="border border-gray-200 rounded flex-1 min-h-0 overflow-auto">
        <div className="bg-gray-50 px-4 py-3 grid grid-cols-3 gap-4 font-semibold text-sm">
          <div>Symbol</div>
          <div>Qty</div>
          <div>P&L</div>
        </div>
        {positions.map((position, index) => (
          <div key={index} className="px-4 py-3 grid grid-cols-3 gap-4 border-t border-gray-100">
            <div>{position.symbol}</div>
            <div>{position.qty}</div>
            <div className={position.positive ? 'text-green-600' : 'text-red-600'}>
              {position.pl}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Positions; 