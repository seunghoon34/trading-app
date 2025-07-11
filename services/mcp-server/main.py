import random
from fastmcp import FastMCP
import requests
import os

mcp = FastMCP(name="pandora-server")
AUTH_HEADER = os.getenv("ALPACA_AUTH_HEADER")


@mcp.tool
def get_user_trade_position(account_id: str):
    """Get the user's trade position for a given account ID."""

    url = f"https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/{account_id}/positions"

    headers = {
        "accept": "application/json",
        "authorization": AUTH_HEADER
}

    response = requests.get(url, headers=headers)
    return response.text

@mcp.tool
def buy_user_qty_stock(symbol: str, qty: str , account_id: str):
    """Buy an amounts(quantity) of stocks for the User with the given account ID a"""
    url = f"https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/{account_id}/orders"

    payload = {
        "type": "market",
        "time_in_force": "day",
        "commission_type": "notional",
        "symbol": symbol,
        "qty": qty,
        "side": "buy"
    }
    headers = {
        "accept": "application/json",
        "content-type": "application/json",
        "authorization": AUTH_HEADER
}
@mcp.tool
def buy_user_amount_stock(symbol: str, amount: str , account_id: str):
    """Buy an amounts(dollars) of stocks for the User with the given account ID a"""
    url = f"https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/{account_id}/orders"

    payload = {
        "type": "market",
        "time_in_force": "day",
        "commission_type": "notional",
        "symbol": symbol,
        "amount": amount,
        "side": "buy"
    }
    headers = {
        "accept": "application/json",
        "content-type": "application/json",
        "authorization": AUTH_HEADER
}
@mcp.tool
def sell_user_qty_stock(symbol: str, qty: str , account_id: str):
    """Buy an amounts(quantity) of stocks for the User with the given account ID a"""
    url = f"https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/{account_id}/orders"

    payload = {
        "type": "market",
        "time_in_force": "day",
        "commission_type": "notional",
        "symbol": symbol,
        "qty": qty,
        "side": "sell"
    }
    headers = {
        "accept": "application/json",
        "content-type": "application/json",
        "authorization": AUTH_HEADER
}
@mcp.tool
def sell_user_amount_stock(symbol: str, amount: str , account_id: str):
    """Buy an amounts(dollars) of stocks for the User with the given account ID a"""
    url = f"https://broker-api.sandbox.alpaca.markets/v1/trading/accounts/{account_id}/orders"

    payload = {
        "type": "market",
        "time_in_force": "day",
        "commission_type": "notional",
        "symbol": symbol,
        "amount": amount,
        "side": "sell"
    }
    headers = {
        "accept": "application/json",
        "content-type": "application/json",
        "authorization": AUTH_HEADER
    }
    response = requests.post(url, json=payload, headers=headers)

    return response.text

    

if __name__ == "__main__":
    mcp.run(transport="sse", port=3003)