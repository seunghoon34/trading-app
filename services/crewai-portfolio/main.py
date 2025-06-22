from fastapi import FastAPI, HTTPException
import os
from crewai import Crew, Process

# Import modular components
from models import UserProfile, PortfolioResponse, RebalanceRequest, RebalanceResponse
from agents import create_portfolio_agents, create_rebalance_agents
from tasks import create_portfolio_tasks, create_rebalance_tasks
from utils import get_llm, parse_portfolio_json, validate_portfolio_response, my_stocks

app = FastAPI(title="CrewAI Portfolio Manager", version="1.0.0")

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
        if not os.getenv("ANTHROPIC_API_KEY") and not os.getenv("OPENAI_API_KEY"):
            raise HTTPException(status_code=500, detail="Either ANTHROPIC_API_KEY or OPENAI_API_KEY must be configured")
        
        os.environ["SERPER_API_KEY"] = os.getenv("SERPER_API_KEY")
        if not os.environ.get("SERPER_API_KEY"):
            raise HTTPException(status_code=500, detail="SERPER_API_KEY not configured")
        
        print(f"Starting portfolio generation for user profile: {user_profile}")
        
        # Create agents and tasks
        agents = create_portfolio_agents()
        tasks = create_portfolio_tasks(user_profile)
        
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
            portfolio_response = parse_portfolio_json(str(result))
            portfolio, explanation = validate_portfolio_response(portfolio_response)
        except Exception as e:
            print(f"Error parsing portfolio JSON: {str(e)}")
            # Fallback if parsing fails
            portfolio = [{"symbol": "SAMPLE", "weight": 1.0}]
            explanation = "Sample portfolio - analysis could not be completed due to parsing error."
        
        return PortfolioResponse(
            portfolio=portfolio,
            explanation=explanation,
            status="success",
            message="Portfolio generated successfully"
        )
        
    except Exception as e:
        print(f"Error generating portfolio: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error generating portfolio: {str(e)}")

@app.post("/rebalance-portfolio")
async def rebalance_portfolio(rebalance_request: RebalanceRequest):
    """Analyze and rebalance existing portfolio based on current conditions"""
    
    try:
        # Set environment variables for CrewAI
        if not os.getenv("ANTHROPIC_API_KEY") and not os.getenv("OPENAI_API_KEY"):
            raise HTTPException(status_code=500, detail="Either ANTHROPIC_API_KEY or OPENAI_API_KEY must be configured")
        
        os.environ["SERPER_API_KEY"] = os.getenv("SERPER_API_KEY")
        if not os.environ.get("SERPER_API_KEY"):
            raise HTTPException(status_code=500, detail="SERPER_API_KEY not configured")
        
        print(f"Starting portfolio rebalancing for request: {rebalance_request}")
        
        # Create agents and tasks for rebalancing
        agents = create_rebalance_agents()
        tasks = create_rebalance_tasks(rebalance_request)
        
        # Create crew
        crew = Crew(
            agents=agents,
            tasks=tasks,
            process=Process.sequential,
            verbose=True
        )
        
        print("Executing CrewAI rebalancing workflow...")
        result = crew.kickoff()
        
        print(f"CrewAI Rebalancing Result: {result}")
        
        # Try to parse the result as JSON rebalancing response
        try:
            rebalance_response = parse_portfolio_json(str(result))
            
            # Validate the rebalance response structure
            if isinstance(rebalance_response, dict):
                rebalance_needed = rebalance_response.get("rebalance_needed", False)
                portfolio = rebalance_response.get("portfolio", [])
                explanation = rebalance_response.get("explanation", "Portfolio analysis completed.")
                changes_made = rebalance_response.get("changes_made", [])
            else:
                # Fallback structure
                rebalance_needed = False
                portfolio = []
                explanation = "Portfolio analysis completed - no rebalancing needed."
                changes_made = []
                
        except Exception as e:
            print(f"Error parsing rebalancing JSON: {str(e)}")
            # Fallback if parsing fails
            rebalance_needed = False
            portfolio = []
            explanation = "Portfolio analysis could not be completed due to parsing error."
            changes_made = []
        
        return RebalanceResponse(
            portfolio=portfolio,
            explanation=explanation,
            status="success",
            message="Portfolio rebalancing analysis completed",
            rebalance_needed=rebalance_needed,
            changes_made=changes_made
        )
        
    except Exception as e:
        print(f"Error rebalancing portfolio: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error rebalancing portfolio: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)