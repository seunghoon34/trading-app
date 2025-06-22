import os
import warnings
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Warning control
warnings.filterwarnings('ignore')

# Import both Anthropic and OpenAI for fallback
try:
    from langchain_anthropic import ChatAnthropic
    ANTHROPIC_AVAILABLE = True
except ImportError:
    ANTHROPIC_AVAILABLE = False
    print("Warning: langchain_anthropic not installed. Install with: pip install langchain-anthropic")

from langchain_openai import ChatOpenAI

def get_llm():
    """Get LLM instance with Anthropic as primary and OpenAI as fallback"""
    
    # Try Anthropic first
    if ANTHROPIC_AVAILABLE and os.getenv("ANTHROPIC_API_KEY"):
        try:
            llm = ChatAnthropic(
                model="claude-3-5-sonnet-20241022",
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
                model=os.getenv("OPENAI_MODEL_NAME", "gpt-4o"),
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

# Stock tracking
my_stocks = {}

def parse_portfolio_json(result_str):
    """Parse portfolio JSON from CrewAI result string"""
    import json
    import re
    
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
        return json.loads(json_str)
    else:
        # Try to find JSON object directly
        json_match = re.search(r'\{.*\}', result_str, re.DOTALL)
        if json_match:
            json_str = json_match.group(0).strip()
            # Clean up the JSON string
            json_str = re.sub(r'\s+', ' ', json_str)
            json_str = re.sub(r'\s*([{}[\],:])\s*', r'\1', json_str)
            print(f"Found and cleaned JSON object: {json_str}")
            return json.loads(json_str)
        else:
            raise ValueError("No JSON found in result")

def validate_portfolio_response(portfolio_response):
    """Validate and normalize portfolio response structure"""
    if isinstance(portfolio_response, dict) and "portfolio" in portfolio_response and "explanation" in portfolio_response:
        return portfolio_response["portfolio"], portfolio_response["explanation"]
    elif isinstance(portfolio_response, list):
        # Fallback for old format
        return portfolio_response, "Portfolio generated based on comprehensive analysis of market conditions, user risk profile, and investment objectives."
    else:
        return [{"symbol": "SAMPLE", "weight": 1.0}], "Sample portfolio - analysis could not be completed." 