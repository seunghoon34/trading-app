import React from 'react';
import PortfolioPerformance from './PortfolioPerformance';
import RecentOrders from './RecentOrders';
import Positions from './Positions';
import MarketData from './MarketData';

const Dashboard = () => {
  return (
    <div className="h-full flex flex-col p-6 gap-4">
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-3 gap-4 min-h-0">
        {/* Main Section */}
        <div className="lg:col-span-2 flex flex-col gap-4 min-h-0">
          <div className="flex-1 flex flex-col min-h-0">
            <PortfolioPerformance />
            <div className="flex-1 min-h-0 flex flex-col">
              <RecentOrders />
            </div>
          </div>
        </div>
        {/* Sidebar Section */}
        <div className="flex flex-col gap-4 min-h-0 flex-1">
          <Positions />
          <MarketData />
        </div>
      </div>
    </div>
  );
};

export default Dashboard; 