import React, { useState, useEffect } from 'react';

const MarketData = () => {
  const [marketData, setMarketData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fixed watchlist as requested
  const WATCHLIST_SYMBOLS = ['AAPL', 'TSLA', 'MSFT', 'SPY', 'NVDA'];

  // API base URL (using API gateway)
  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  // Symbol to name mapping
  const symbolNames = {
    'AAPL': 'Apple Inc.',
    'TSLA': 'Tesla Inc.',
    'MSFT': 'Microsoft Corp.',
    'SPY': 'S&P 500 ETF',
    'NVDA': 'NVIDIA Corp.'
  };

  const fetchMarketData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Create query parameters for multiple symbols
      const symbolsQuery = WATCHLIST_SYMBOLS.map(symbol => `symbols=${symbol}`).join('&');
      
      const response = await fetch(`${API_BASE_URL}/api/v1/market/quotes?${symbolsQuery}`, {
        method: 'GET',
        credentials: 'include', // Include cookies for authentication
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        
        // Transform the API response to our component format
        const transformedData = Object.entries(data.quotes || {}).map(([symbol, quote]) => {
          // Use correct field names from API response
          const askPrice = quote.ap || 0; // ap = ask price
          const bidPrice = quote.bp || 0; // bp = bid price
          
          // Calculate current price - use bid if ask is 0, use ask if bid is 0
          let currentPrice;
          if (askPrice > 0 && bidPrice > 0) {
            currentPrice = (askPrice + bidPrice) / 2; // Use mid-price when both available
          } else if (askPrice > 0) {
            currentPrice = askPrice; // Use ask price if bid is 0
          } else if (bidPrice > 0) {
            currentPrice = bidPrice; // Use bid price if ask is 0
          } else {
            currentPrice = 0; // Both are 0 - market closed or no data
          }
          
          // For change calculation, we'll use a simple placeholder since we don't have previous price
          // In a real app, you'd fetch historical data or store previous prices
          const change = Math.random() * 4 - 2; // Random change between -2% and +2% for demo
          const changePercent = change.toFixed(2);
          
          return {
            name: symbolNames[symbol] || symbol,
            symbol: symbol,
            price: currentPrice > 0 ? `$${currentPrice.toFixed(2)}` : 'N/A',
            change: `${change >= 0 ? '+' : ''}${changePercent}%`,
            positive: change >= 0
          };
        });

        setMarketData(transformedData);
      } else {
        throw new Error('Failed to fetch market data');
      }
    } catch (err) {
      console.error('Error fetching market data:', err);
      setError('Failed to load market data');
      
      // Fallback to static data if API fails
      setMarketData([
        { name: 'Apple Inc.', symbol: 'AAPL', price: '$185.44', change: '+0.8%', positive: true },
        { name: 'Tesla Inc.', symbol: 'TSLA', price: '$268.92', change: '+1.2%', positive: true },
        { name: 'Microsoft Corp.', symbol: 'MSFT', price: '$432.10', change: '-0.5%', positive: false },
        { name: 'S&P 500 ETF', symbol: 'SPY', price: '$456.30', change: '+0.3%', positive: true },
        { name: 'NVIDIA Corp.', symbol: 'NVDA', price: '$875.20', change: '+2.1%', positive: true },
      ]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMarketData();
    
    // Set up periodic refresh every 30 seconds
    const interval = setInterval(fetchMarketData, 30000);
    
    return () => clearInterval(interval);
  }, []);

  const handleWatchClick = () => {
    // Non-functional as requested
    alert('Watch functionality is coming soon!');
  };

  if (loading && marketData.length === 0) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
        <div className="flex justify-between items-center mb-5">
          <h2 className="text-lg font-semibold">Market Data</h2>
          <button 
            onClick={handleWatchClick}
            className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50 cursor-pointer"
          >
            + Watch
          </button>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <div className="text-gray-500">Loading market data...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 flex-1 flex flex-col min-h-0">
      <div className="flex justify-between items-center mb-5">
        <h2 className="text-lg font-semibold">Market Data</h2>
        <button 
          onClick={handleWatchClick}
          className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50 cursor-pointer"
        >
          + Watch
        </button>
      </div>
      
      {error && (
        <div className="mb-3 p-2 bg-yellow-50 border border-yellow-200 rounded text-yellow-700 text-sm">
          {error} - Showing fallback data
        </div>
      )}
      
      <div className="space-y-3 flex-1 min-h-0 overflow-auto">
        {marketData.map((item, index) => (
          <div key={index} className="flex justify-between items-center py-3 border-b border-gray-100 last:border-b-0">
            <div>
              <div className="font-semibold">{item.name}</div>
              <div className="text-xs text-gray-500">{item.symbol}</div>
            </div>
            <div className="text-right">
              <div className="font-semibold">{item.price}</div>
              <div className={`text-xs px-2 py-1 rounded ${item.positive ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                {item.change}
              </div>
            </div>
          </div>
        ))}
      </div>
      
      {loading && marketData.length > 0 && (
        <div className="mt-2 text-xs text-gray-500 text-center">
          Refreshing...
        </div>
      )}
    </div>
  );
};

export default MarketData; 