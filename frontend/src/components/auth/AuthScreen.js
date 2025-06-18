import React, { useState } from 'react';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';

const AuthScreen = () => {
  const [isLogin, setIsLogin] = useState(true);

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-500 to-blue-500 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        {/* Logo and Title */}
        <div className="p-8 text-center">
          <div className="flex items-center justify-center gap-3 mb-3">
            <div className="w-12 h-12 rounded-xl bg-gradient-to-r from-purple-500 via-blue-500 to-indigo-500 flex items-center justify-center text-white font-bold text-2xl shadow-lg">
              P
            </div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-purple-500 via-blue-500 to-indigo-500 bg-clip-text text-transparent">
              Pandora Wealth
            </h1>
          </div>
          <p className="text-gray-600 mt-2">
            {isLogin ? 'Welcome back!' : 'Create your account'}
          </p>
        </div>

        {/* Form Container */}
        <div className="px-8 pb-8">
          {isLogin ? (
            <LoginForm />
          ) : (
            <RegisterForm />
          )}

          {/* Toggle between login and register */}
          <div className="text-center mt-6">
            <p className="text-gray-600">
              {isLogin ? "Don't have an account?" : "Already have an account?"}
              <button
                onClick={() => setIsLogin(!isLogin)}
                className="ml-2 text-purple-600 hover:text-purple-700 font-medium"
              >
                {isLogin ? 'Sign up' : 'Sign in'}
              </button>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuthScreen; 