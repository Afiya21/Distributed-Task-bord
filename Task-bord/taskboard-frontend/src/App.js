import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import AdminDashboard from './pages/AdminDashboard';
import UserDashboard from './pages/UserDashboard';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        {/* Navigation removed for cleaner auth flow, or keep consistent? */}
        {/* Let's keep it minimal or just remove generic links since we redirect */}
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/admin" element={<AdminDashboard />} />
          <Route path="/user" element={<UserDashboard />} />
          <Route path="/" element={
            <div className="auth-container">
              <div className="card" style={{ textAlign: 'center', maxWidth: '500px' }}>
                <h1 className="logo" style={{ fontSize: '3rem', marginBottom: '1rem' }}>TaskBoard</h1>
                <p style={{ color: 'var(--text-secondary)', marginBottom: '2rem', lineHeight: '1.6' }}>
                  Manage your tasks efficiently with our distributed task management system.
                  Assign tasks, track progress, and stay organized.
                </p>
                <div className="flex" style={{ justifyContent: 'center', gap: '1rem' }}>
                  <Link to="/login">
                    <button>Login</button>
                  </Link>
                  <Link to="/register">
                    <button style={{ background: 'transparent', border: '1px solid var(--border-color)' }}>Register</button>
                  </Link>
                </div>
              </div>
            </div>
          } />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
