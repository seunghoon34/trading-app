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
                model="claude-3-5-sonnet-20241022",  # <-- ANTHROPIC MODEL SPECIFIED HERE
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
                model=os.getenv("OPENAI_MODEL_NAME", "gpt-4o"),  # <-- OPENAI MODEL SPECIFIED HERE
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