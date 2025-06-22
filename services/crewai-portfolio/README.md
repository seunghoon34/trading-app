# CrewAI Portfolio Manager

A modular AI-powered portfolio management system using CrewAI with Anthropic Claude 3.5 Sonnet and OpenAI fallback.

## Features

- **Portfolio Generation**: Generate new portfolios based on user risk profiles
- **Portfolio Rebalancing**: Analyze and rebalance existing portfolios
- **Modular Architecture**: Clean separation of concerns across multiple files
- **Multiple LLM Support**: Primary Anthropic Claude with OpenAI fallback
- **Real-time Market Analysis**: Uses web scraping and search tools for current market data

## Project Structure

```
services/crewai-portfolio/
├── main.py              # FastAPI application with endpoints
├── models.py            # Pydantic models for requests/responses
├── agents.py            # CrewAI agent definitions
├── tasks.py             # CrewAI task definitions
├── tools.py             # Custom tools and utilities
├── utils.py             # LLM configuration and helper functions
├── example_usage.py     # Example usage scripts
├── requirements.txt     # Python dependencies
└── README.md           # This file
```

## API Endpoints

### 1. Portfolio Generation
**POST** `/generate-portfolio`

Generate a new portfolio based on user profile.

```json
{
  "risk_tolerance": "high",
  "investment_timeline": "long_term",
  "financial_goals": ["wealth_building", "retirement"],
  "age_bracket": "26-35",
  "annual_income_bracket": "1000000-2000000",
  "investment_experience": "high",
  "risk_capacity": "high"
}
```

**Response:**
```json
{
  "portfolio": [
    {"symbol": "VTI", "weight": 0.40},
    {"symbol": "AAPL", "weight": 0.20},
    {"symbol": "QQQ", "weight": 0.25},
    {"symbol": "MSFT", "weight": 0.15}
  ],
  "explanation": "Portfolio strategy explanation...",
  "status": "success",
  "message": "Portfolio generated successfully"
}
```

### 2. Portfolio Rebalancing
**POST** `/rebalance-portfolio`

Analyze and potentially rebalance an existing portfolio.

```json
{
  "user_profile": {
    "risk_tolerance": "high",
    "investment_timeline": "long_term",
    "financial_goals": ["wealth_building"],
    "age_bracket": "26-35",
    "investment_experience": "high",
    "risk_capacity": "high"
  },
  "current_portfolio": [
    {"symbol": "VTI", "weight": 0.30, "performance": "+12%"},
    {"symbol": "AAPL", "weight": 0.35, "performance": "+28%"},
    {"symbol": "QQQ", "weight": 0.20, "performance": "+15%"},
    {"symbol": "MSFT", "weight": 0.15, "performance": "+22%"}
  ],
  "time_since_last_rebalance": "6_months"
}
```

**Response:**
```json
{
  "portfolio": [
    {"symbol": "VTI", "weight": 0.40},
    {"symbol": "AAPL", "weight": 0.25},
    {"symbol": "QQQ", "weight": 0.20},
    {"symbol": "MSFT", "weight": 0.15}
  ],
  "explanation": "Rebalancing rationale...",
  "status": "success",
  "message": "Portfolio rebalancing analysis completed",
  "rebalance_needed": true,
  "changes_made": [
    "Reduced AAPL from 35% to 25% due to concentration risk",
    "Increased VTI allocation for better diversification"
  ]
}
```

## Setup and Installation

1. **Environment Variables**
   Create a `.env` file with:
   ```env
   ANTHROPIC_API_KEY=your_anthropic_key_here
   OPENAI_API_KEY=your_openai_key_here  # fallback
   SERPER_API_KEY=your_serper_key_here
   ```

2. **Install Dependencies**
   ```bash
   pip install -r requirements.txt
   ```

3. **Run the API**
   ```bash
   python main.py
   ```
   The API will start on http://localhost:8000

4. **Test the API**
   ```bash
   python example_usage.py
   ```

## Architecture Components

### Agents (agents.py)

**Portfolio Generation Agents:**
- Economic Research Agent: Analyzes macroeconomic trends
- Market Research Agent: Identifies investment opportunities
- Fundamental Analyst: Evaluates securities fundamentals
- Risk Analyst: Assesses portfolio risk metrics
- Portfolio Manager: Constructs optimal allocations
- Validation Agent: Validates tickers and formatting

**Portfolio Rebalancing Agents:**
- Portfolio Performance Analyst: Analyzes current performance
- Market Conditions Analyst: Assesses market environment changes
- Risk Drift Analyst: Evaluates risk characteristic changes
- Rebalance Strategist: Determines rebalancing strategy
- Rebalance Validation Agent: Validates recommendations

### Tasks (tasks.py)

Tasks are created dynamically based on user profiles and current portfolios. Each task has detailed instructions for the agents and expected outputs.

### Models (models.py)

- `UserProfile`: User risk and investment preferences
- `PortfolioHolding`: Individual portfolio position with performance
- `RebalanceRequest`: Complete rebalancing request structure
- `PortfolioResponse`: Portfolio generation response
- `RebalanceResponse`: Portfolio rebalancing response

### Tools (tools.py)

- Web scraping tool for market data
- Search tool for research
- Custom stock purchase tool for tracking

### Utils (utils.py)

- LLM configuration with Anthropic/OpenAI fallback
- JSON parsing utilities
- Portfolio validation functions

## Usage Examples

### Generate Portfolio
```python
import requests

user_profile = {
    "risk_tolerance": "medium",
    "investment_timeline": "medium_term",
    "financial_goals": ["wealth_building"],
    "investment_experience": "medium"
}

response = requests.post("http://localhost:8000/generate-portfolio", json=user_profile)
portfolio = response.json()
```

### Rebalance Portfolio
```python
rebalance_request = {
    "user_profile": user_profile,
    "current_portfolio": [
        {"symbol": "VTI", "weight": 0.40, "performance": "+10%"},
        {"symbol": "AAPL", "weight": 0.35, "performance": "+25%"},
        {"symbol": "MSFT", "weight": 0.25, "performance": "+15%"}
    ],
    "time_since_last_rebalance": "3_months"
}

response = requests.post("http://localhost:8000/rebalance-portfolio", json=rebalance_request)
result = response.json()
```

## Risk Profile Options

- **risk_tolerance**: "low", "medium", "high"
- **investment_timeline**: "short_term", "medium_term", "long_term"
- **investment_experience**: "low", "medium", "high"
- **risk_capacity**: "low", "medium", "high"

## Contributing

1. Follow the modular structure when adding new features
2. Add new agents to `agents.py` with proper role definitions
3. Create corresponding tasks in `tasks.py`
4. Update models in `models.py` for new data structures
5. Add utility functions to `utils.py` as needed

## License

This project is part of the trading-app ecosystem. 