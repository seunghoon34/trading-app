from pydantic import BaseModel
from typing import List, Dict, Any

class UserProfile(BaseModel):
    risk_tolerance: str = "high"  # low, medium, high
    investment_timeline: str = "short_term"  # short_term, medium_term, long_term
    financial_goals: List[str] = ["wealth_building"]
    age_bracket: str = "26-35"
    annual_income_bracket: str = "1000000-2000000"
    investment_experience: str = "high"  # low, medium, high
    risk_capacity: str = "high"  # low, medium, high

class PortfolioHolding(BaseModel):
    symbol: str
    weight: float
    performance: str  # e.g., "+10%", "-5%"

class RebalanceRequest(BaseModel):
    user_profile: UserProfile
    current_portfolio: List[PortfolioHolding]
    time_since_last_rebalance: str = "3_months"  # 1_month, 3_months, 6_months, 1_year

class PortfolioResponse(BaseModel):
    portfolio: List[Dict[str, Any]]
    explanation: str
    status: str
    message: str

class RebalanceResponse(BaseModel):
    portfolio: List[Dict[str, Any]]
    explanation: str
    status: str
    message: str
    rebalance_needed: bool
    changes_made: List[str] = [] 