from pydantic import BaseModel, Field
from typing import List, Dict, Any

class Position(BaseModel):
    symbol: str = Field(..., description="Stock or ETF symbol")
    weight: float = Field(..., ge=0, le=1, description="Portfolio weight between 0 and 1")
    perf: float = Field(..., description="Performance as a decimal (e.g., 0.2 for +20%, -0.01 for -1%)")

class CurrentPortfolio(BaseModel):
    positions: List[Position] = Field(..., description="List of current portfolio positions")
    total_value: float = Field(..., gt=0, description="Total portfolio value in USD")

class UserProfile(BaseModel):
    risk_tolerance: str = "high"  # low, medium, high
    investment_timeline: str = "short_term"  # short_term, medium_term, long_term
    financial_goals: List[str] = ["wealth_building"]
    age_bracket: str = "26-35"
    annual_income_bracket: str = "1000000-2000000"
    investment_experience: str = "high"  # low, medium, high
    risk_capacity: str = "high"  # low, medium, high
    current_portfolio: CurrentPortfolio = None  # Optional current portfolio

class RebalanceRequest(BaseModel):
    current_portfolio: CurrentPortfolio = Field(..., description="Current portfolio positions and value")
    user_profile: UserProfile = Field(..., description="User profile and preferences")

class PortfolioResponse(BaseModel):
    portfolio: List[Dict[str, Any]]
    explanation: str
    status: str
    message: str 