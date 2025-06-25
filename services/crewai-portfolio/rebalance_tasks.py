from typing import Dict, List
from crewai import Task, Agent
from models import UserProfile, CurrentPortfolio

def create_rebalance_tasks(agents: Dict[str, Agent], user_profile: UserProfile, current_portfolio: CurrentPortfolio) -> List[Task]:
    """Create and return all tasks for portfolio rebalancing"""
    
    economic_analysis_task = Task(
        description=f"""
        Analyze current economic conditions and their impact on the existing portfolio positions.

        Current Portfolio Performance:
        {repr(current_portfolio.positions)}

        Your analysis should cover:

        **Economic Impact Assessment:**
        1. How current economic conditions affect each position
        2. Federal Reserve policy impact on portfolio holdings
        3. Sector exposure analysis given economic trends
        4. Currency and international exposure evaluation
        5. Interest rate sensitivity of current positions

        **Forward-Looking Analysis:**
        1. Expected economic changes that may require position adjustments
        2. Sector rotation opportunities based on economic cycle
        3. Asset class reallocation needs given economic outlook
        4. Geographic exposure recommendations
        5. Duration and fixed income positioning needs

        **Position-Specific Recommendations:**
        1. Identify positions most vulnerable to economic changes
        2. Highlight positions well-positioned for expected conditions
        3. Suggest specific position adjustments based on economic outlook
        4. Recommend sector weight modifications

        Provide specific rebalancing recommendations based on economic factors.
        """,
        agent=agents['economic_research'],
        expected_output="""Economic impact analysis including:
        - Impact of current conditions on existing positions
        - Forward-looking economic risks and opportunities
        - Specific position adjustment recommendations
        - Sector and asset class reallocation suggestions"""
    )

    performance_analysis_task = Task(
        description=f"""
        Analyze the performance of current portfolio positions and identify rebalancing needs.

        Current Portfolio:
        {repr(current_portfolio.positions)}

        **Performance Analysis:**
        1. Calculate absolute and relative performance for each position
        2. Identify significant performance outliers
        3. Analyze tracking error for ETF positions
        4. Calculate portfolio drift from target allocations
        5. Assess contribution and attribution metrics

        **Drift Analysis:**
        1. Calculate current vs target weights for all positions
        2. Identify positions requiring immediate rebalancing
        3. Analyze sector and asset class drift
        4. Evaluate style drift (growth vs value, market cap)

        **Rebalancing Triggers:**
        1. List positions exceeding drift thresholds
        2. Identify underweight and overweight positions
        3. Calculate rebalancing trade sizes needed
        4. Prioritize rebalancing actions

        Provide specific rebalancing recommendations based on performance and drift.
        """,
        agent=agents['performance_analyst'],
        expected_output="Detailed performance analysis with drift calculations and rebalancing recommendations"
    )

    market_conditions_task = Task(
        description=f"""
        Assess how current market conditions affect the portfolio and identify necessary adjustments.

        Current Portfolio:
        {repr(current_portfolio.positions)}

        **Market Impact Analysis:**
        1. Current market regime assessment
        2. Sector rotation impacts on holdings
        3. Factor exposure analysis
        4. Market sentiment impact on positions
        5. Liquidity conditions assessment
        6. Identify new sectors or assets that should be added for better market positioning

        **Position-Specific Analysis:**
        1. Market sensitivity of each holding
        2. Momentum and trend analysis
        3. Technical analysis signals
        4. Volume and liquidity patterns
        5. Volatility regime impact
        6. Identify specific new stocks or ETFs that could enhance the portfolio

        **Adjustment Recommendations:**
        1. Identify positions requiring adjustment due to market conditions
        2. Suggest tactical shifts based on market regime
        3. Recommend position size adjustments
        4. Provide entry/exit timing considerations
        5. Propose new positions that would benefit from current market conditions
        6. Suggest positions to completely exit if no longer suitable

        Focus on actionable rebalancing recommendations based on market conditions.
        Include specific recommendations for both existing and potential new positions.
        """,
        agent=agents['market_conditions'],
        expected_output="Market conditions analysis with specific rebalancing recommendations"
    )

    risk_analysis_task = Task(
        description=f"""
        Analyze portfolio risk metrics and recommend risk-based rebalancing actions.

        Current Portfolio:
        {repr(current_portfolio.positions)}

        User Risk Profile:
        - Risk tolerance: {user_profile.risk_tolerance}
        - Risk capacity: {user_profile.risk_capacity}

        **Risk Metrics Analysis:**
        1. Portfolio beta and correlation changes
        2. Value at Risk (VaR) analysis
        3. Volatility assessment
        4. Concentration risk evaluation
        5. Factor exposure analysis

        **Risk-Based Rebalancing Needs:**
        1. Identify positions with excessive risk
        2. Analyze portfolio-level risk metrics
        3. Evaluate diversification effectiveness
        4. Assess hedging requirements

        Provide specific rebalancing recommendations to optimize risk profile.
        """,
        agent=agents['risk_monitor'],
        expected_output="Risk analysis with rebalancing recommendations to maintain target risk levels"
    )

    rebalance_strategy_task = Task(
        description=f"""
        Develop comprehensive rebalancing strategy based on all analyses.

        Current Portfolio:
        {repr(current_portfolio.positions)}

        **Strategy Development:**
        1. Prioritize rebalancing actions including new position additions
        2. Calculate optimal trade sizes for both existing and new positions
        3. Determine trade sequencing including position exits and entries
        4. Consider trading costs and market impact
        5. Plan implementation timeline
        6. Evaluate liquidity for new position entries

        **Trade Planning:**
        1. List specific trades needed (including new positions)
        2. Provide target weights for each position (existing and new)
        3. Include price limits and timing considerations
        4. Account for position liquidity
        5. Plan position exits if completely removing holdings

        CRITICAL: Output must be in this exact JSON format:
        {{
            "portfolio": [
                {{"symbol": "AAPL", "weight": 0.XX, "action": "reduce/increase/hold/exit"}},
                {{"symbol": "VTI", "weight": 0.XX, "action": "reduce/increase/hold/exit"}},
                {{"symbol": "NVDA", "weight": 0.XX, "action": "new"}}
            ],
            "explanation": "Detailed explanation of rebalancing rationale, including specific reasons for each action, new position additions, complete exits, and how the changes align with economic outlook, market conditions, and risk parameters."
        }}

        Portfolio weights must sum to exactly 1.0.
        The explanation should be a concise summary (3-4 sentences).
        Use "new" action for newly added positions and "exit" for complete removals.
        """,
        agent=agents['rebalance_strategist'],
        expected_output="JSON object with rebalancing actions and detailed explanation"
    )

    validation_task = Task(
        description=f"""
        Validate and finalize the rebalancing recommendations.

        Current Portfolio:
        {repr(current_portfolio.positions)}

        **Validation Requirements:**
        1. Verify all symbols are valid and actively trading (including new additions)
        2. Ensure target weights sum to 1.0
        3. Validate trade sizes against liquidity constraints
        4. Check position size limits
        5. Verify diversification requirements
        6. Validate feasibility of new position entries
        7. Confirm exit positions can be fully liquidated

        **Final Output Requirements:**
        - Must be valid JSON format with "portfolio" and "explanation" keys
        - Portfolio array must include symbol, weight, and action
        - Weights must sum to exactly 1.0
        - Actions must be "reduce", "increase", "hold", "exit", or "new"
        - Explanation must be clear and concise
        - New positions must be validated for trading status and liquidity

        CRITICAL: Your final output must be EXACTLY this format:
        {{
            "portfolio": [
                {{"symbol": "AAPL", "weight": 0.20, "action": "reduce"}},
                {{"symbol": "VTI", "weight": 0.50, "action": "increase"}},
                {{"symbol": "MSFT", "weight": 0.30, "action": "new"}}
            ],
            "explanation": "Rebalancing actions focus on reducing overweight tech exposure, increasing broad market allocation, and adding new positions in high-quality growth companies. Changes consider recent economic data and maintain target risk levels while optimizing for minimal trading impact."
        }}
        """,
        agent=agents['execution_validator'],
        expected_output="Final validated JSON object with rebalancing actions and explanation"
    )

    return [
        economic_analysis_task,
        performance_analysis_task,
        market_conditions_task,
        risk_analysis_task,
        rebalance_strategy_task,
        validation_task
    ]