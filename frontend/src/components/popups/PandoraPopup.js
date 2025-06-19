import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';

// API service functions
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

const apiService = {
  // Risk Profile APIs
  getRiskProfile: async () => {
    const response = await fetch(`${API_BASE_URL}/api/v1/investment-strategy/risk-profile`, {
      method: 'GET',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (response.status === 404) {
      return null; // No risk profile found
    }
    
    if (!response.ok) {
      throw new Error('Failed to fetch risk profile');
    }
    
    return await response.json();
  },

  createRiskProfile: async (riskProfileData) => {
    console.log('Sending risk profile data:', riskProfileData);
    const response = await fetch(`${API_BASE_URL}/api/v1/investment-strategy/risk-profile`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(riskProfileData),
    });
    
    if (!response.ok) {
      const error = await response.json();
      console.error('Risk profile creation error:', error);
      throw new Error(error.error || 'Failed to create risk profile');
    }
    
    return await response.json();
  },

  updateRiskProfile: async (riskProfileData) => {
    const response = await fetch(`${API_BASE_URL}/api/v1/investment-strategy/risk-profile`, {
      method: 'PUT',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(riskProfileData),
    });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to update risk profile');
    }
    
    return await response.json();
  },

  // Portfolio Generation API
  generatePortfolio: async (userProfile) => {
    const response = await fetch(`${API_BASE_URL}/api/v1/crewai-portfolio/generate-portfolio`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userProfile),
    });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.detail || 'Failed to generate portfolio');
    }
    
    return await response.json();
  },

  // Portfolio Management APIs
  createOrUpdatePortfolio: async (portfolioData) => {
    // Try to create first, if it exists, update it
    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/investment-strategy/portfolio`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(portfolioData),
      });
      
      if (response.status === 409) {
        // Portfolio exists, update it
        const updateResponse = await fetch(`${API_BASE_URL}/api/v1/investment-strategy/portfolio`, {
          method: 'PUT',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(portfolioData),
        });
        
        if (!updateResponse.ok) {
          const error = await updateResponse.json();
          throw new Error(error.error || 'Failed to update portfolio');
        }
        
        return await updateResponse.json();
      }
      
      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to create portfolio');
      }
      
      return await response.json();
    } catch (error) {
      throw error;
    }
  },

  // Purchase Portfolio API
  purchasePortfolio: async () => {
    const response = await fetch(`${API_BASE_URL}/api/v1/investment-strategy/portfolio/purchase`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to purchase portfolio');
    }
    
    return await response.json();
  },
};

// Loading Screen Component for Portfolio Generation
const LoadingScreen = () => {
  return (
    <div className="space-y-8 text-center py-12">
      {/* Animated Icon */}
      <div className="relative mx-auto w-24 h-24">
        <div className="absolute inset-0 rounded-full border-4 border-purple-200"></div>
        <div className="absolute inset-0 rounded-full border-4 border-transparent border-t-purple-500 animate-spin"></div>
        <div className="absolute inset-2 rounded-full border-4 border-transparent border-t-blue-500 animate-spin" style={{ animationDirection: 'reverse' }}></div>
        <div className="absolute inset-0 flex items-center justify-center">
          <svg className="w-8 h-8 text-purple-600" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M4 4a2 2 0 00-2 2v4a2 2 0 002 2V6h10a2 2 0 00-2-2H4zm2 6a2 2 0 012-2h8a2 2 0 012 2v4a2 2 0 01-2 2H8a2 2 0 01-2-2v-4zm6 4a2 2 0 100-4 2 2 0 000 4z" clipRule="evenodd" />
          </svg>
        </div>
      </div>

      {/* Loading Text */}
      <div className="space-y-2">
        <h3 className="text-xl font-bold text-gray-800">Generating Your Portfolio</h3>
        <p className="text-gray-600">Our AI is analyzing your preferences and market data...</p>
      </div>

      {/* Progress Animation */}
      <div className="max-w-md mx-auto">
        <div className="bg-gray-200 rounded-full h-2 overflow-hidden">
          <div className="bg-gradient-to-r from-purple-500 to-blue-500 h-full rounded-full animate-pulse"></div>
        </div>
      </div>

      {/* Fun Facts */}
      <div className="bg-purple-50 rounded-lg p-4 max-w-md mx-auto">
        <p className="text-sm text-purple-800">
          💡 <strong>Did you know?</strong> Our AI considers over 100 market factors to create your personalized portfolio
        </p>
      </div>

      {/* Loading Message */}
      <div className="text-center">
        <p className="text-gray-500 text-sm">
          Please wait while we create your personalized portfolio...
        </p>
      </div>
    </div>
  );
};

// Question Section Component
const QuestionSection = ({ title, options, selected, onSelect, description }) => {
  return (
    <div>
      <h3 className="text-base font-semibold mb-2 text-gray-800">{title}</h3>
      {description && <p className="text-sm text-gray-600 mb-3">{description}</p>}
      <div className="space-y-2">
        {options.map((option) => (
          <div
            key={option.label}
            onClick={() => onSelect(option.label)}
            className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${
              selected === option.label 
                ? 'border-purple-500 bg-gradient-to-r from-purple-50 to-blue-50' 
                : 'border-gray-200 hover:border-purple-500 hover:bg-purple-50'
            }`}
          >
            <div className="font-medium">{option.label}</div>
            <div className="text-xs text-gray-600">{option.desc}</div>
          </div>
        ))}
      </div>
    </div>
  );
};

// Multi-select Question Component
const MultiSelectSection = ({ title, options, selected, onToggle, maxSelections, description }) => {
  return (
    <div>
      <h3 className="text-base font-semibold mb-2 text-gray-800">{title}</h3>
      {description && <p className="text-sm text-gray-600 mb-3">{description}</p>}
      <div className="space-y-2">
        {options.map((option) => (
          <div
            key={option.label}
            onClick={() => onToggle(option.label)}
            className={`p-3 border-2 rounded-lg cursor-pointer transition-all ${
              selected.includes(option.label) 
                ? 'border-purple-500 bg-gradient-to-r from-purple-50 to-blue-50' 
                : selected.length >= maxSelections 
                  ? 'border-gray-200 bg-gray-50 cursor-not-allowed' 
                  : 'border-gray-200 hover:border-purple-500 hover:bg-purple-50'
            }`}
          >
            <div className="font-medium">{option.label}</div>
            <div className="text-xs text-gray-600">{option.desc}</div>
          </div>
        ))}
      </div>
      <p className="text-xs text-gray-500 mt-2">
        {selected.length} of {maxSelections} selected
      </p>
    </div>
  );
};

// Portfolio Display Component
const PortfolioDisplay = ({ portfolio, explanation, onRegenerate, onContinue, loading }) => {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-xl font-bold text-gray-800 mb-4">Generated Portfolio</h3>
        
        {/* Portfolio Positions */}
        <div className="bg-gray-50 rounded-lg p-4 mb-4">
          <h4 className="font-semibold mb-3">Portfolio Allocation</h4>
          <div className="space-y-2">
            {portfolio.map((position, index) => (
              <div key={index} className="flex justify-between items-center py-2 px-3 bg-white rounded">
                <span className="font-medium">{position.symbol}</span>
                <span className="text-purple-600 font-semibold">
                  {(position.weight * 100).toFixed(1)}%
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Explanation */}
        <div className="bg-blue-50 rounded-lg p-4">
          <h4 className="font-semibold mb-2">Investment Strategy</h4>
          <p className="text-sm text-gray-700 leading-relaxed">{explanation}</p>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="flex gap-3">
        <button
          onClick={onRegenerate}
          disabled={loading}
          className="flex-1 px-4 py-2 border border-purple-500 text-purple-500 rounded-lg font-semibold hover:bg-purple-50 transition-colors disabled:opacity-50"
        >
          {loading ? 'Generating...' : 'Regenerate Portfolio'}
        </button>
        <button
          onClick={onContinue}
          disabled={loading}
          className="flex-1 px-4 py-2 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-lg font-semibold hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-200 disabled:opacity-50"
        >
          Continue to Purchase
        </button>
      </div>
    </div>
  );
};

// Purchase Result Component
const PurchaseResult = ({ result, onClose, onPurchaseComplete }) => {
  return (
    <div className="space-y-6">
      <div className="text-center">
        <h3 className="text-xl font-bold text-gray-800 mb-2">Purchase Complete!</h3>
        <p className="text-gray-600">Your portfolio has been purchased successfully.</p>
      </div>

      <div className="bg-green-50 rounded-lg p-4">
        <h4 className="font-semibold mb-3 text-green-800">Purchase Summary</h4>
        <div className="space-y-2 text-sm">
          <div className="flex justify-between">
            <span>Total Buying Power:</span>
            <span className="font-semibold">${result.result.total_buying_power}</span>
          </div>
          <div className="flex justify-between">
            <span>Successful Orders:</span>
            <span className="font-semibold">{result.result.success_count}</span>
          </div>
          {result.result.failure_count > 0 && (
            <div className="flex justify-between">
              <span>Failed Orders:</span>
              <span className="font-semibold text-red-600">{result.result.failure_count}</span>
            </div>
          )}
        </div>
      </div>

      {/* Order Details */}
      <div className="space-y-2">
        <h4 className="font-semibold">Order Details</h4>
        {result.result.order_results.map((order, index) => (
          <div
            key={index}
            className={`p-3 rounded-lg ${
              order.success ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'
            }`}
          >
            <div className="flex justify-between items-center">
              <span className="font-medium">{order.symbol}</span>
              <span className="text-sm">${order.notional}</span>
            </div>
            {order.error && (
              <p className="text-xs text-red-600 mt-1">{order.error}</p>
            )}
            {order.order_id && (
              <p className="text-xs text-gray-500 mt-1">Order ID: {order.order_id}</p>
            )}
          </div>
        ))}
      </div>

      <button
        onClick={() => {
          onPurchaseComplete && onPurchaseComplete();
          onClose();
        }}
        className="w-full px-4 py-2 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-lg font-semibold hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-200"
      >
        Close
      </button>
    </div>
  );
};

const PandoraPopup = ({ onClose, onPurchaseComplete }) => {
  const { user } = useAuth();
  const [currentStep, setCurrentStep] = useState('loading'); // loading, riskProfile, editProfile, portfolioGenerating, portfolio, purchase, complete
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [existingRiskProfile, setExistingRiskProfile] = useState(null);
  const [generatedPortfolio, setGeneratedPortfolio] = useState(null);
  const [purchaseResult, setPurchaseResult] = useState(null);

  // Form state
  const [formData, setFormData] = useState({
    risk_tolerance: '',
    investment_timeline: '',
    financial_goals: [],
    age_bracket: '',
    annual_income_bracket: '',
    investment_experience: '',
    risk_capacity: ''
  });

  useEffect(() => {
    checkExistingRiskProfile();
  }, []);

  const mapApiToFormData = (apiData) => {
    // Reverse mapping from API values to frontend display values
    const reverseMapping = {
      risk_tolerance: {
        'conservative': 'Sell some positions',
        'moderate': 'Hold my positions',
        'aggressive': 'Buy more'
      },
      investment_timeline: {
        'short_term': 'Short-term (< 2 years)',
        'medium_term': 'Medium-term (2-10 years)',
        'long_term': 'Long-term (10+ years)'
      },
      financial_goals: {
        'capital_preservation': 'Wealth Preservation',
        'wealth_building': 'Steady Growth',
        'income_generation': 'Income Generation',
        'retirement': 'Retirement',
        'education': 'Education',
        'home_purchase': 'Home Purchase'
      },
      investment_experience: {
        'beginner': 'Beginner',
        'intermediate': 'Intermediate',
        'advanced': 'Advanced'
      }
    };

    return {
      risk_tolerance: reverseMapping.risk_tolerance[apiData.risk_tolerance] || '',
      investment_timeline: reverseMapping.investment_timeline[apiData.investment_timeline] || '',
      financial_goals: apiData.financial_goals?.map(goal => reverseMapping.financial_goals[goal]).filter(Boolean) || [],
      age_bracket: apiData.age_bracket || '',
      annual_income_bracket: apiData.annual_income_bracket || '',
      investment_experience: reverseMapping.investment_experience[apiData.investment_experience] || '',
      risk_capacity: apiData.risk_capacity || ''
    };
  };

  const checkExistingRiskProfile = async () => {
    try {
      setLoading(true);
      const riskProfile = await apiService.getRiskProfile();
      
      if (riskProfile) {
        setExistingRiskProfile(riskProfile);
        // Convert API data to form display values
        const formDisplayData = mapApiToFormData(riskProfile);
        setFormData(formDisplayData);
        setCurrentStep('riskProfile');
      } else {
        setCurrentStep('riskProfile');
      }
    } catch (error) {
      console.error('Error checking risk profile:', error);
      setError('Failed to load risk profile. Please try again.');
      setCurrentStep('riskProfile');
    } finally {
      setLoading(false);
    }
  };

  const handleOptionSelect = (field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleFinancialGoalsToggle = (goal) => {
    setFormData(prev => ({
      ...prev,
      financial_goals: prev.financial_goals.includes(goal)
        ? prev.financial_goals.filter(g => g !== goal)
        : prev.financial_goals.length < 3
        ? [...prev.financial_goals, goal]
        : prev.financial_goals
    }));
  };

  const mapFormDataToApiFormat = (data) => {
    // Map frontend form values to API expected values
    const mapping = {
      risk_tolerance: {
        'Sell everything immediately': 'conservative',
        'Sell some positions': 'conservative',
        'Hold my positions': 'moderate',
        'Buy more': 'aggressive'
      },
      investment_timeline: {
        'Short-term (< 2 years)': 'short_term',
        'Medium-term (2-10 years)': 'medium_term',
        'Long-term (10+ years)': 'long_term'
      },
      financial_goals: {
        'Wealth Preservation': 'capital_preservation',
        'Steady Growth': 'wealth_building',
        'Aggressive Growth': 'wealth_building',
        'Income Generation': 'income_generation',
        'Retirement': 'retirement',
        'Education': 'education',
        'Home Purchase': 'home_purchase'
      },
      investment_experience: {
        'Beginner': 'beginner',
        'Intermediate': 'intermediate',
        'Advanced': 'advanced'
      }
    };

    // Validate that all mappings are successful
    const mappedData = {
      risk_tolerance: mapping.risk_tolerance[data.risk_tolerance],
      investment_timeline: mapping.investment_timeline[data.investment_timeline],
      financial_goals: data.financial_goals.map(goal => mapping.financial_goals[goal]).filter(Boolean),
      age_bracket: data.age_bracket,
      annual_income_bracket: data.annual_income_bracket,
      investment_experience: mapping.investment_experience[data.investment_experience],
      risk_capacity: data.risk_capacity
    };

    // Check for missing mappings
    if (!mappedData.risk_tolerance) {
      throw new Error(`Invalid risk tolerance value: ${data.risk_tolerance}`);
    }
    if (!mappedData.investment_timeline) {
      throw new Error(`Invalid investment timeline value: ${data.investment_timeline}`);
    }
    if (!mappedData.investment_experience) {
      throw new Error(`Invalid investment experience value: ${data.investment_experience}`);
    }
    if (mappedData.financial_goals.length !== data.financial_goals.length) {
      throw new Error(`Some financial goals could not be mapped: ${data.financial_goals}`);
    }

    return mappedData;
  };

  const handleContinue = async () => {
    try {
      setLoading(true);
      setError('');

      // Validate form
      if (!formData.risk_tolerance || !formData.investment_timeline || 
          formData.financial_goals.length === 0 || !formData.age_bracket || 
          !formData.annual_income_bracket || !formData.investment_experience || 
          !formData.risk_capacity) {
        setError('Please fill in all fields');
        return;
      }

      // Map form data to API format
      const apiFormData = mapFormDataToApiFormat(formData);
      
      // Debug: Log the data being sent
      console.log('Form Data (before mapping):', formData);
      console.log('API Data (after mapping):', apiFormData);

      // Create or update risk profile
      if (existingRiskProfile) {
        await apiService.updateRiskProfile(apiFormData);
      } else {
        await apiService.createRiskProfile(apiFormData);
      }

      // Switch to portfolio generating screen
      setCurrentStep('portfolioGenerating');
      setLoading(false); // Reset loading for the button

      // Generate portfolio
      const portfolioResponse = await apiService.generatePortfolio(apiFormData);
      setGeneratedPortfolio(portfolioResponse);
      setCurrentStep('portfolio');

    } catch (error) {
      console.error('Error:', error);
      setError(error.message || 'An error occurred. Please try again.');
      setCurrentStep('riskProfile'); // Go back to form on error
    } finally {
      setLoading(false);
    }
  };

  const handleRegenerate = async () => {
    try {
      setLoading(true);
      setError('');

      // Switch to portfolio generating screen
      setCurrentStep('portfolioGenerating');
      setLoading(false); // Reset loading for the loading screen

      const apiFormData = mapFormDataToApiFormat(formData);
      const portfolioResponse = await apiService.generatePortfolio(apiFormData);
      setGeneratedPortfolio(portfolioResponse);
      setCurrentStep('portfolio');

    } catch (error) {
      console.error('Error regenerating portfolio:', error);
      setError(error.message || 'Failed to regenerate portfolio. Please try again.');
      setCurrentStep('portfolio'); // Go back to portfolio on error
    } finally {
      setLoading(false);
    }
  };

  const handlePortfolioContinue = async () => {
    try {
      setLoading(true);
      setError('');

      // Store portfolio
      const portfolioData = {
        positions: generatedPortfolio.portfolio.map(pos => ({
          symbol: pos.symbol,
          weight: pos.weight
        }))
      };

      await apiService.createOrUpdatePortfolio(portfolioData);

      // Purchase portfolio
      const purchaseResponse = await apiService.purchasePortfolio();
      setPurchaseResult(purchaseResponse);
      setCurrentStep('complete');

    } catch (error) {
      console.error('Error purchasing portfolio:', error);
      setError(error.message || 'Failed to purchase portfolio. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const renderContent = () => {
    if (currentStep === 'loading') {
      return (
        <div className="text-center py-8">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading your profile...</p>
        </div>
      );
    }

    if (currentStep === 'portfolioGenerating') {
      return <LoadingScreen />;
    }

    if (currentStep === 'complete') {
      return <PurchaseResult result={purchaseResult} onClose={onClose} onPurchaseComplete={onPurchaseComplete} />;
    }

    if (currentStep === 'portfolio') {
      return (
        <PortfolioDisplay
          portfolio={generatedPortfolio.portfolio}
          explanation={generatedPortfolio.explanation}
          onRegenerate={handleRegenerate}
          onContinue={handlePortfolioContinue}
          loading={loading}
        />
      );
    }

    // Risk Profile form
    return (
      <div className="space-y-6">
        {existingRiskProfile && (
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p className="text-blue-800 text-sm">
              <strong>Existing Profile Found:</strong> You can edit your settings below or continue with your current profile.
            </p>
          </div>
        )}

        {/* Risk Tolerance */}
        <QuestionSection
          title="How would you react if your portfolio lost 20% in a month?"
          options={[
            { label: 'Sell everything immediately', desc: "I can't handle large losses" },
            { label: 'Sell some positions', desc: 'Reduce risk but stay invested' },
            { label: 'Hold my positions', desc: 'Markets recover over time' },
            { label: 'Buy more', desc: 'Great opportunity to invest at lower prices' }
          ]}
          selected={formData.risk_tolerance}
          onSelect={(value) => handleOptionSelect('risk_tolerance', value)}
        />

        {/* Investment Timeline */}
        <QuestionSection
          title="What is your investment time horizon?"
          options={[
            { label: 'Short-term (< 2 years)', desc: 'Need access to funds soon' },
            { label: 'Medium-term (2-10 years)', desc: 'Major purchase or life event' },
            { label: 'Long-term (10+ years)', desc: 'Retirement or distant goals' }
          ]}
          selected={formData.investment_timeline}
          onSelect={(value) => handleOptionSelect('investment_timeline', value)}
        />

        {/* Financial Goals */}
        <MultiSelectSection
          title="What are your primary financial goals?"
          description="Select up to 3 goals"
          options={[
            { label: 'Wealth Preservation', desc: 'Protect my current wealth with minimal risk' },
            { label: 'Steady Growth', desc: 'Consistent returns with moderate risk' },
            { label: 'Aggressive Growth', desc: 'Maximum returns, willing to accept higher risk' },
            { label: 'Income Generation', desc: 'Regular dividend and interest income' },
            { label: 'Retirement', desc: 'Building wealth for retirement' },
            { label: 'Education', desc: 'Saving for education expenses' },
            { label: 'Home Purchase', desc: 'Saving for a home down payment' }
          ]}
          selected={formData.financial_goals}
          onToggle={handleFinancialGoalsToggle}
          maxSelections={3}
        />

        {/* Age Bracket */}
        <QuestionSection
          title="What is your age bracket?"
          options={[
            { label: '18-25', desc: 'Young professional' },
            { label: '26-35', desc: 'Early career' },
            { label: '36-45', desc: 'Mid career' },
            { label: '46-55', desc: 'Pre-retirement' },
            { label: '56-65', desc: 'Near retirement' },
            { label: '65+', desc: 'Retirement age' }
          ]}
          selected={formData.age_bracket}
          onSelect={(value) => handleOptionSelect('age_bracket', value)}
        />

        {/* Annual Income Bracket */}
        <QuestionSection
          title="What is your annual income bracket?"
          options={[
            { label: '0-25000', desc: 'Entry level income' },
            { label: '25000-50000', desc: 'Lower middle income' },
            { label: '50000-75000', desc: 'Middle income' },
            { label: '75000-100000', desc: 'Upper middle income' },
            { label: '100000-150000', desc: 'High income' },
            { label: '150000+', desc: 'Very high income' }
          ]}
          selected={formData.annual_income_bracket}
          onSelect={(value) => handleOptionSelect('annual_income_bracket', value)}
        />

        {/* Investment Experience */}
        <QuestionSection
          title="What is your investment experience level?"
          options={[
            { label: 'Beginner', desc: 'New to investing' },
            { label: 'Intermediate', desc: 'Some investment experience' },
            { label: 'Advanced', desc: 'Experienced investor' }
          ]}
          selected={formData.investment_experience}
          onSelect={(value) => handleOptionSelect('investment_experience', value)}
        />

        {/* Risk Capacity */}
        <QuestionSection
          title="How much risk can you afford to take?"
          options={[
            { label: 'low', desc: 'I need my money to be safe' },
            { label: 'medium', desc: 'I can handle some volatility' },
            { label: 'high', desc: 'I can afford significant losses' }
          ]}
          selected={formData.risk_capacity}
          onSelect={(value) => handleOptionSelect('risk_capacity', value)}
        />
      </div>
    );
  };

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      onClick={currentStep === 'portfolioGenerating' ? undefined : (e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div className="bg-white rounded-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="relative p-6 text-center border-b border-gray-100">
          <button
            onClick={currentStep === 'portfolioGenerating' ? undefined : onClose}
            disabled={currentStep === 'portfolioGenerating'}
            className={`absolute top-4 right-4 w-8 h-8 flex items-center justify-center rounded-full transition-colors ${
              currentStep === 'portfolioGenerating' 
                ? 'text-gray-300 cursor-not-allowed' 
                : 'hover:bg-gray-100 text-gray-500 hover:text-gray-700 cursor-pointer'
            }`}
          >
            ×
          </button>
          <h2 className="text-2xl font-bold bg-gradient-to-r from-purple-500 to-blue-500 bg-clip-text text-transparent mb-2">
            Pandora Investment AI
          </h2>
          <p className="text-sm text-gray-600">
            {currentStep === 'portfolio' 
              ? 'Your AI-generated investment portfolio' 
              : currentStep === 'portfolioGenerating'
              ? 'AI is creating your personalized portfolio'
              : currentStep === 'complete'
              ? 'Portfolio purchase completed'
              : existingRiskProfile 
              ? 'Review and update your investment profile'
              : 'Create your personalized investment profile'
            }
          </p>
        </div>

        {/* Error Display */}
        {error && (
          <div className="mx-6 mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-700 text-sm">{error}</p>
          </div>
        )}

        {/* Content */}
        <div className="p-6">
          {renderContent()}
        </div>

        {/* Footer */}
        {currentStep === 'riskProfile' && (
        <div className="px-6 py-4 border-t border-gray-100 flex gap-3 justify-end">
          <button
            onClick={onClose}
            className="px-5 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
              onClick={handleContinue}
              disabled={loading}
              className="px-5 py-2 bg-gradient-to-r from-purple-500 to-blue-500 text-white rounded-lg font-semibold hover:transform hover:-translate-y-0.5 hover:shadow-lg transition-all duration-200 disabled:opacity-50"
          >
              {loading ? 'Saving Profile...' : existingRiskProfile ? 'Continue' : 'Create Profile & Generate Portfolio'}
          </button>
        </div>
        )}
      </div>
    </div>
  );
};

export default PandoraPopup; 