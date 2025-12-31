import React, { useState } from 'react';
import Auth from './components/auth';
import TaskList from './components/tasklist';
import Notification from './components/Notification';
import { loginUser } from './services/authService';

const App = () => {
  const [user, setUser] = useState(null);

  const handleLogin = async (email, password) => {
    const userData = await loginUser(email, password);
    if (userData) {
      setUser(userData);
    }
  };

  return (
    <div>
      {!user ? (
        <Auth setUser={handleLogin} />
      ) : (
        <div>
          <h1>Welcome, {user.email}</h1>
          <TaskList />
          <Notification userId={user.user_id} />
        </div>
      )}
    </div>
  );
};

export default App;
