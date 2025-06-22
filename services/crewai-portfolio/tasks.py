from crewai import Task
from models import UserProfile, RebalanceRequest
from agents import create_portfolio_agents, create_rebalance_agents

def create_portfolio_tasks(user_profile: UserProfile):
    """Create and return all tasks for portfolio generation with user profile"""
    
    agents = create_portfolio_agents()
    economic_research_agent, market_research_agent, fundamental_analyst_agent, risk_analyst_agent, portfolio_manager_agent, validation_agent = agents
    
    economy_research_task = Task(
        description="""
        Conduct a comprehensive macroeconomic analysis to inform investment strategy decisions.
        Your analysis should cover current economic state, future outlook, and sector predictions.
        """,
        agent=economic_research_agent,
        expected_output="Comprehensive economic analysis report with sector predictions and investment themes"
    )

    research_task = Task(
        description=f"""
        Research and identify 8-12 promising US equity investment opportunities suitable for:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Investment timeline: {user_profile.investment_timeline}
        - Financial goals: {user_profile.financial_goals}
        
        Consider both individual stocks AND ETFs based on user profile.
        """,
        agent=market_research_agent,
        expected_output="List of 8-12 US equity symbols with investment rationale"
    )

    fundamental_analysis_task = Task(
        description="Perform fundamental analysis on the identified stocks and ETFs.",
        agent=fundamental_analyst_agent,
        expected_output="Fundamental analysis scores and rankings for each security"
    )

    risk_analysis_task = Task(
        description=f"""
        Perform comprehensive risk analysis considering user risk profile:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Risk capacity: {user_profile.risk_capacity}
        """,
        agent=risk_analyst_agent,
        expected_output="Risk analysis with volatility metrics and position sizing recommendations"
    )

    portfolio_construction_task = Task(
        description=f"""
        Construct an optimal portfolio using insights from previous analysis.
        
        User Profile:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Investment timeline: {user_profile.investment_timeline}
        - Investment experience: {user_profile.investment_experience}
        - Financial goals: {user_profile.financial_goals}

        CRITICAL: Output must be in this exact JSON format:
        {{
            "portfolio": [
                {{"symbol": "STOCK1", "weight": 0.XX}},
                {{"symbol": "STOCK2", "weight": 0.XX}}
            ],
            "explanation": "Detailed explanation of the portfolio strategy and rationale."
        }}
        """,
        agent=portfolio_manager_agent,
        expected_output="JSON object with portfolio array and detailed explanation"
    )

    validation_task = Task(
        description="""
        Validate and finalize the portfolio output. Ensure all tickers are valid US symbols.
        Return ONLY a valid JSON object with real tickers and weights summing to 1.0.
        """,
        agent=validation_agent,
        expected_output="Final validated JSON object with real stock tickers"
    )

    return [economy_research_task, research_task, fundamental_analysis_task, 
            risk_analysis_task, portfolio_construction_task, validation_task]

def create_rebalance_tasks(rebalance_request: RebalanceRequest):
    """Create and return all tasks for portfolio rebalancing"""
    
    agents = create_rebalance_agents()
    portfolio_performance_analyst, market_conditions_analyst, risk_drift_analyst, rebalance_strategist, rebalance_validation_agent = agents
    
    user_profile = rebalance_request.user_profile
    current_portfolio = rebalance_request.current_portfolio
    
    # Format current portfolio for tasks
    portfolio_summary = ""
    for holding in current_portfolio:
        portfolio_summary += f"- {holding.symbol}: {holding.weight*100:.1f}% allocation, {holding.performance} performance\n"
    
    performance_analysis_task = Task(
        description=f"""
        Analyze the performance of the current portfolio and identify drift or imbalances.

        Current Portfolio Holdings:
        {portfolio_summary}
        
        User Profile:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Investment timeline: {user_profile.investment_timeline}

        Assess performance attribution, allocation drift, and concentration risk.
        """,
        agent=portfolio_performance_analyst,
        expected_output="Portfolio performance analysis with specific adjustment recommendations"
    )

    market_conditions_task = Task(
        description=f"""
        Assess current market conditions and determine if tactical adjustments are needed.
        Consider market regime, sector rotation, and risk environment changes.
        User has {user_profile.risk_tolerance} risk tolerance.
        """,
        agent=market_conditions_analyst,
        expected_output="Market conditions assessment with tactical allocation recommendations"
    )

    risk_drift_task = Task(
        description=f"""
        Analyze how the portfolio's risk characteristics have evolved.
        
        Current Portfolio:
        {portfolio_summary}
        
        Target Risk Profile: {user_profile.risk_tolerance} risk tolerance
        
        Assess volatility, correlations, and concentration risk changes.
        """,
        agent=risk_drift_analyst,
        expected_output="Risk drift analysis with risk-based rebalancing recommendations"
    )

    rebalance_strategy_task = Task(
        description=f"""
        Determine the optimal rebalancing strategy based on all previous analysis.

        CRITICAL: Output must be in this EXACT format:
        {{
            "rebalance_needed": true/false,
            "portfolio": [
                {{"symbol": "TICKER1", "weight": 0.XX}}
            ],
            "explanation": "Explanation of changes made and rationale",
            "changes_made": ["List of specific changes and reasons"]
        }}

        If no rebalancing needed, set rebalance_needed to false and portfolio to empty array.
        """,
        agent=rebalance_strategist,
        expected_output="Rebalancing strategy recommendation in specified JSON format"
    )

    rebalance_validation_task = Task(
        description="""
        Validate the rebalancing recommendation and ensure proper formatting.
        Verify all tickers are valid US symbols and weights sum to 1.0 if rebalancing.
        
        Return ONLY a valid JSON object in the exact format specified.
        """,
        agent=rebalance_validation_agent,
        expected_output="Final validated JSON object with rebalancing recommendation"
    )

    return [performance_analysis_task, market_conditions_task, risk_drift_task, 
            rebalance_strategy_task, rebalance_validation_task] 