"""
Example usage of the CrewAI Portfolio Manager API
This demonstrates both portfolio generation and rebalancing functionality
"""

import requests
import json

# API base URL
BASE_URL = "http://localhost:8000"

def generate_portfolio_example():
    """Example of generating a new portfolio"""
    print("=== Portfolio Generation Example ===")
    
    # User profile for portfolio generation
    user_profile = {
        "risk_tolerance": "high",
        "investment_timeline": "long_term",
        "financial_goals": ["wealth_building", "retirement"],
        "age_bracket": "26-35",
        "annual_income_bracket": "1000000-2000000",
        "investment_experience": "high",
        "risk_capacity": "high"
    }
    
    try:
        response = requests.post(f"{BASE_URL}/generate-portfolio", json=user_profile)
        if response.status_code == 200:
            result = response.json()
            print("Portfolio Generated Successfully!")
            print(f"Status: {result['status']}")
            print(f"Message: {result['message']}")
            print("\nPortfolio:")
            for holding in result['portfolio']:
                print(f"  {holding['symbol']}: {holding['weight']*100:.1f}%")
            print(f"\nExplanation: {result['explanation']}")
            return result['portfolio']  # Return for rebalancing example
        else:
            print(f"Error: {response.status_code} - {response.text}")
            return None
    except Exception as e:
        print(f"Request failed: {e}")
        return None

def rebalance_portfolio_example(current_portfolio=None):
    """Example of rebalancing an existing portfolio"""
    print("\n=== Portfolio Rebalancing Example ===")
    
    # Use generated portfolio or sample portfolio
    if current_portfolio is None:
        current_portfolio = [
            {"symbol": "VTI", "weight": 0.40, "performance": "+12%"},
            {"symbol": "AAPL", "weight": 0.25, "performance": "+28%"},  # Grew significantly
            {"symbol": "QQQ", "weight": 0.20, "performance": "+15%"},
            {"symbol": "MSFT", "weight": 0.15, "performance": "+22%"}
        ]
    else:
        # Add sample performance data to the generated portfolio
        performances = ["+12%", "+28%", "+15%", "+22%", "-5%", "+8%", "+18%", "+3%"]
        for i, holding in enumerate(current_portfolio):
            holding["performance"] = performances[i % len(performances)]
    
    # Rebalancing request
    rebalance_request = {
        "user_profile": {
            "risk_tolerance": "high",
            "investment_timeline": "long_term",
            "financial_goals": ["wealth_building", "retirement"],
            "age_bracket": "26-35",
            "annual_income_bracket": "1000000-2000000",
            "investment_experience": "high",
            "risk_capacity": "high"
        },
        "current_portfolio": current_portfolio,
        "time_since_last_rebalance": "6_months"
    }
    
    try:
        response = requests.post(f"{BASE_URL}/rebalance-portfolio", json=rebalance_request)
        if response.status_code == 200:
            result = response.json()
            print("Portfolio Rebalancing Analysis Completed!")
            print(f"Status: {result['status']}")
            print(f"Message: {result['message']}")
            print(f"Rebalance Needed: {result['rebalance_needed']}")
            
            if result['rebalance_needed']:
                print("\nRecommended New Portfolio:")
                for holding in result['portfolio']:
                    print(f"  {holding['symbol']}: {holding['weight']*100:.1f}%")
                print(f"\nChanges Made:")
                for change in result['changes_made']:
                    print(f"  - {change}")
            else:
                print("\nNo rebalancing needed - portfolio remains optimal")
                
            print(f"\nExplanation: {result['explanation']}")
        else:
            print(f"Error: {response.status_code} - {response.text}")
    except Exception as e:
        print(f"Request failed: {e}")

def health_check():
    """Check if the API is running"""
    try:
        response = requests.get(f"{BASE_URL}/health")
        if response.status_code == 200:
            print("✓ API is healthy and running")
            return True
        else:
            print(f"✗ API health check failed: {response.status_code}")
            return False
    except Exception as e:
        print(f"✗ Cannot connect to API: {e}")
        return False

if __name__ == "__main__":
    print("CrewAI Portfolio Manager - Example Usage")
    print("=" * 50)
    
    # Check if API is running
    if not health_check():
        print("\nPlease make sure the API is running:")
        print("cd services/crewai-portfolio")
        print("python main.py")
        exit(1)
    
    print()
    
    # Generate a new portfolio
    portfolio = generate_portfolio_example()
    
    # Rebalance the portfolio (or use sample data)
    rebalance_portfolio_example(portfolio)
    
    print("\n" + "=" * 50)
    print("Example completed!") 