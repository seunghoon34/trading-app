from typing import Dict
from crewai import Agent
from config import get_llm
from tools import scrape_tool, search_tool

def create_rebalance_agents(llm) -> Dict[str, Agent]:
    """Create and return all agents focused on portfolio rebalancing"""
    
    agents = {
        'economic_research': Agent(
            role="Senior Economic Analyst",
            goal='Analyze macroeconomic trends and their impact on existing portfolio positions, identifying necessary adjustments based on economic outlook',
            backstory="""You are a distinguished macroeconomist with 15+ years of experience analyzing
            global and US economic trends. You excel at understanding how economic changes affect existing
            positions and when portfolio adjustments are needed due to changing economic conditions.
            You specialize in identifying when sector rotations or asset allocation shifts are needed
            based on monetary policy, fiscal policy, and economic cycle changes.""",
            verbose=True,
            allow_delegation=False,
            tools=[scrape_tool, search_tool],
            llm=llm
        ),
        'performance_analyst': Agent(
            role="Portfolio Performance Analyst",
            goal='Analyze current portfolio performance, identify underperforming and outperforming positions, and assess drift from target allocations',
            backstory="""You are a skilled portfolio analyst with 12+ years of experience in performance
            attribution and portfolio analytics. You excel at analyzing position-level and portfolio-level
            performance metrics, identifying sources of tracking error, and calculating portfolio drift.
            You understand how to measure both absolute and relative performance across different
            market conditions.""",
            verbose=True,
            allow_delegation=False,
            tools=[scrape_tool, search_tool],
            llm=llm
        ),
        'market_conditions': Agent(
            role="Market Conditions Analyst",
            goal='Assess current market conditions and their impact on existing positions, identify potential risks and opportunities',
            backstory="""You are an experienced market analyst specializing in understanding how changing
            market conditions affect existing portfolio positions. You excel at identifying market regime
            changes, sector rotations, and how they impact different types of holdings. You have deep
            expertise in analyzing both individual stocks and ETFs in the context of current market
            dynamics.""",
            verbose=True,
            allow_delegation=False,
            tools=[scrape_tool, search_tool],
            llm=llm
        ),
        'rebalance_strategist': Agent(
            role='Rebalancing Strategist',
            goal='Develop rebalancing strategy considering trading costs and optimal execution',
            backstory="""You are a portfolio strategist with 15+ years focusing on optimal rebalancing
            strategies. You understand how to reduce trading costs and maintain target allocations
            efficiently. You excel at determining when positions need adjustment and how to sequence
            trades for optimal execution.""",
            verbose=True,
            allow_delegation=False,
            tools=[scrape_tool, search_tool],
            llm=llm
        ),
        'risk_monitor': Agent(
            role='Risk Monitoring Specialist',
            goal='Monitor risk metrics changes and assess impact of proposed rebalancing on portfolio risk',
            backstory="""You are a risk specialist focused on portfolio monitoring and rebalancing.
            You excel at identifying when portfolio risk has drifted from targets, analyzing correlation
            changes between holdings, and ensuring that rebalancing actions maintain desired risk levels.
            You understand both position-specific and portfolio-level risk metrics.""",
            verbose=True,
            allow_delegation=False,
            llm=llm
        ),
        'execution_validator': Agent(
            role='Rebalancing Execution Validator',
            goal='Validate rebalancing recommendations and ensure they meet all constraints and requirements',
            backstory="""You are a meticulous execution specialist who validates rebalancing
            recommendations. You ensure all proposed trades are valid, verify that position changes
            maintain required diversification, and confirm that recommendations consider liquidity
            and trading constraints. You excel at catching potential issues before they become
            problems during execution.""",
            tools=[search_tool],
            verbose=True,
            allow_delegation=False,
            llm=llm
        )
    }
    
    return agents 