import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import PortfolioChart from '../PortfolioChart';

// Metric Card Component
const MetricCard = ({ label, value, isLoading = false, isPositive = null }) => {
  const getTextColor = () => {
    if (isPositive === null) return 'bg-gradient-to-r from-purple-500 to-blue-500 bg-clip-text text-transparent';
    return isPositive ? 'text-green-600' : 'text-red-600';
  };

  return (
    <div className="text-center p-4 border border-gray-200 rounded-lg">
      <div className={`text-2xl font-bold ${getTextColor()}`}>
        {isLoading ? '...' : value}
      </div>
      <div className="text-xs text-gray-500 mt-1">{label}</div>
    </div>
  );
};

const PortfolioPerformance = () => {
  const { isAuthenticated } = useAuth();
  const [selectedPeriod, setSelectedPeriod] = useState('1D');
  const [performanceData, setPerformanceData] = useState(null);
  const [chartData, setChartData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  // Fetch current portfolio performance (for metrics)
  const fetchPortfolioPerformance = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/portfolio/performance`, {
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
      setPerformanceData(data);
    } catch (err) {
      console.error('Error fetching portfolio performance:', err);
      setError(err.message);
    }
  };

  // Fetch multi-timeframe data for charts
  const fetchChartData = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/portfolio/performance/all`, {
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
      setChartData(data);
    } catch (err) {
      console.error('Error fetching chart data:', err);
      setError(err.message);
    }
  };

  useEffect(() => {
    if (isAuthenticated) {
      const loadData = async () => {
        setLoading(true);
        setError(null);
        try {
          await Promise.all([
            fetchPortfolioPerformance(),
            fetchChartData()
          ]);
        } catch (err) {
          setError('Failed to load portfolio data');
        } finally {
          setLoading(false);
        }
      };

      loadData();
    }
  }, [isAuthenticated]);

  const formatCurrency = (value) => {
    if (value === null || value === undefined) return '$0.00';
    const numValue = typeof value === 'string' ? parseFloat(value) : value;
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
    }).format(numValue);
  };

  const formatPercentage = (value) => {
    if (value === null || value === undefined) return '0.00%';
    const numValue = typeof value === 'string' ? parseFloat(value) : value;
    const sign = numValue >= 0 ? '+' : '';
    return `${sign}${numValue.toFixed(2)}%`;
  };

  if (!isAuthenticated) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-6 flex items-center justify-center">
        <p className="text-gray-500">Please log in to view portfolio performance</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg border border-gray-200 p-6 flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-2">Error loading portfolio data</p>
          <p className="text-gray-500 text-sm">{error}</p>
          <button 
            onClick={() => window.location.reload()} 
            className="mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 flex flex-col">
      <div className="flex justify-between items-center mb-5">
        <h2 className="text-lg font-semibold">Portfolio Performance</h2>
        <div className="flex gap-2">
          {['1D', '1W', '1M', '1Y'].map((period) => (
            <button 
              key={period} 
              onClick={() => setSelectedPeriod(period)}
              className={`px-3 py-1.5 text-xs border rounded hover:bg-gray-50 ${
                selectedPeriod === period 
                  ? 'border-blue-500 bg-blue-50 text-blue-600' 
                  : 'border-gray-200'
              }`}
            >
              {period}
            </button>
          ))}
        </div>
      </div>
      
      {/* Metrics */}
      <div className="grid grid-cols-3 gap-5 mb-5">
        <MetricCard 
          label="Total Value" 
          value={formatCurrency(performanceData?.total_market_value)} 
          isLoading={loading}
        />
        <MetricCard 
          label="Today's Change" 
          value={formatCurrency(performanceData?.daily_pl)} 
          isLoading={loading}
          isPositive={performanceData?.daily_pl >= 0}
        />
        <MetricCard 
          label="Total Return" 
          value={formatPercentage(performanceData?.total_plpc)} 
          isLoading={loading}
          isPositive={performanceData?.total_plpc >= 0}
        />
      </div>
      
      {/* Chart */}
      <div className="flex-1 min-h-0">
        <PortfolioChart 
          data={chartData} 
          selectedPeriod={selectedPeriod} 
          loading={loading}
        />
      </div>
    </div>
  );
};

export default PortfolioPerformance; 