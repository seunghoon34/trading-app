from fastapi import FastAPI, HTTPException
import os
import warnings
import json
import re
from crewai import Crew, Process

# Import local modules
from config import get_llm
from models import UserProfile, PortfolioResponse
from agents import create_agents
from tasks import create_tasks

# Warning control
warnings.filterwarnings('ignore')

app = FastAPI(title="CrewAI Portfolio Manager", version="1.0.0")

def parse_portfolio_result(result_str: str) -> dict:
    """Parse the CrewAI result into a structured portfolio response"""
    
    try:
        # Clean the result string
        result_str = str(result_str).strip()
        
        # Try to extract JSON from markdown code blocks
        json_match = re.search(r'```(?:json)?\s*(.*?)\s*```', result_str, re.DOTALL | re.IGNORECASE)
        if json_match:
            json_str = json_match.group(1).strip()
        else:
            # Try to find JSON object directly
            json_match = re.search(r'\{.*\}', result_str, re.DOTALL)
            if json_match:
                json_str = json_match.group(0).strip()
            else:
                raise ValueError("No JSON found in result")
        
        # Parse the JSON
        portfolio_data = json.loads(json_str)
        
        # Validate structure
        if isinstance(portfolio_data, dict) and "portfolio" in portfolio_data and "explanation" in portfolio_data:
            return {
                "portfolio": portfolio_data["portfolio"],
                "explanation": portfolio_data["explanation"]
            }
        elif isinstance(portfolio_data, list):
            # Fallback for old format
            return {
                "portfolio": portfolio_data,
                "explanation": "Portfolio generated based on comprehensive analysis of market conditions, user risk profile, and investment objectives."
            }
        else:
            raise ValueError("Invalid portfolio structure")
            
    except Exception as e:
        print(f"Error parsing portfolio result: {e}")
        # Return fallback portfolio
        return {
            "portfolio": [{"symbol": "VTI", "weight": 1.0}],
            "explanation": "Sample portfolio - analysis could not be completed due to parsing error."
        }

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
        # Validate environment variables
        if not os.getenv("ANTHROPIC_API_KEY") and not os.getenv("OPENAI_API_KEY"):
            raise HTTPException(status_code=500, detail="Either ANTHROPIC_API_KEY or OPENAI_API_KEY must be configured")
        
        if not os.getenv("SERPER_API_KEY"):
            raise HTTPException(status_code=500, detail="SERPER_API_KEY not configured")
        
        print(f"Starting portfolio generation for user profile: {user_profile}")
        
        # Create agents and tasks
        llm = get_llm()
        agents = create_agents(llm)
        tasks = create_tasks(agents, user_profile)
        
        # Create crew
        crew = Crew(
            agents=list(agents.values()),
            tasks=tasks,
            process=Process.sequential,
            verbose=True
        )
        
        print("Executing CrewAI workflow...")
        result = crew.kickoff()
        
        print(f"CrewAI Result: {result}")
        
        # Parse the result
        portfolio_data = parse_portfolio_result(result)
        
        return PortfolioResponse(
            portfolio=portfolio_data["portfolio"],
            explanation=portfolio_data["explanation"],
            status="success",
            message="Portfolio generated successfully"
        )
        
    except Exception as e:
        print(f"Error generating portfolio: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error generating portfolio: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)