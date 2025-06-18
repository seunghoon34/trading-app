// Debug script for Pandora form issues
// Run this in the browser console when on the Pandora popup page

function debugPandoraForm() {
  console.log('🔮 Debugging Pandora Form...');
  
  // Sample form data as it would appear in the form
  const sampleFormData = {
    risk_tolerance: 'Hold my positions',
    investment_timeline: 'Medium-term (2-10 years)',
    financial_goals: ['Steady Growth', 'Retirement'],
    age_bracket: '26-35',
    annual_income_bracket: '50000-75000',
    investment_experience: 'Intermediate',
    risk_capacity: 'medium'
  };

  console.log('Sample form data:', sampleFormData);

  // Test the mapping function
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

  const mappedData = {
    risk_tolerance: mapping.risk_tolerance[sampleFormData.risk_tolerance],
    investment_timeline: mapping.investment_timeline[sampleFormData.investment_timeline],
    financial_goals: sampleFormData.financial_goals.map(goal => mapping.financial_goals[goal]),
    age_bracket: sampleFormData.age_bracket,
    annual_income_bracket: sampleFormData.annual_income_bracket,
    investment_experience: mapping.investment_experience[sampleFormData.investment_experience],
    risk_capacity: sampleFormData.risk_capacity
  };

  console.log('Mapped API data:', mappedData);

  // Check if all required fields are properly mapped
  const requiredFields = [
    'risk_tolerance',
    'investment_timeline', 
    'financial_goals',
    'age_bracket',
    'annual_income_bracket',
    'investment_experience',
    'risk_capacity'
  ];

  console.log('Field validation:');
  requiredFields.forEach(field => {
    const value = mappedData[field];
    const isValid = value && (Array.isArray(value) ? value.length > 0 : true);
    console.log(`${field}: ${JSON.stringify(value)} - ${isValid ? '✅' : '❌'}`);
  });

  // Backend expected enum values
  const backendValidation = {
    risk_tolerance: ['conservative', 'moderate', 'aggressive'],
    investment_timeline: ['short_term', 'medium_term', 'long_term'],
    financial_goals: ['retirement', 'wealth_building', 'income_generation', 'capital_preservation', 'education', 'home_purchase'],
    age_bracket: ['18-25', '26-35', '36-45', '46-55', '56-65', '65+'],
    annual_income_bracket: ['0-25000', '25000-50000', '50000-75000', '75000-100000', '100000-150000', '150000+'],
    investment_experience: ['beginner', 'intermediate', 'advanced'],
    risk_capacity: ['low', 'medium', 'high']
  };

  console.log('Backend validation check:');
  Object.keys(backendValidation).forEach(field => {
    const value = mappedData[field];
    const allowedValues = backendValidation[field];
    
    if (Array.isArray(value)) {
      const allValid = value.every(v => allowedValues.includes(v));
      console.log(`${field}: ${JSON.stringify(value)} - ${allValid ? '✅' : '❌'}`);
      if (!allValid) {
        const invalid = value.filter(v => !allowedValues.includes(v));
        console.log(`  Invalid values: ${JSON.stringify(invalid)}`);
      }
    } else {
      const isValid = allowedValues.includes(value);
      console.log(`${field}: "${value}" - ${isValid ? '✅' : '❌'}`);
      if (!isValid) {
        console.log(`  Expected one of: ${JSON.stringify(allowedValues)}`);
      }
    }
  });

  return {
    formData: sampleFormData,
    mappedData: mappedData,
    validation: backendValidation
  };
}

// Test the debug function
console.log('Run debugPandoraForm() to test the form mapping');

// Export for use
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { debugPandoraForm };
} 