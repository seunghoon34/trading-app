import React, { useState, useRef, useEffect } from 'react';

const Chatbot = () => {
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState([
    {
      type: 'assistant',
      content: "Hello! I'm your investment AI assistant. I can help you with portfolio analysis, market insights, investment research, and answer questions about your holdings. What would you like to explore today?"
    },
    {
      type: 'user',
      content: "Can you analyze my current portfolio performance and suggest any optimizations?"
    },
    {
      type: 'assistant',
      content: "I'd be happy to analyze your portfolio! Based on your current holdings, I can see you have a well-diversified mix across tech, energy, and traditional markets. Here are some key observations:\n\n**Performance Summary:**\n• Your portfolio is up 12.4% overall, outperforming the S&P 500 by 3.2%\n• AAPL and MSFT are your strongest performers\n• TSLA position is currently underperforming\n\n**Optimization Suggestions:**\n1. Consider rebalancing your tech allocation (currently 65% of portfolio)\n2. Your portfolio lacks exposure to international markets\n3. Consider adding some defensive positions given current market volatility\n\nWould you like me to dive deeper into any of these areas or run a specific analysis?"
    }
  ]);

  // Add refs for auto-scroll
  const messagesEndRef = useRef(null);
  const messagesContainerRef = useRef(null);

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  const handleSend = () => {
    if (message.trim()) {
      setMessages([...messages, { type: 'user', content: message }]);
      setMessage('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="h-full flex flex-col max-h-screen">
      {/* Messages - This will be scrollable */}
      <div 
        className="flex-1 overflow-y-auto px-8 py-5 flex justify-center"
        ref={messagesContainerRef}
      >
        <div className="w-full max-w-4xl space-y-6">
          {messages.map((msg, index) => (
            <div key={index} className={msg.type === 'user' ? 'flex justify-end' : ''}>
              {msg.type === 'assistant' ? (
                <div className="prose max-w-none">
                  <div className="whitespace-pre-line text-gray-800 leading-relaxed">
                    {msg.content}
                  </div>
                </div>
              ) : (
                <div className="bg-gray-100 rounded-2xl px-4 py-3 max-w-xl">
                  <div className="text-gray-800">
                    {msg.content}
                  </div>
                </div>
              )}
            </div>
          ))}
          {/* Invisible div at the bottom to scroll to */}
          <div ref={messagesEndRef} />
        </div>
      </div>
        
      {/* Input Area - Floating like Claude */}
      <div className="bg-gray-50 px-8 py-6">
        <div className="flex justify-center">
          <div className="w-full max-w-3xl">
            {/* Floating Input Bubble */}
            <div className="relative bg-white border border-gray-300 rounded-3xl shadow-lg px-4 py-3 flex items-end gap-3">
              <textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Ask me about investments, portfolio analysis, market trends..."
                className="flex-1 resize-none outline-none bg-transparent min-h-6 max-h-32 placeholder-gray-500"
                rows={1}
                style={{
                  lineHeight: '1.5',
                  fontSize: '16px'
                }}
              />
              
              {/* Send Button */}
              <button
                onClick={handleSend}
                disabled={!message.trim()}
                className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center transition-all ${
                  message.trim()
                    ? 'bg-black text-white hover:bg-gray-800'
                    : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                }`}
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="m22 2-7 20-4-9-9-4z"/>
                  <path d="M22 2 11 13"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Chatbot; 