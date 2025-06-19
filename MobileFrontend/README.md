# Trading App - Mobile Frontend

A React Native mobile application that mirrors the functionality and design of the web frontend, providing a seamless trading experience on mobile devices.

## Features

### 🔐 Authentication
- **Login & Registration**: Secure user authentication with JWT tokens
- **Biometric Authentication**: Support for Face ID/Touch ID (future enhancement)
- **Auto-login**: Persistent authentication using AsyncStorage

### 📊 Dashboard
- **Portfolio Performance**: Real-time portfolio value with interactive charts
- **Live Market Data**: Stock prices, gainers, losers, and trending stocks
- **Position Management**: View and manage your current positions
- **Order History**: Track recent trades and order status
- **Pull-to-refresh**: Refresh data with simple pull gesture

### 💬 AI Assistant
- **Trading Chatbot**: Interactive AI assistant for trading guidance
- **Market Insights**: Get real-time market analysis and recommendations
- **Portfolio Advice**: Personalized investment suggestions

### 📱 Mobile-Optimized Features
- **Touch-friendly Interface**: Optimized for mobile touch interactions
- **Responsive Design**: Adapts to different screen sizes
- **Native Performance**: Built with React Native for optimal performance
- **Offline Support**: Basic functionality works offline (future enhancement)

## Technology Stack

- **React Native**: Cross-platform mobile development
- **React Navigation**: Navigation between screens
- **AsyncStorage**: Local data persistence
- **React Native Chart Kit**: Beautiful charts and graphs
- **NativeWind**: Tailwind CSS for React Native styling
- **TypeScript**: Type-safe development

## Installation

### Prerequisites
- Node.js (v18 or later)
- React Native development environment
- iOS Simulator (for iOS development)
- Android Studio & Emulator (for Android development)

### Setup
1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **iOS Setup** (macOS only):
   ```bash
   cd ios && pod install && cd ..
   ```

3. **Start Metro Bundler**:
   ```bash
   npm start
   ```

4. **Run on iOS**:
   ```bash
   npm run ios
   ```

5. **Run on Android**:
   ```bash
   npm run android
   ```

## Project Structure

```
MobileFrontend/
├── src/
│   ├── components/
│   │   ├── auth/           # Authentication screens
│   │   │   ├── AuthScreen.js
│   │   │   ├── LoginForm.js
│   │   │   └── RegisterForm.js
│   │   ├── dashboard/      # Dashboard components
│   │   │   ├── Dashboard.js
│   │   │   ├── PortfolioPerformance.js
│   │   │   ├── Positions.js
│   │   │   ├── RecentOrders.js
│   │   │   └── MarketData.js
│   │   ├── chatbot/        # AI Assistant
│   │   │   └── Chatbot.js
│   │   └── layout/         # Shared components
│   ├── contexts/           # React contexts
│   │   └── AuthContext.js
│   ├── services/           # API services
│   ├── utils/              # Utility functions
│   └── hooks/              # Custom hooks
├── App.tsx                 # Main app component
└── package.json
```

## API Integration

The mobile app connects to the same backend services as the web frontend:

- **Authentication Service**: `/api/v1/auth/*`
- **Portfolio Service**: `/api/v1/portfolio/*`
- **Trading Service**: `/api/v1/trading/*`
- **Market Data Service**: `/api/v1/market/*`

### Configuration

Update the API base URL in the components:

```javascript
const API_BASE_URL = __DEV__ 
  ? 'http://localhost:3000'          // Development
  : 'https://your-production-api.com'; // Production
```

## Design System

The mobile app follows the same design principles as the web frontend:

- **Colors**: Purple primary (`#7c3aed`), with green/red for gains/losses
- **Typography**: Clear hierarchy with bold headings and readable body text
- **Spacing**: Consistent 16px base spacing with 8px increments
- **Components**: Cards with subtle shadows and rounded corners
- **Interactions**: Smooth animations and haptic feedback

## Development

### Running Tests
```bash
npm test
```

### Linting
```bash
npm run lint
```

### Building for Production

**iOS**:
```bash
cd ios
xcodebuild -workspace MobileFrontend.xcworkspace -scheme MobileFrontend archive
```

**Android**:
```bash
cd android
./gradlew assembleRelease
```

## Features Roadmap

### Short Term
- [ ] Order placement functionality
- [ ] Real-time price updates via WebSocket
- [ ] Push notifications for price alerts
- [ ] Biometric authentication

### Medium Term
- [ ] Advanced charting with technical indicators
- [ ] Options trading interface
- [ ] Social trading features
- [ ] Watchlist management

### Long Term
- [ ] Paper trading mode
- [ ] Educational content integration
- [ ] Advanced portfolio analytics
- [ ] Widget support for home screen

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -am 'Add new feature'`
4. Push to branch: `git push origin feature/new-feature`
5. Submit a pull request

## Support

For issues and feature requests, please use the GitHub issues tracker.

## License

This project is part of the larger Trading App ecosystem. See the main project for licensing information.
