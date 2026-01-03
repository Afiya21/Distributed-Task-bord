import React, { useState, useEffect } from 'react';
import { jwtDecode } from 'jwt-decode';
import api from '../api';
import websocketService from '../services/websocketService';
import Notifications from '../components/Notifications';

const UserDashboard = () => {
    const [tasks, setTasks] = useState([]);
    const [loading, setLoading] = useState(false);
    const [userId, setUserId] = useState('');
    const [tab, setTab] = useState('tasks');
    const [username, setUsername] = useState('');

    // Theme state
    const [theme, setTheme] = useState('light');

    useEffect(() => {
        // Decode token to get user ID
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decoded = jwtDecode(token);
                setUserId(decoded.user_id);
                fetchTasks(decoded.user_id);
                fetchUserProfile(decoded.user_id); // Fetch profile data (theme/name)
                // Connect WebSocket
                websocketService.connect(decoded.user_id);
            } catch (e) {
                console.error("Invalid token", e);
            }
        }
        return () => {
            websocketService.disconnect();
        };
    }, []);

    const fetchUserProfile = async (uid) => {
        try {
            // Need a new endpoint to get specific user profile, or loop from GetAllUsers (inefficient) or use Auth data?
            // User Service usually has GetUserByID. Let's assume api.getUser(uid) exists or similar.
            // If not, we might need to rely on what we have.
            // Wait, we can use api.getUsers() and filter? Or simpler: 
            // The user service likely supports GET /users/:id. Let's check api.js later.
            // For now, let's assume we can fetch it. If not, we'll default to light.

            // Actually, let's add a proper fetch here if the API supports it. 
            // Checking api.js... existing one is getUsers -> list.
            // Let's rely on standard practice: we need current user data.
            // We can just rely on the user manually setting it for now if API is missing, 
            // BUT to fix "first thing not working", we must fetch it.
            // Let's implement fetchUser(uid).

            // Temporarily using getUsers and filtering (not ideal but works for small scale)
            const res = await api.getUsers();
            const me = res.data.find(u => u.id === uid);
            if (me) {
                setUsername(me.username || '');
                if (me.theme) {
                    setTheme(me.theme);
                    applyTheme(me.theme);
                }
            }
        } catch (err) {
            console.error("Failed to fetch profile", err);
        }
    };

    const fetchTasks = async (uid) => {
        try {
            // Filter by assignedTo = current user ID
            const res = await api.getTasks({ assignedTo: uid });
            setTasks(res.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    const applyTheme = (t) => {
        if (t === 'dark') {
            document.body.style.backgroundColor = '#0f172a'; // Match CSS var
            document.body.style.color = '#f8fafc';
            document.body.classList.add('dark-mode');
        } else {
            document.body.style.backgroundColor = '#f4f7fa'; // Light gray
            document.body.style.color = '#1e293b'; // Slate 800
            document.body.classList.remove('dark-mode');
        }
    };

    const handleStatusUpdate = async (taskId, newStatus) => {
        setLoading(true);
        try {
            await api.updateTaskStatus(taskId, newStatus, userId);
            fetchTasks(userId); // Refresh to see update
        } catch (err) {
            alert('Failed to update status');
        } finally {
            setLoading(false);
        }
    };

    const getStatusBadge = (status) => {
        const s = status ? status.toLowerCase() : 'open';
        return <span className={`badge badge-${s}`}>{status}</span>;
    };

    const handleThemeChange = (newTheme) => {
        setTheme(newTheme);
        applyTheme(newTheme); // Apply visually immediately (Preview)
    };

    return (
        <div>
            <nav className="navbar">
                <div className="logo">TaskBoard <span style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', fontWeight: 'normal' }}>User</span></div>
                <div className="flex">
                    <button
                        onClick={() => setTab('tasks')}
                        style={{ background: tab === 'tasks' ? 'var(--accent-color)' : 'transparent' }}
                    >Tasks</button>
                    <button
                        onClick={() => setTab('settings')}
                        style={{ background: tab === 'settings' ? 'var(--accent-color)' : 'transparent' }}
                    >Settings</button>
                    <Notifications userId={userId} />
                    <button onClick={() => { localStorage.clear(); window.location.href = '/login'; }} style={{ background: 'transparent', border: '1px solid var(--border-color)' }}>Logout</button>
                </div>
            </nav>

            <div className="container">
                {tab === 'tasks' && (
                    <>
                        <h3>My Tasks</h3>
                        {tasks.length === 0 ? (
                            <p style={{ color: 'var(--text-secondary)' }}>No tasks assigned to you yet.</p>
                        ) : (
                            <div className="task-list">
                                {tasks.map(task => (
                                    <div key={task.id} className="task-card">
                                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '1rem' }}>
                                            <h4 style={{ margin: 0 }}>{task.title}</h4>
                                            {getStatusBadge(task.status)}
                                        </div>

                                        <hr style={{ borderColor: 'var(--border-color)', margin: '1rem 0' }} />

                                        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '0.5rem', marginTop: '1rem' }}>
                                            <button
                                                className={task.status === "OPEN" ? "" : "secondary"}
                                                onClick={() => handleStatusUpdate(task.id, 'OPEN')}
                                                style={{ opacity: task.status === 'OPEN' ? 1 : 0.5, fontSize: '0.8rem', padding: '0.5rem' }}
                                                disabled={loading}
                                            >
                                                To Open
                                            </button>
                                            <button
                                                className={task.status === "IN_PROGRESS" ? "" : "secondary"}
                                                onClick={() => handleStatusUpdate(task.id, 'IN_PROGRESS')}
                                                style={{ opacity: task.status === 'IN_PROGRESS' ? 1 : 0.5, fontSize: '0.8rem', padding: '0.5rem' }}
                                                disabled={loading}
                                            >
                                                Progress
                                            </button>
                                            <button
                                                className={task.status === "DONE" ? "" : "secondary"}
                                                onClick={() => handleStatusUpdate(task.id, 'DONE')}
                                                style={{ opacity: task.status === 'DONE' ? 1 : 0.5, fontSize: '0.8rem', padding: '0.5rem' }}
                                                disabled={loading}
                                            >
                                                Done
                                            </button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </>
                )}

                {tab === 'settings' && (
                    <div className="card" style={{ maxWidth: '600px', margin: '0 auto' }}>
                        <h3>Profile Settings</h3>

                        <div className="form-group">
                            <label>Display Name</label>
                            <input
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                placeholder="Enter your display name"
                            />
                        </div>

                        <div className="form-group">
                            <label>Appearance</label>
                            <div className="settings-grid">
                                <div
                                    className={`theme-card ${theme === 'light' ? 'active' : ''}`}
                                    onClick={() => handleThemeChange('light')}
                                >
                                    <span className="theme-icon">☀️</span>
                                    <h4>Light Mode</h4>
                                </div>
                                <div
                                    className={`theme-card ${theme === 'dark' ? 'active' : ''}`}
                                    onClick={() => handleThemeChange('dark')}
                                >
                                    <span className="theme-icon">🌙</span>
                                    <h4>Dark Mode</h4>
                                </div>
                            </div>
                        </div>

                        <div style={{ marginTop: '2rem', display: 'flex', justifyContent: 'flex-end' }}>
                            <button className="btn btn-primary" onClick={async () => {
                                try {
                                    await api.updateUserProfile(userId, { username, theme });
                                    setTab('tasks'); // Redirect to main page
                                } catch (e) {
                                    // alert('Failed to update profile'); 
                                }
                            }}>
                                Save Changes
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default UserDashboard;
