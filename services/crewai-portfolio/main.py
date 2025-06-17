from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Dict, Any
import os
import warnings
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Warning control
warnings.filterwarnings('ignore')

# Import CrewAI components
from crewai import Agent, Task, Crew, Process
from crewai_tools import BaseTool, ScrapeWebsiteTool, SerperDevTool

# Import both Anthropic and OpenAI for fallback
try:
    from langchain_anthropic import ChatAnthropic
    ANTHROPIC_AVAILABLE = True
except ImportError:
    ANTHROPIC_AVAILABLE = False
    print("Warning: langchain_anthropic not installed. Install with: pip install langchain-anthropic")

from langchain_openai import ChatOpenAI

app = FastAPI(title="CrewAI Portfolio Manager", version="1.0.0")

# Pydantic models for request/response
class UserProfile(BaseModel):
    risk_tolerance: str = "high"  # low, medium, high
    investment_timeline: str = "short_term"  # short_term, medium_term, long_term
    financial_goals: List[str] = ["wealth_building"]
    age_bracket: str = "26-35"
    annual_income_bracket: str = "1000000-2000000"
    investment_experience: str = "high"  # low, medium, high
    risk_capacity: str = "high"  # low, medium, high

class PortfolioResponse(BaseModel):
    portfolio: List[Dict[str, Any]]
    status: str
    message: str

# Initialize tools
search_tool = SerperDevTool()
scrape_tool = ScrapeWebsiteTool()

# Stock tracking
my_stocks = {}

class PurchaseStockTool(BaseTool):
    name: str = "Stock purchasing Tool"
    description: str = "Purchases Stock"

    def _run(self, stock, amount) -> str:
        my_stocks[stock] = my_stocks.get(stock, 0) + amount
        return True

purchase_stock_tool = PurchaseStockTool()

def get_llm():
    """Get LLM instance with Anthropic as primary and OpenAI as fallback"""
    
    # Try Anthropic first
    if ANTHROPIC_AVAILABLE and os.getenv("ANTHROPIC_API_KEY"):
        try:
            llm = ChatAnthropic(
                model="claude-3-5-sonnet-20241022",  # <-- ANTHROPIC MODEL SPECIFIED HERE
                anthropic_api_key=os.getenv("ANTHROPIC_API_KEY"),
                temperature=0.1,
                max_tokens=4000,
                timeout=60,
            )
            print("Using Anthropic Claude 3.5 Sonnet")
            return llm
        except Exception as e:
            print(f"Failed to initialize Anthropic: {e}")
    
    # Fallback to OpenAI
    if os.getenv("OPENAI_API_KEY"):
        try:
            llm = ChatOpenAI(
                model=os.getenv("OPENAI_MODEL_NAME", "gpt-4o"),  # <-- OPENAI MODEL SPECIFIED HERE
                openai_api_key=os.getenv("OPENAI_API_KEY"),
                temperature=0.1,
                max_tokens=4000,
                timeout=60,
                max_retries=2
            )
            print("Using OpenAI as fallback")
            return llm
        except Exception as e:
            print(f"Failed to initialize OpenAI: {e}")
    
    raise Exception("No valid LLM configuration found. Please set ANTHROPIC_API_KEY or OPENAI_API_KEY")

def create_agents():
    """Create and return all agents with Anthropic/OpenAI LLM"""
    
    # Get the LLM instance
    llm = get_llm()
    
    economic_research_agent = Agent(
        role="Senior Economic Analyst",
        goal='Analyze macroeconomic trends, assess current economic conditions, and forecast future economic scenarios to guide investment strategy',
        backstory="""You are a distinguished macroeconomist with 15+ years of experience analyzing
        global and US economic trends. You have worked at top-tier investment banks and central banks,
        specializing in economic forecasting, monetary policy analysis, and sector rotation strategies.
        You excel at connecting macroeconomic indicators to market performance and identifying
        sectors that will outperform or underperform based on economic cycles.""",
        verbose=True,
        allow_delegation=False,
        tools=[scrape_tool, search_tool],
        llm=llm
    )

    market_research_agent = Agent(
        role="Market Research Analyst",
        goal='Research and identify promising US equities based on market trends, news, and fundamental factors',
        backstory="""You are an expert market researcher with 10+ years of experience in identifying
        investment opportunities in US equity markets. You excel at analyzing market trends,
        company fundamentals, and industry dynamics.""",
        verbose=True,
        allow_delegation=False,
        tools=[scrape_tool, search_tool],
        llm=llm
    )

    fundamental_analyst_agent = Agent(
        role='Fundamental Analyst',
        goal='Analyze financial metrics, valuation, and company fundamentals for selected stocks',
        backstory="""You are a seasoned fundamental analyst with expertise in financial statement
        analysis, valuation models, and company assessment. You can evaluate P/E ratios,
        revenue growth, profit margins, and other key financial metrics.""",
        verbose=True,
        allow_delegation=False,
        tools=[scrape_tool, search_tool],
        llm=llm
    )

    risk_analyst_agent = Agent(
        role='Risk Management Analyst',
        goal='Assess risk metrics, volatility, correlations, and risk-adjusted returns for portfolio construction',
        backstory="""You are a quantitative risk analyst with deep expertise in portfolio theory,
        risk metrics, and volatility analysis. You specialize in calculating beta, correlation matrices,
        and risk-adjusted performance measures.""",
        verbose=True,
        allow_delegation=False,
        llm=llm
    )

    portfolio_manager_agent = Agent(
        role='Portfolio Manager',
        goal='Construct optimal portfolio weights based on research, fundamentals, risk analysis, and user risk profile',
        backstory="""You are an experienced portfolio manager with 15+ years in asset management.
        You excel at combining quantitative analysis with qualitative insights to create
        well-balanced portfolios that match investor risk profiles and objectives.""",
        verbose=True,
        allow_delegation=False,
        llm=llm
    )

    validation_agent = Agent(
        role='Portfolio Validation Specialist',
        goal='Validate and correct stock tickers, ensure proper formatting, and verify portfolio output',
        backstory="""You are a meticulous validation specialist with expertise in US stock market
        tickers and portfolio formatting. You ensure all stock symbols are valid NYSE/NASDAQ tickers
        and that portfolio outputs meet exact specifications. You have comprehensive knowledge of
        US stock ticker symbols and can identify and correct any invalid or generic references.""",
        tools=[search_tool],
        verbose=True,
        allow_delegation=False,
        llm=llm
    )
    
    return [economic_research_agent, market_research_agent, fundamental_analyst_agent, 
            risk_analyst_agent, portfolio_manager_agent, validation_agent]

def create_tasks(user_profile: UserProfile):
    """Create and return all tasks with user profile"""
    
    economy_research_task = Task(
        description="""
        Conduct a comprehensive macroeconomic analysis to inform investment strategy decisions.

        Your analysis should cover:

        **Current Economic State Assessment:**
        1. GDP growth trends and forecasts
        2. Inflation rates and Federal Reserve policy stance
        3. Employment data and labor market conditions
        4. Consumer spending patterns and confidence indices
        5. Corporate earnings trends across major sectors
        6. Interest rate environment and yield curve analysis

        **Future Economic Outlook (6-18 months):**
        1. Recession probability and economic cycle positioning
        2. Expected Federal Reserve actions and monetary policy trajectory
        3. Fiscal policy impacts and government spending priorities
        4. Global economic factors affecting the US (trade, geopolitics, etc.)
        5. Key economic risks and opportunities

        **Sector Analysis and Predictions:**
        1. Identify 3-4 sectors likely to outperform in the current/forecasted environment
        2. Identify 2-3 sectors likely to underperform or face headwinds
        3. Explain the economic drivers behind these sector predictions
        4. Consider both cyclical and secular trends affecting each sector

        **Investment Implications:**
        1. Recommend investment themes aligned with economic outlook
        2. Suggest defensive vs. growth-oriented positioning
        3. Comment on small-cap vs. large-cap preferences
        4. Provide guidance on sector allocation priorities

        Use current economic data, recent Fed communications, earnings reports, and reputable
        economic research sources. Focus on actionable insights that can guide portfolio construction.
        """,
        agent=create_agents()[0],  # economic_research_agent
        expected_output="""Comprehensive economic analysis report including:
        - Current economic state summary with key metrics
        - 6-18 month economic forecast with probability assessments
        - Sector outperformer/underperformer predictions with rationale
        - Investment themes and portfolio positioning recommendations
        - Key economic risks and catalysts to monitor"""
    )

    research_task = Task(
        description=f"""
        Research and identify 8-12 promising US equity investment opportunities suitable for a user with:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Investment timeline: {user_profile.investment_timeline}
        - Financial goals: {user_profile.financial_goals}
        - Age bracket: {user_profile.age_bracket}
        - Investment experience: {user_profile.investment_experience}

        Focus on:
        1. Large and mid-cap US stocks
        2. Different sectors for diversification
        3. Companies with strong fundamentals and growth prospects
        4. Consider current market conditions and trends

        Provide a list of stock symbols with brief rationale for each selection.
        """,
        agent=create_agents()[1],  # market_research_agent
        expected_output="List of 8-12 US stock symbols with brief investment rationale for each"
    )

    fundamental_analysis_task = Task(
        description="""
        Perform fundamental analysis on the stocks identified by the Market Research Analyst.
        For each stock, analyze:

        1. Financial metrics (P/E, P/B, ROE, debt-to-equity)
        2. Revenue and earnings growth trends
        3. Profit margins and efficiency ratios
        4. Competitive position and market share
        5. Management quality and corporate governance

        Rank the stocks based on fundamental strength and provide scores (1-10) for each.
        """,
        agent=create_agents()[2],  # fundamental_analyst_agent
        expected_output="Fundamental analysis scores and rankings for each stock with key metrics"
    )

    risk_analysis_task = Task(
        description=f"""
        Perform comprehensive risk analysis on the selected stocks considering user risk profile:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Risk capacity: {user_profile.risk_capacity}

        Calculate and analyze:
        1. Historical volatility (1-year and 3-year)
        2. Beta coefficients relative to S&P 500
        3. Correlation matrix between stocks
        4. Maximum drawdown analysis
        5. Risk-adjusted returns (Sharpe ratio)
        6. Sector concentration risk

        Provide risk scores and recommendations for position sizing based on user's risk profile.
        """,
        agent=create_agents()[3],  # risk_analyst_agent
        expected_output="Risk analysis with volatility metrics, correlations, and position sizing recommendations"
    )

    portfolio_construction_task = Task(
        description=f"""
        Construct an optimal portfolio using insights from market research, fundamental analysis, and risk assessment.

        User Profile:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Investment timeline: {user_profile.investment_timeline}
        - Investment experience: {user_profile.investment_experience}
        - Financial goals: {user_profile.financial_goals}

        Portfolio Requirements:
        1. Select 5-8 best stocks from the analyzed list
        2. Assign weights that sum to exactly 1.0
        3. Balance growth potential with risk management
        4. Consider diversification across sectors
        5. Match the user's risk tolerance and timeline

        CRITICAL: Output must be in this exact JSON format:
        [
            {{"symbol": "STOCK1", "weight": 0.XX}},
            {{"symbol": "STOCK2", "weight": 0.XX}},
            {{"symbol": "STOCK3", "weight": 0.XX}}
        ]

        Weights must sum to exactly 1.0.
        """,
        agent=create_agents()[4],  # portfolio_manager_agent
        expected_output="JSON array with stock symbols and weights that sum to 1.0"
    )

    validation_task = Task(
        description="""
        Validate and finalize the portfolio output from the Portfolio Manager.

        Your critical responsibilities:
        1. Verify ALL stock symbols are valid US ticker symbols (NYSE/NASDAQ)
        2. Replace any generic references like "Stock A", "Company X", "STOCK1" with actual tickers
        3. Ensure all symbols are properly formatted (uppercase, no spaces)
        4. Verify weights sum to exactly 1.0
        5. Ensure 5-8 stocks maximum in final portfolio

        If you find invalid tickers:
        - Research and replace with appropriate real US stock tickers
        - Maintain the same sector/style allocation intended by the portfolio manager
        - Keep the same relative weight proportions

        FINAL OUTPUT REQUIREMENTS:
        - Must be valid JSON array format
        - Only real, tradeable US stock tickers
        - Weights must sum to exactly 1.0
        - 5-8 stocks maximum

        CRITICAL: Your final output must be EXACTLY this format:
        [
            {"symbol": "AAPL", "weight": 0.30},
            {"symbol": "MSFT", "weight": 0.25},
            {"symbol": "GOOGL", "weight": 0.45}
        ]

        Return ONLY the JSON array, no additional text or explanation.
        """,
        agent=create_agents()[5],  # validation_agent
        expected_output="Final validated JSON array with real stock tickers and weights summing to 1.0"
    )

    return [economy_research_task, research_task, fundamental_analysis_task, 
            risk_analysis_task, portfolio_construction_task, validation_task]

@app.get("/")
async def root():
    return {"message": "CrewAI Portfolio Manager API with Anthropic", "status": "running"}

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

@app.post("/generate-portfolio")
async def generate_portfolio(user_profile: UserProfile):
    """Generate portfolio recommendations based on user profile"""
    
    try:
        # Set environment variables for CrewAI
        # Anthropic will be used if available, otherwise fallback to OpenAI
        if not os.getenv("ANTHROPIC_API_KEY") and not os.getenv("OPENAI_API_KEY"):
            raise HTTPException(status_code=500, detail="Either ANTHROPIC_API_KEY or OPENAI_API_KEY must be configured")
        
        os.environ["SERPER_API_KEY"] = os.getenv("SERPER_API_KEY")
        if not os.environ.get("SERPER_API_KEY"):
            raise HTTPException(status_code=500, detail="SERPER_API_KEY not configured")
        
        print(f"Starting portfolio generation for user profile: {user_profile}")
        
        # Create agents and tasks
        agents = create_agents()
        tasks = create_tasks(user_profile)
        
        # Create crew
        crew = Crew(
            agents=agents,
            tasks=tasks,
            process=Process.sequential,
            verbose=True
        )
        
        print("Executing CrewAI workflow...")
        result = crew.kickoff()
        
        print(f"CrewAI Result: {result}")
        print(f"My Stocks: {my_stocks}")
        
        # Try to parse the result as JSON portfolio
        try:
            import json
            import re
            
            result_str = str(result)
            print(f"Raw result string: {result_str}")
            
            # Extract JSON from markdown code blocks if present
            json_match = re.search(r'```json\s*(.*?)\s*```', result_str, re.DOTALL)
            if json_match:
                json_str = json_match.group(1).strip()
                print(f"Extracted JSON: {json_str}")
                # Clean up the JSON string - remove extra whitespace and newlines
                json_str = re.sub(r'\s+', ' ', json_str)  # Replace multiple whitespace with single space
                json_str = re.sub(r'\s*([{}[\],:])\s*', r'\1', json_str)  # Remove spaces around JSON syntax
                print(f"Cleaned JSON: {json_str}")
                portfolio_data = json.loads(json_str)
            else:
                # Try to find JSON array directly
                json_match = re.search(r'\[.*\]', result_str, re.DOTALL)
                if json_match:
                    json_str = json_match.group(0).strip()
                    # Clean up the JSON string
                    json_str = re.sub(r'\s+', ' ', json_str)
                    json_str = re.sub(r'\s*([{}[\],:])\s*', r'\1', json_str)
                    print(f"Found and cleaned JSON array: {json_str}")
                    portfolio_data = json.loads(json_str)
                else:
                    raise ValueError("No JSON found in result")
            
            if isinstance(portfolio_data, list) and len(portfolio_data) > 0:
                portfolio = portfolio_data
            else:
                portfolio = [{"symbol": "SAMPLE", "weight": 1.0}]
                
        except Exception as e:
            print(f"Error parsing portfolio JSON: {str(e)}")
            # Fallback if parsing fails
            portfolio = [{"symbol": "SAMPLE", "weight": 1.0}]
        
        return portfolio
        
    except Exception as e:
        print(f"Error generating portfolio: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error generating portfolio: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)