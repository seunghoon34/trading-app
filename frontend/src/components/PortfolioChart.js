import React, { useMemo } from 'react';

const PortfolioChart = ({ data, selectedPeriod, loading }) => {
  // Process the data for the selected period
  const chartData = useMemo(() => {
    if (!data || !data[selectedPeriod] || loading) {
      return null;
    }

    const periodData = data[selectedPeriod];
    const { timestamp, equity } = periodData;

    if (!timestamp || !equity || timestamp.length === 0 || equity.length === 0) {
      return null;
    }

    // Convert timestamps to dates and combine with equity values
    return timestamp.map((ts, index) => ({
      timestamp: ts,
      date: new Date(ts * 1000), // Convert Unix timestamp to Date
      value: equity[index],
    }));
  }, [data, selectedPeriod, loading]);

  // Calculate min and max for scaling
  const { minValue, maxValue, range } = useMemo(() => {
    if (!chartData || chartData.length === 0) {
      return { minValue: 0, maxValue: 100000, range: 100000 };
    }

    const values = chartData.map(d => d.value);
    const min = Math.min(...values);
    const max = Math.max(...values);
    const calculatedRange = max - min;
    
    return {
      minValue: min,
      maxValue: max,
      range: calculatedRange > 0 ? calculatedRange : 1, // Prevent division by zero
    };
  }, [chartData]);

  // Determine if the overall trend is positive
  const isPositiveTrend = useMemo(() => {
    if (!chartData || chartData.length < 2) return true;
    const firstValue = chartData[0].value;
    const lastValue = chartData[chartData.length - 1].value;
    return lastValue >= firstValue;
  }, [chartData]);

  // Create SVG path for the line
  const createPath = () => {
    if (!chartData || chartData.length === 0) return '';

    const width = 100; // percentage
    const height = 70; // Leave some margin
    const stepX = chartData.length > 1 ? width / (chartData.length - 1) : 0;
    
    return chartData.map((point, index) => {
      const x = index * stepX;
      const y = 15 + (70 - ((point.value - minValue) / range) * height); // 15% margin on top
      return `${index === 0 ? 'M' : 'L'} ${x} ${y}`;
    }).join(' ');
  };

  // Create area path for colored shadow
  const createAreaPath = () => {
    if (!chartData || chartData.length === 0) return '';

    const width = 100;
    const height = 70;
    const stepX = chartData.length > 1 ? width / (chartData.length - 1) : 0;
    
    let path = chartData.map((point, index) => {
      const x = index * stepX;
      const y = 15 + (70 - ((point.value - minValue) / range) * height);
      return `${index === 0 ? 'M' : 'L'} ${x} ${y}`;
    }).join(' ');
    
    // Close the area to the bottom
    path += ` L 100 85 L 0 85 Z`;
    return path;
  };

  // Format date based on selected period
  const formatDate = (date) => {
    if (!date) return '';
    
    switch (selectedPeriod) {
      case '1D':
        return date.toLocaleTimeString('en-US', { 
          hour: '2-digit', 
          minute: '2-digit' 
        });
      case '1W':
        return date.toLocaleDateString('en-US', { 
          weekday: 'short',
          month: 'short',
          day: 'numeric'
        });
      case '1M':
        return date.toLocaleDateString('en-US', { 
          month: 'short', 
          day: 'numeric' 
        });
      case '1Y':
        return date.toLocaleDateString('en-US', { 
          month: 'short',
          year: '2-digit'
        });
      default:
        return date.toLocaleDateString('en-US');
    }
  };

  // Loading state
  if (loading) {
    return (
      <div className="h-48 border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
          <p className="text-gray-500 text-sm">Loading chart data...</p>
        </div>
      </div>
    );
  }

  // No data state
  if (!chartData || chartData.length === 0) {
    return (
      <div className="h-48 border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <p className="text-gray-500 text-sm">No chart data available for {selectedPeriod}</p>
          <p className="text-gray-400 text-xs mt-1">Try a different time period</p>
        </div>
      </div>
    );
  }

  const lineColor = isPositiveTrend ? '#22c55e' : '#ef4444'; // green or red
  const gradientId = isPositiveTrend ? 'positiveGradient' : 'negativeGradient';

  return (
    <div className="h-48 border border-gray-200 rounded-lg bg-white relative overflow-hidden">
      {/* Period info */}
      <div className="absolute top-2 left-3 text-xs text-gray-500 z-10">
        {selectedPeriod} Performance
      </div>
      
      {/* Performance indicator */}
      {chartData.length > 1 && (
        <div className="absolute top-2 right-3 text-xs z-10">
          <span className={`font-medium ${isPositiveTrend ? 'text-green-600' : 'text-red-600'}`}>
            {isPositiveTrend ? '↗' : '↘'} {Math.abs(((chartData[chartData.length - 1].value - chartData[0].value) / chartData[0].value) * 100).toFixed(2)}%
          </span>
        </div>
      )}

      {/* SVG Chart */}
      <svg 
        viewBox="0 0 100 100" 
        className="w-full h-full p-4"
        preserveAspectRatio="none"
      >
        {/* Gradient definitions */}
        <defs>
          <linearGradient id="positiveGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#22c55e" stopOpacity="0.3"/>
            <stop offset="100%" stopColor="#22c55e" stopOpacity="0.05"/>
          </linearGradient>
          <linearGradient id="negativeGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#ef4444" stopOpacity="0.3"/>
            <stop offset="100%" stopColor="#ef4444" stopOpacity="0.05"/>
          </linearGradient>
        </defs>
        
        {/* Colored shadow area */}
        <path
          d={createAreaPath()}
          fill={`url(#${gradientId})`}
        />
        
        {/* Performance line */}
        <path
          d={createPath()}
          fill="none"
          stroke={lineColor}
          strokeWidth="1.5"
          vectorEffect="non-scaling-stroke"
        />
      </svg>

      {/* Time labels */}
      {chartData.length > 1 && (
        <div className="absolute bottom-1 left-3 right-3 flex justify-between text-xs text-gray-400">
          <span>{formatDate(chartData[0].date)}</span>
          <span>{formatDate(chartData[chartData.length - 1].date)}</span>
        </div>
      )}
    </div>
  );
};

export default PortfolioChart;