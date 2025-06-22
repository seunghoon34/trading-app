from crewai_tools import BaseTool, ScrapeWebsiteTool, SerperDevTool
from utils import my_stocks

# Initialize tools
search_tool = SerperDevTool()
scrape_tool = ScrapeWebsiteTool()

