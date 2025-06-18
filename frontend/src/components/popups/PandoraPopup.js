import React, { useState } from 'react';

// Question Section Component
const QuestionSection = ({ title, options, selected, onSelect }) => {
  return (
    <div>
      <h3 className="text-base font-semibold mb-3 text-gray-800">{title}</h3>
      <div className="space-y-2">
        {options.map((option) => (
          <div
            key={option.label}
            onClick={() => onSelect(option.label)}
            className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${selected === option.label ? 'border-purple-500 bg-gradient-to-r from-purple-50 to-blue-50' : 'border-gray-200 hover:border-purple-500 hover:bg-purple-50'}`}
          >
            <div className="font-medium">{option.label}</div>
            <div className="text-xs text-gray-600">{option.desc}</div>
          </div>
        ))}
      </div>
    </div>
  );
};

const PandoraPopup = ({ onClose }) => {
  const [formData, setFormData] = useState({
    goal: '',
    timeHorizon: '',
    riskTolerance: '',
    investmentPercentage: 50,
    sectors: []
  });

  const handleOptionSelect = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSectorToggle = (sector) => {
    setFormData(prev => ({
      ...prev,
      sectors: prev.sectors.includes(sector)
        ? prev.sectors.filter(s => s !== sector)
        : prev.sectors.length < 3
        ? [...prev.sectors, sector]
        : prev.sectors
    }));
  };

  const handleSubmit = () => {
    console.log('Pandora Form Data:', formData);
    alert('Investment profile submitted! Pandora will now generate personalized recommendations based on your preferences.');
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="relative p-6 text-center border-b border-gray-100">
          <button
            onClick={onClose}
            className="absolute top-4 right-4 w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 text-gray-500 hover:text-gray-700"
          >
            ×
          </button>
          <h2 className="text-2xl font-bold bg-gradient-to-r from-purple-500 to-blue-500 bg-clip-text text-transparent mb-2">
            Pandora Investment Profile
          </h2>
          <p className="text-sm text-gray-600">
            Help us understand your investment goals and preferences to provide personalized recommendations
          </p>
        </div>

        {/* Form Content */}
        <div className="p-6 space-y-8">
          {/* Investment Goal */}
          <QuestionSection
            title="What is your primary investment goal?"
            options={[
              { label: 'Wealth Preservation', desc: 'Protect my current wealth with minimal risk' },
              { label: 'Steady Growth', desc: 'Consistent returns with moderate risk' },
              { label: 'Aggressive Growth', desc: 'Maximum returns, willing to accept higher risk' },
              { label: 'Income Generation', desc: 'Regular dividend and interest income' }
            ]}
            selected={formData.goal}
            onSelect={(value) => handleOptionSelect('goal', value)}
          />

          {/* Time Horizon */}
          <QuestionSection
            title="What is your investment time horizon?"
            options={[
              { label: 'Short-term (< 2 years)', desc: 'Need access to funds soon' },
              { label: 'Medium-term (2-10 years)', desc: 'Major purchase or life event' },
              { label: 'Long-term (10+ years)', desc: 'Retirement or distant goals' }
            ]}
            selected={formData.timeHorizon}
            onSelect={(value) => handleOptionSelect('timeHorizon', value)}
          />

          {/* Risk Tolerance */}
          <QuestionSection
            title="How would you react if your portfolio lost 20% in a month?"
            options={[
              { label: 'Sell everything immediately', desc: "I can't handle large losses" },
              { label: 'Sell some positions', desc: 'Reduce risk but stay invested' },
              { label: 'Hold my positions', desc: 'Markets recover over time' },
              { label: 'Buy more', desc: 'Great opportunity to invest at lower prices' }
            ]}
            selected={formData.riskTolerance}
            onSelect={(value) => handleOptionSelect('riskTolerance', value)}
          />

          {/* Investment Percentage */}
          <div>
            <h3 className="text-base font-semibold mb-3 text-gray-800">
              What percentage of your total assets should be invested?
            </h3>
            <input
              type="range"
              min="10"
              max="100"
              value={formData.investmentPercentage}
              onChange={(e) => handleOptionSelect('investmentPercentage', e.target.value)}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
            <div className="flex justify-between text-xs text-gray-500 mt-2">
              <span>10%</span>
              <span className="font-semibold text-purple-600">{formData.investmentPercentage}%</span>
              <span>100%</span>
            </div>
          </div>

          {/* Sectors */}
          <div>
            <h3 className="text-base font-semibold mb-3 text-gray-800">
              Which sectors interest you most? (Select up to 3)
            </h3>
            <div className="space-y-2">
              {[
                { label: 'Technology', desc: 'Software, hardware, AI, semiconductors' },
                { label: 'Healthcare', desc: 'Pharmaceuticals, biotech, medical devices' },
                { label: 'Financial Services', desc: 'Banks, insurance, fintech' },
                { label: 'Energy', desc: 'Oil, gas, renewables, utilities' },
                { label: 'Consumer Goods', desc: 'Retail, food, beverages, luxury' },
                { label: 'Real Estate', desc: 'REITs, property development' }
              ].map((sector) => (
                <div
                  key={sector.label}
                  onClick={() => handleSectorToggle(sector.label)}
                  className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${formData.sectors.includes(sector.label) ? 'border-purple-500 bg-gradient-to-r from-purple-50 to-blue-50' : 'border-gray-200 hover:border-purple-500 hover:bg-purple-50'}`}
                >
                  <div className="font-medium">{sector.label}</div>
                  <div className="text-xs text-gray-600">{sector.desc}</div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
          <button
            onClick={onClose}
            className="px-5 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            className="px-5 py-2 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-lg font-semibold hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-200"
          >
            Generate Recommendations
          </button>
        </div>
      </div>
    </div>
  );
};

export default PandoraPopup; 