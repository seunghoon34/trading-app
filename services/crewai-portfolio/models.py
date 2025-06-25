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

class PortfolioResponse(BaseModel):
    portfolio: List[Dict[str, Any]]
    explanation: str
    status: str
    message: str 