from typing import Dict, List
from crewai import Task, Agent
from models import UserProfile

def create_tasks(agents: Dict[str, Agent], user_profile: UserProfile) -> List[Task]:
    """Create and return all tasks with user profile and pre-created agents"""
    
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
        agent=agents['economic_research'],
        expected_output="""Comprehensive economic analysis report including:
        - Current economic state summary with key metrics
        - 6-18 month economic forecast with probability assessments
        - Sector outperformer/underperformer predictions with rationale
        - Investment themes and portfolio positioning recommendations
        - Key economic risks and catalysts to monitor"""
    )

    research_task = Task(
        description=f"""
        Research and identify 8-12 promising US equity investment opportunities (stocks and ETFs) suitable for a user with:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Investment timeline: {user_profile.investment_timeline}
        - Financial goals: {user_profile.financial_goals}
        - Age bracket: {user_profile.age_bracket}
        - Investment experience: {user_profile.investment_experience}

        Consider both individual stocks AND ETFs based on user profile:
        
        **For Individual Stocks:**
        1. Large and mid-cap US stocks with strong fundamentals
        2. Growth and value opportunities across different sectors
        3. Companies with competitive advantages and growth prospects
        
        **For ETFs:**
        1. Broad market ETFs (S&P 500, Total Market) for core holdings
        2. Sector-specific ETFs for targeted exposure (technology, healthcare, etc.)
        3. International developed market ETFs for global diversification
        4. Bond ETFs for stability (especially for conservative profiles)
        5. Thematic ETFs (clean energy, emerging technologies) for growth-oriented profiles
        
        **Selection Criteria:**
        - Match ETF vs stock allocation to user's experience level and risk tolerance
        - Consider expense ratios for ETFs (prefer low-cost options)
        - Balance individual stock picks with diversified ETF exposure
        - For less experienced investors, favor ETFs for core positions
        - For experienced investors, consider more individual stock picks
        
        Provide a mixed list of stock and ETF symbols with brief rationale for each selection.
        Include ticker symbols and specify whether each is a stock or ETF.
        """,
        agent=agents['market_research'],
        expected_output="List of 8-12 US equity symbols (mix of stocks and ETFs) with investment rationale and type specification for each"
    )

    fundamental_analysis_task = Task(
        description="""
        Perform fundamental analysis on the stocks and ETFs identified by the Market Research Analyst.
        
        **For Individual Stocks, analyze:**
        1. Financial metrics (P/E, P/B, ROE, debt-to-equity)
        2. Revenue and earnings growth trends
        3. Profit margins and efficiency ratios
        4. Competitive position and market share
        5. Management quality and corporate governance
        
        **For ETFs, analyze:**
        1. Expense ratios and fee structure
        2. Tracking error and performance vs benchmark
        3. Assets under management (AUM) and liquidity
        4. Holdings concentration and diversification
        5. Fund methodology and management approach
        6. Underlying asset quality and sector allocations
        
        **Additional ETF Considerations:**
        - For sector ETFs: sector outlook and cyclical positioning
        - For broad market ETFs: market cap exposure and style tilts
        - For international ETFs: geographic diversification benefits
        - For bond ETFs: duration, credit quality, and yield characteristics
        
        Rank all securities (stocks and ETFs) based on fundamental strength and appropriateness 
        for the user profile. Provide scores (1-10) for each with specific reasoning for stocks vs ETFs.
        """,
        agent=agents['fundamental_analyst'],
        expected_output="Fundamental analysis scores and rankings for each stock and ETF with appropriate metrics for each type"
    )

    risk_analysis_task = Task(
        description=f"""
        Perform comprehensive risk analysis on the selected stocks and ETFs considering user risk profile:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Risk capacity: {user_profile.risk_capacity}

        **For All Securities (Stocks & ETFs):**
        1. Historical volatility (1-year and 3-year)
        2. Beta coefficients relative to S&P 500
        3. Correlation matrix between all selected securities
        4. Maximum drawdown analysis
        5. Risk-adjusted returns (Sharpe ratio)
        
        **Additional Risk Considerations:**
        6. Sector concentration risk across the entire portfolio
        7. Geographic concentration risk (especially for international ETFs)
        8. Style risk (growth vs value, large vs small cap)
        9. ETF-specific risks: tracking error, liquidity risk, counterparty risk
        10. Single stock risk vs diversified ETF risk assessment
        
        **Position Sizing Recommendations:**
        - Consider ETFs for larger allocations due to diversification benefits
        - Limit individual stock positions based on concentration risk
        - Account for user's experience level in risk tolerance
        - Balance between diversified ETF exposure and focused stock picks
        - Consider correlation benefits of mixing stocks and ETFs

        Provide risk scores and position sizing recommendations that account for the different 
        risk characteristics of individual stocks versus ETFs.
        """,
        agent=agents['risk_analyst'],
        expected_output="Risk analysis with volatility metrics, correlations, and position sizing recommendations for stocks and ETFs"
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
        1. Select 5-8 best securities (mix of stocks and ETFs) from the analyzed list
        2. Assign weights that sum to exactly 1.0
        3. Balance growth potential with risk management
        4. Consider diversification across sectors and asset types
        5. Match the user's risk tolerance and timeline
        6. Use ETFs for core diversified exposure and individual stocks for targeted opportunities
        7. Consider user experience level when balancing ETF vs stock allocation

        CRITICAL: Output must be in this exact JSON format:
        {{
            "portfolio": [
                {{"symbol": "STOCK1", "weight": 0.XX}},
                {{"symbol": "STOCK2", "weight": 0.XX}},
                {{"symbol": "STOCK3", "weight": 0.XX}}
            ],
            "explanation": "Detailed explanation of stock selection rationale, weight allocation reasoning, risk profile alignment, and how this portfolio meets the user's financial goals and risk tolerance. Include specific reasons for each major holding and overall portfolio strategy."
        }}

        Portfolio weights must sum to exactly 1.0.
        The explanation should be a concise summary (3-4 sentences) covering:
        - Brief rationale for the stock/ETF mix and key holdings
        - How the portfolio aligns with user's risk tolerance and goals
        - Expected risk/return profile for the investment timeline
        """,
        agent=agents['portfolio_manager'],
        expected_output="JSON object with portfolio array and detailed explanation of investment rationale"
    )

    validation_task = Task(
        description="""
        Validate and finalize the portfolio output from the Portfolio Manager.

        Your critical responsibilities:
        1. Verify ALL symbols are valid US ticker symbols (NYSE/NASDAQ/ARCA for stocks and ETFs)
        2. Replace any generic references like "Stock A", "Company X", "STOCK1", "ETF1" with actual tickers
        3. Ensure all symbols are properly formatted (uppercase, no spaces)
        4. Verify weights sum to exactly 1.0
        5. Ensure 5-8 securities (mix of stocks and ETFs) maximum in final portfolio
        6. Validate that the explanation is comprehensive and accurately describes the portfolio
        7. Ensure appropriate balance between individual stocks and ETFs based on user profile

        If you find invalid tickers:
        - Research and replace with appropriate real US stock or ETF tickers
        - Maintain the same sector/style/asset type allocation intended by the portfolio manager
        - Keep the same relative weight proportions
        - For ETF replacements, ensure expense ratios are reasonable (typically under 0.75%)
        - Update the explanation to reflect any ticker changes

        FINAL OUTPUT REQUIREMENTS:
        - Must be valid JSON object format with "portfolio" and "explanation" keys
        - Only real, tradeable US stock and ETF tickers in portfolio array
        - Weights must sum to exactly 1.0
        - 5-8 securities (mix of stocks and ETFs) maximum
        - Explanation must be a concise 3-4 sentence summary

        CRITICAL: Your final output must be EXACTLY this format:
        {
            "portfolio": [
                {"symbol": "VTI", "weight": 0.40},
                {"symbol": "AAPL", "weight": 0.20},
                {"symbol": "QQQ", "weight": 0.25},
                {"symbol": "MSFT", "weight": 0.15}
            ],
            "explanation": "This portfolio combines broad market ETF exposure with high-quality individual stocks to match your risk tolerance. The allocation emphasizes growth potential while maintaining diversification through ETFs. Expected to deliver competitive returns over your investment timeline with managed risk."
        }

        Return ONLY the JSON object, no additional text or explanation outside the JSON.
        """,
        agent=agents['validation'],
        expected_output="Final validated JSON object with real stock tickers, weights summing to 1.0, and comprehensive explanation"
    )

    return [economy_research_task, research_task, fundamental_analysis_task, 
            risk_analysis_task, portfolio_construction_task, validation_task]