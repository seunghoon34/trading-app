from typing import Dict
from crewai_tools import BaseTool, ScrapeWebsiteTool, SerperDevTool

def create_tools() -> Dict[str, BaseTool]:
    """Create and return all tools used by agents"""
    tools = {
        'search': SerperDevTool(),
        'scrape': ScrapeWebsiteTool()
    }
    return tools

# Initialize standard tools for direct import
search_tool = SerperDevTool()
scrape_tool = ScrapeWebsiteTool()

