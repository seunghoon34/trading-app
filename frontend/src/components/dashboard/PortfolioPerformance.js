import React from 'react';
import PortfolioChart from '../PortfolioChart';

// Metric Card Component
const MetricCard = ({ label, value }) => {
  return (
    <div className="text-center p-4 border border-gray-200 rounded-lg">
      <div className="text-2xl font-bold bg-gradient-to-r from-purple-500 to-blue-500 bg-clip-text text-transparent">
        {value}
      </div>
      <div className="text-xs text-gray-500 mt-1">{label}</div>
    </div>
  );
};

const PortfolioPerformance = () => {
  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 flex flex-col">
      <div className="flex justify-between items-center mb-5">
        <h2 className="text-lg font-semibold">Portfolio Performance</h2>
        <div className="flex gap-2">
          {['1D', '1W', '1M', '1Y'].map((period) => (
            <button key={period} className="px-3 py-1.5 text-xs border border-gray-200 rounded hover:bg-gray-50">
              {period}
            </button>
          ))}
        </div>
      </div>
      {/* Metrics */}
      <div className="grid grid-cols-3 gap-5 mb-5">
        <MetricCard label="Total Value" value="$124,567" />
        <MetricCard label="Today's Change" value="+$8,432" />
        <MetricCard label="Total Return" value="+12.4%" />
      </div>
      {/* Chart */}
      <div className="flex-1 min-h-0">
        <PortfolioChart/>
      </div>
    </div>
  );
};

export default PortfolioPerformance; 