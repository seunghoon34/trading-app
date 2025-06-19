import React, { useState, useRef, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';

const Chatbot = () => {
  const { isAuthenticated } = useAuth();
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState([
    {
      role: 'assistant',
      content: "Hello! I'm your investment AI assistant Zeus. I can help you with portfolio analysis, market insights, investment research, and answer questions about your holdings. What would you like to explore today?"
    }
  ]);
  const [loading, setLoading] = useState(false);

  // Add refs for auto-scroll
  const messagesEndRef = useRef(null);
  const messagesContainerRef = useRef(null);

  // API base URL
  const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  const handleSend = async () => {
    if (!message.trim() || loading || !isAuthenticated) return;
    
    const userMessage = { role: 'user', content: message };
    const newMessages = [...messages, userMessage];
    setMessages(newMessages);
    setMessage('');
    setLoading(true);

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/zeus/chat`, {
        method: 'POST',
        credentials: 'include', // Include cookies for authentication
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          messages: newMessages,
          url: 'https://72ec-2001-fb1-41-a82a-2dae-45ee-bdf5-94b5.ngrok-free.app'
        })
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      setMessages([...newMessages, { role: 'assistant', content: data.response }]);
    } catch (error) {
      console.error('Error:', error);
      setMessages([...newMessages, { 
        role: 'assistant', 
        content: 'Sorry, I encountered an error while processing your request. Please try again.' 
      }]);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  // Show authentication required message if not authenticated
  if (!isAuthenticated) {
    return (
      <div className="h-full flex flex-col max-h-screen items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">Authentication Required</h2>
          <p className="text-gray-600">Please log in to use the Zeus AI assistant.</p>
        </div>
      </div>
    );
  }

      return (
      <div className="h-full flex flex-col max-h-screen">

      {/* Messages - This will be scrollable */}
      <div 
        className="flex-1 overflow-y-auto px-8 py-5 flex justify-center"
        ref={messagesContainerRef}
      >
        <div className="w-full max-w-4xl space-y-6">
          {messages.map((msg, index) => (
            <div key={index} className={msg.role === 'user' ? 'flex justify-end' : ''}>
              {msg.role === 'assistant' ? (
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
            {loading && (
            <div className="prose max-w-none">
              <div className="whitespace-pre-line text-gray-800 leading-relaxed">
                <div className="flex items-center space-x-2">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-900"></div>
                  <span>Thinking...</span>
                </div>
              </div>
            </div>
          )}
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
                disabled={loading}
                style={{
                  lineHeight: '1.5',
                  fontSize: '16px'
                }}
              />
              
              {/* Send Button */}
              <button
                onClick={handleSend}
                disabled={!message.trim() || loading}
                className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center transition-all ${
                  message.trim() && !loading
                    ? 'bg-black text-white hover:bg-gray-800'
                    : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                }`}
              >
                {loading ? (
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                ) : (
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="m22 2-7 20-4-9-9-4z"/>
                    <path d="M22 2 11 13"/>
                  </svg>
                )}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Chatbot; 