import React, { useState, useEffect } from 'react';
import api from '../api';
import websocketService from '../services/websocketService';
import Notifications from '../components/Notifications';

const AdminDashboard = () => {
    const [tasks, setTasks] = useState([]);
    const [users, setUsers] = useState([]);
    const [title, setTitle] = useState('');
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [activeTab, setActiveTab] = useState('tasks');
    const [currentUserId, setCurrentUserId] = useState(null); // To store logged in admin ID for notifications

    const [username, setUsername] = useState('');
    const [theme, setTheme] = useState('light');

    // ... existing filters state ...
    const [filterStatus, setFilterStatus] = useState('');
    const [filterUser, setFilterUser] = useState('');

    useEffect(() => {
        if (currentUserId) {
            websocketService.connect(currentUserId);
        }
        return () => {
            websocketService.disconnect();
        };
    }, [currentUserId]);

    useEffect(() => {
        fetchTasks();
        fetchUsers();
        // Decode token to get current user ID
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const base64Url = token.split('.')[1];
                const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
                const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
                    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
                }).join(''));
                const decoded = JSON.parse(jsonPayload);
                setCurrentUserId(decoded.user_id);
                fetchUserProfile(decoded.user_id);
            } catch (error) {
                console.error("Error decoding token", error);
            }
        }
    }, [filterStatus, filterUser]); // Refetch when filters change

    const fetchUserProfile = async (uid) => {
        try {
            // Reusing the same strategy as User Dashboard
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

    const applyTheme = (t) => {
        if (t === 'dark') {
            document.body.style.backgroundColor = '#0f172a';
            document.body.style.color = '#f8fafc';
            document.body.classList.add('dark-mode');
        } else {
            document.body.style.backgroundColor = '#f4f7fa';
            document.body.style.color = '#1e293b';
            document.body.classList.remove('dark-mode');
        }
    };

    const handleThemeChange = (newTheme) => {
        setTheme(newTheme);
        applyTheme(newTheme);
    };

    const fetchTasks = async () => {
        try {
            const params = {};
            if (filterStatus) params.status = filterStatus;
            if (filterUser) params.assignedTo = filterUser;

            const res = await api.getTasks(params);
            setTasks(res.data || []);
        } catch (err) {
            console.error(err);
        }
    };

    const fetchUsers = async () => {
        try {
            const res = await api.getUsers();
            setUsers(res.data);
        } catch (err) {
            console.error(err);
        }
    };

    const handleCreateTask = async (e) => {
        e.preventDefault();
        try {
            await api.createTask({ title, assignedTo: selectedUsers });
            fetchTasks();
            setTitle('');
            setSelectedUsers([]);
            alert('Task created!');
        } catch (err) {
            alert('Failed to create task');
        }
    };

    const getStatusBadge = (status) => {
        const s = status ? status.toLowerCase() : 'open';
        return <span className={`badge badge-${s}`}>{status}</span>;
    };

    return (
        <div>
            <nav className="navbar">
                <div className="logo">TaskBoard <span style={{ fontSize: '0.8rem', color: 'var(--text-secondary)', fontWeight: 'normal' }}>Admin</span></div>
                <div className="flex">
                    <button
                        onClick={() => setActiveTab('tasks')}
                        style={{ background: activeTab === 'tasks' ? 'var(--accent-color)' : 'transparent' }}
                    >Tasks</button>
                    <button
                        onClick={() => setActiveTab('users')}
                        style={{ background: activeTab === 'users' ? 'var(--accent-color)' : 'transparent' }}
                    >Users</button>
                    <button
                        onClick={() => setActiveTab('settings')}
                        style={{ background: activeTab === 'settings' ? 'var(--accent-color)' : 'transparent' }}
                    >Settings</button>
                    <Notifications userId={currentUserId} />
                    <button onClick={() => { localStorage.clear(); window.location.href = '/login'; }} style={{ background: 'transparent', border: '1px solid var(--border-color)' }}>Logout</button>
                </div>
            </nav>

            <div className="container">
                {activeTab === 'tasks' && (
                    <>
                        <div className="card">
                            <h3>Create New Task</h3>
                            <form onSubmit={handleCreateTask} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                                <input
                                    value={title}
                                    onChange={(e) => setTitle(e.target.value)}
                                    placeholder="What needs to be done?"
                                    required
                                />

                                <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
                                    <select
                                        onChange={(e) => {
                                            if (e.target.value && !selectedUsers.includes(e.target.value)) {
                                                setSelectedUsers([...selectedUsers, e.target.value]);
                                            }
                                        }}
                                        style={{ flex: 1 }}
                                    >
                                        <option value="">Assign to User...</option>
                                        {users.map(u => (
                                            <option key={u.id} value={u.id}>
                                                {u.username && u.username.trim() !== "" ? u.username : u.email} ({u.role})
                                            </option>
                                        ))}
                                    </select>
                                </div>

                                {selectedUsers.length > 0 && (
                                    <div className="flex" style={{ flexWrap: 'wrap' }}>
                                        <span style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>Assigned:</span>
                                        {selectedUsers.map(id => {
                                            const u = users.find(user => user.id === id);
                                            return (
                                                <span key={id} className="badge" style={{ background: 'rgba(255,255,255,0.1)', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                                    {u ? (u.username && u.username.trim() !== "" ? u.username : u.email) : id}
                                                    <span
                                                        onClick={() => setSelectedUsers(selectedUsers.filter(uid => uid !== id))}
                                                        style={{ cursor: 'pointer', fontWeight: 'bold' }}
                                                    >×</span>
                                                </span>
                                            );
                                        })}
                                    </div>
                                )}

                                <button type="submit" className="btn btn-primary" style={{ alignSelf: 'flex-start' }}>Create Task</button>
                            </form>
                        </div>

                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
                            <h3>Task Board</h3>
                            <div className="flex">
                                <select value={filterStatus} onChange={(e) => setFilterStatus(e.target.value)} style={{ width: '150px' }}>
                                    <option value="">All Statuses</option>
                                    <option value="OPEN">Open</option>
                                    <option value="IN_PROGRESS">In Progress</option>
                                    <option value="DONE">Done</option>
                                </select>
                                <select value={filterUser} onChange={(e) => setFilterUser(e.target.value)} style={{ width: '200px' }}>
                                    <option value="">All Users</option>
                                    {users.map(u => (
                                        <option key={u.id} value={u.id}>{u.email}</option>
                                    ))}
                                </select>
                            </div>
                        </div>

                        {/* Summary Stats */}
                        <div className="stat-grid">
                            <div className="stat-card">
                                <span style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>👥</span>
                                <div className="stat-value">{users.length}</div>
                                <div className="stat-label">Total Users</div>
                            </div>
                            <div className="stat-card">
                                <span style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>📝</span>
                                <div className="stat-value">{tasks.filter(t => (t.status || 'OPEN') === 'OPEN').length}</div>
                                <div className="stat-label">Open Tasks</div>
                            </div>
                            <div className="stat-card">
                                <span style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>🚧</span>
                                <div className="stat-value">{tasks.filter(t => t.status === 'IN_PROGRESS').length}</div>
                                <div className="stat-label">In Progress</div>
                            </div>
                            <div className="stat-card">
                                <span style={{ fontSize: '1.5rem', marginBottom: '0.5rem' }}>✅</span>
                                <div className="stat-value">{tasks.filter(t => t.status === 'DONE').length}</div>
                                <div className="stat-label">Completed</div>
                            </div>
                        </div>

                        <div className="kanban-board">
                            {['OPEN', 'IN_PROGRESS', 'DONE'].map(status => {
                                const columnTasks = tasks.filter(t => (t.status || 'OPEN') === status);
                                const statusLabel = status === 'IN_PROGRESS' ? 'In Progress' : status.charAt(0) + status.slice(1).toLowerCase();

                                return (
                                    <div key={status} className="kanban-column">
                                        <div className="kanban-header">
                                            <span>{statusLabel}</span>
                                            <span className="badge" style={{ background: 'rgba(255,255,255,0.1)', fontSize: '0.7rem' }}>{columnTasks.length}</span>
                                        </div>

                                        {columnTasks.map(task => (
                                            <div key={task.id} className="kanban-task-card">
                                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', marginBottom: '0.5rem' }}>
                                                    <h4 style={{ margin: 0, fontSize: '1rem' }}>{task.title}</h4>
                                                </div>

                                                <div style={{ fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
                                                    <div style={{ marginBottom: '0.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                                                        <span>Assigned to:</span>
                                                        <div className="flex" style={{ flexWrap: 'wrap', gap: '0.25rem' }}>
                                                            {(task.assignedTo || []).length > 0 ? (task.assignedTo || []).map(id => {
                                                                const u = users.find(user => user.id === id);
                                                                return (
                                                                    <span key={id} style={{
                                                                        background: 'rgba(255,255,255,0.1)',
                                                                        padding: '2px 6px',
                                                                        borderRadius: '4px',
                                                                        fontSize: '0.75rem',
                                                                        color: 'var(--text-primary)'
                                                                    }}>
                                                                        {u ? (u.username || u.email.split('@')[0]) : 'User'}
                                                                    </span>
                                                                );
                                                            }) : <span style={{ fontStyle: 'italic' }}>Unassigned</span>}
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        ))}

                                        {columnTasks.length === 0 && (
                                            <div style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-secondary)', fontStyle: 'italic', fontSize: '0.9rem' }}>
                                                No tasks
                                            </div>
                                        )}
                                    </div>
                                );
                            })}
                        </div>
                    </>
                )}

                {activeTab === 'settings' && (
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
                                    await api.updateUserProfile(currentUserId, { username, theme });
                                    setActiveTab('tasks');
                                } catch (e) {
                                    // Silent fail
                                }
                            }}>
                                Save Changes
                            </button>
                        </div>
                    </div>
                )}

                {activeTab === 'users' && (
                    <div className="card">
                        <h3>User Management</h3>
                        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                            <thead>
                                <tr style={{ textAlign: 'left', borderBottom: '1px solid var(--border-color)' }}>
                                    <th style={{ padding: '1rem' }}>Email</th>
                                    <th style={{ padding: '1rem' }}>Role</th>
                                    <th style={{ padding: '1rem' }}>ID</th>
                                    <th style={{ padding: '1rem' }}>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {users.map(u => (
                                    <tr key={u.id} style={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                                        <td style={{ padding: '1rem' }}>{u.email}</td>
                                        <td style={{ padding: '1rem' }}><span className="badge" style={{ background: u.role === 'admin' ? 'var(--accent-color)' : 'rgba(255,255,255,0.1)' }}>{u.role}</span></td>
                                        <td style={{ padding: '1rem', fontFamily: 'monospace', fontSize: '0.8rem' }}>{u.id}</td>
                                        <td style={{ padding: '1rem' }}>
                                            {u.role !== 'admin' && (
                                                <button
                                                    style={{ padding: '0.5rem 1rem', fontSize: '0.8rem' }}
                                                    onClick={async () => {
                                                        if (window.confirm(`Promote ${u.email} to Admin?`)) {
                                                            try {
                                                                await api.updateUserRole(u.id, 'admin');
                                                                alert('User promoted!');
                                                                fetchUsers();
                                                            } catch (err) {
                                                                alert('Failed to promote user');
                                                            }
                                                        }
                                                    }}
                                                >Make Admin</button>
                                            )}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
};

export default AdminDashboard;
