from fastapi import FastAPI, HTTPException, Header
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Dict, Optional
import anthropic
import os
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, specify actual domains
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize Anthropic client with error handling
anthropic_api_key = os.getenv("ANTHROPIC_API_KEY")
if not anthropic_api_key:
    print("Warning: ANTHROPIC_API_KEY not found in environment variables")
    claude_client = None
else:
    try:
        # Initialize with minimal configuration to avoid compatibility issues
        claude_client = anthropic.Anthropic(
            api_key=anthropic_api_key,
            max_retries=2,
            timeout=30.0
        )
        print("Anthropic client initialized successfully")
    except Exception as e:
        print(f"Error initializing Anthropic client: {e}")
        claude_client = None

if claude_client:
    print(f"Anthropic version: {anthropic.__version__}")
    print(f"Client type: {type(claude_client)}")
    print(f"Has beta: {hasattr(claude_client, 'beta')}")
    if hasattr(claude_client, 'beta'):
        print(f"Beta type: {type(claude_client.beta)}")
        print(f"Beta attributes: {dir(claude_client.beta)}")
        print(f"Has messages: {hasattr(claude_client.beta, 'messages')}")
    else:
        print("No beta attribute found!")

class ChatRequest(BaseModel):
    messages: List[Dict[str, str]]
    url: str

class ChatResponse(BaseModel):
    response: str
    
def format_anthropic_response(response):
    """Format Anthropic API response with MCP tool calls for easy reading"""
    
    print("=" * 50)
    print("ANTHROPIC RESPONSE")
    print("=" * 50)
    
    for i, block in enumerate(response.content):
        print(f"\n--- Block {i+1}: {block.type} ---")
        
        if block.type == 'text':
            print(f"Text: {block.text}")
            
        elif block.type == 'mcp_tool_use':
            print(f"Tool Called: {block.name}")
            print(f"Server: {block.server_name}")
            print(f"Tool ID: {block.id}")
            print(f"Input: {block.input}")
            
        elif block.type == 'mcp_tool_result':
            print(f"Tool Result (ID: {block.tool_use_id}):")
            print(f"Error: {block.is_error}")
            for content in block.content:
                if content.type == 'text':
                    print(f"Result: {content.text}")
    
    print("\n" + "=" * 50)

# Usage example:
# format_anthropic_response(response)

def get_final_text(response):
    """Extract just the final text response from Claude"""
    text_blocks = [block.text for block in response.content if block.type == 'text']
    return text_blocks[-1] if text_blocks else "No final text found"

# Usage:
# print("Final answer:", get_final_text(response))

@app.post("/chat")
async def chat(request: ChatRequest, x_account_id: Optional[str] = Header(None, alias="X-Account-ID")):
    try:
        if not x_account_id:
            raise HTTPException(status_code=400, detail="Account ID is required")
            
        if claude_client is None:
            raise HTTPException(status_code=500, detail="Anthropic client not initialized. Please check API key configuration.")
            
        print(f"Received request for account: {x_account_id}")  # Debug log
        
        response = claude_client.beta.messages.create(
            model="claude-sonnet-4-20250514",
            max_tokens=1000,
            system=f"You are a portfolio manager with more than 15 years of experience specializing in US equities. Current user account id is {x_account_id}",
            messages=request.messages,
            mcp_servers=[
                {
                    "type": "url",
                    "url": f"{request.url}/sse",
                    "name": "dice-server",
                }
            ],
            tools=[{
                "type": "web_search_20250305",
                "name": "web_search",
                "max_uses": 5
            }],
            extra_headers={
                "anthropic-beta": "mcp-client-2025-04-04"
            }
        )
        print("response:", get_final_text(response))
        return {"response": get_final_text(response)}
    
    except Exception as e:
        print(f"Error occurred: {str(e)}")  # Debug log
        print(f"Error type: {type(e)}")     # Debug log
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health():
    return {"status": "healthy", "service": "zeus-backend"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=3002)