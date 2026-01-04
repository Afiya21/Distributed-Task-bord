import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
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
    const [userEmail, setUserEmail] = useState('');
    const [userRole, setUserRole] = useState('');
    const [theme, setTheme] = useState('dark'); // Default to dark
    const [refreshKey, setRefreshKey] = useState(0); // Added for useEffect dependency
    const navigate = useNavigate();

    useEffect(() => {
        // Decode token to get current user ID
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decoded = jwtDecode(token);
                setUserId(decoded.user_id);
                // Set initial profile from token
                setUsername(decoded.username || '');
                setUserEmail(decoded.email || '');
                setUserRole(decoded.role || 'user');

                fetchTasks(decoded.user_id); // Keep fetching tasks
                fetchUser(decoded.user_id); // Fetch full profile (theme)
                if (websocketService.connect) {
                    websocketService.connect(decoded.user_id);
                }
            } catch (e) {
                console.error("Invalid token", e);
                // Optionally redirect to login if token is invalid
                localStorage.clear();
                navigate('/login');
            }
        } else {
            // If no token, redirect to login
            navigate('/login');
        }
        return () => {
            websocketService.disconnect();
        };
    }, [refreshKey]);

    const fetchUser = async (uid) => {
        try {
            const res = await api.getUsers();
            if (res.data) {
                const me = res.data.find(u => u.id === uid);
                if (me) {
                    setUsername(me.username || '');
                    setUserEmail(me.email || '');
                    setUserRole(me.role || 'user');
                    if (me.theme) {
                        setTheme(me.theme);
                        applyTheme(me.theme);
                    }
                }
            }
        } catch (err) {
            console.error(err);
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
        if (t === 'light') {
            document.body.classList.add('light-mode');
        } else {
            document.body.classList.remove('light-mode');
        }
    };

    const handleStatusUpdate = async (taskId, newStatus) => {
        setLoading(true);
        try {
            await api.updateTaskStatus(taskId, newStatus, userId);
            fetchTasks(userId); // Refresh to see update
        } catch (err) {
            console.error('Failed to update status');
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
                        className={`btn-ghost ${tab === 'tasks' ? 'active' : ''}`}
                    >Tasks</button>
                    <button
                        onClick={() => setTab('settings')}
                        className={`btn-ghost ${tab === 'settings' ? 'active' : ''}`}
                    >Settings</button>
                    <Notifications userId={userId} />
                    <button onClick={() => { localStorage.clear(); window.location.href = '/login'; }} className="btn-ghost" style={{ border: '1px solid var(--border-color)' }}>Logout</button>
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
                            <label>Profile Information</label>
                            <div style={{
                                background: 'rgba(255, 255, 255, 0.05)',
                                padding: '1.5rem',
                                borderRadius: '0.5rem',
                                border: '1px solid var(--border-color)',
                                display: 'flex',
                                flexDirection: 'column',
                                gap: '0.5rem'
                            }}>
                                <div style={{ fontSize: '1.2rem', fontWeight: 'bold' }}>{username}</div>
                                <div style={{ color: 'var(--text-secondary)', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                    <span>📧 {userEmail}</span>
                                    <span className="badge" style={{ background: 'rgba(255, 255, 255, 0.1)', fontSize: '0.7rem' }}>{userRole}</span>
                                </div>
                            </div>
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
