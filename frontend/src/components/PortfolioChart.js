import React from 'react';

const PortfolioChart = () => {
  // Sample data for the portfolio performance
  const data = [
    { date: 'Jan', value: 100000 },
    { date: 'Feb', value: 105000 },
    { date: 'Mar', value: 98000 },
    { date: 'Apr', value: 112000 },
    { date: 'May', value: 108000 },
    { date: 'Jun', value: 124567 },
  ];

  // Calculate min and max for scaling
  const values = data.map(d => d.value);
  const minValue = Math.min(...values);
  const maxValue = Math.max(...values);
  const range = maxValue - minValue;

  // Create SVG path for the line
  const createPath = () => {
    const width = 100; // percentage
    const height = 70; // Leave some margin
    const stepX = width / (data.length - 1);
    
    return data.map((point, index) => {
      const x = index * stepX;
      const y = 15 + (70 - ((point.value - minValue) / range) * height); // 15% margin on top
      return `${index === 0 ? 'M' : 'L'} ${x} ${y}`;
    }).join(' ');
  };

  // Create area path for green shadow
  const createAreaPath = () => {
    const width = 100;
    const height = 70;
    const stepX = width / (data.length - 1);
    
    let path = data.map((point, index) => {
      const x = index * stepX;
      const y = 15 + (70 - ((point.value - minValue) / range) * height);
      return `${index === 0 ? 'M' : 'L'} ${x} ${y}`;
    }).join(' ');
    
    // Close the area to the bottom
    path += ` L 100 85 L 0 85 Z`;
    return path;
  };

  return (
    <div className="h-48 border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center bg-white relative">
      {/* SVG Chart */}
      <svg 
        viewBox="0 0 100 100" 
        className="w-full h-full p-4"
        preserveAspectRatio="none"
      >
        {/* Green shadow area */}
        <path
          d={createAreaPath()}
          fill="url(#greenGradient)"
          fillOpacity="0.2"
        />
        
        {/* Gradient definition for green shadow */}
        <defs>
          <linearGradient id="greenGradient" x1="0%" y1="0%" x2="0%" y2="100%">
            <stop offset="0%" stopColor="#22c55e" stopOpacity="0.4"/>
            <stop offset="100%" stopColor="#22c55e" stopOpacity="0.05"/>
          </linearGradient>
        </defs>
        
        {/* Simple green line */}
        <path
          d={createPath()}
          fill="none"
          stroke="#22c55e"
          strokeWidth="1"
          vectorEffect="non-scaling-stroke"
        />
      </svg>
    </div>
  );
};

export default PortfolioChart;