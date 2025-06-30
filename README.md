# Trading App

A microservices-based trading application built with Python and Go, using the Alpaca API for paper trading and learning advanced software engineering concepts.

## Architecture

- **API Gateway** (Go/Gin) - Request routing, authentication, rate limiting, and CORS handling
- **User Management** (Go/Gin) - User authentication, registration, and profile management
- **Market Data** (Go/Gin) - Real-time market data and quotes from Alpaca
- **Trading Engine** (Go/Gin) - Order placement, execution, and management
- **Portfolio** (Go/Gin) - Position tracking and portfolio analytics
- **Investment Strategy** (Go/Gin) - Portfolio creation, risk profiling, and investment recommendations
- **Payment** (Go/Gin) - Deposit and payment processing
- **Notification** (Go/Gin) - Real-time notifications via Kafka and MongoDB
- **Event Listener** (Go/Gin) - Alpaca SSE event processing and Kafka publishing
- **CrewAI Portfolio** (Python/FastAPI) - AI-powered portfolio management
- **MCP Server** (Python) - Model Context Protocol server
- **Zeus Backend** (Python) - AI assistant backend

## Demo

[![Pandora Demo](https://img.youtube.com/vi/hdR3aePfZQA/0.jpg)](https://youtu.be/hdR3aePfZQA?si=-Bb-GyyPCH6Qp42o)