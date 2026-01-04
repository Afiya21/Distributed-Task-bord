import React, { useState, useEffect } from 'react';
import api from '../api';
import websocketService from '../services/websocketService';
import Notifications from '../components/Notifications';
import Modal from '../components/Modal';

const AdminDashboard = () => {
    const [tasks, setTasks] = useState([]);
    const [users, setUsers] = useState([]);
    const [title, setTitle] = useState('');
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [currentSelection, setCurrentSelection] = useState('');
    const [activeTab, setActiveTab] = useState('tasks');
    const [showSuccess, setShowSuccess] = useState(false);
    const [error, setError] = useState(null); // Add error state if we want inline errors later
    const [currentUserId, setCurrentUserId] = useState(null); // To store logged in admin ID for notifications

    const [username, setUsername] = useState('');
    const [userEmail, setUserEmail] = useState('');
    const [userRole, setUserRole] = useState('');
    const [theme, setTheme] = useState('dark');

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
                // Set initial profile from token
                setUsername(decoded.username || '');
                setUserEmail(decoded.email || '');
                setUserRole(decoded.role || 'admin');

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
                setUserEmail(me.email || '');
                setUserRole(me.role || 'admin');
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
            document.body.classList.remove('light-mode');
        } else {
            document.body.classList.add('light-mode');
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
            fetchTasks(); // Refresh list
            setTitle(''); // Clear title
            setSelectedUsers([]); // Clear assigned users
            setCurrentSelection(''); // Reset dropdown
            setShowSuccess(true); // Show success modal
            setTimeout(() => setShowSuccess(false), 2000); // Auto close
        } catch (err) {
            console.error('Failed to create task');
        }
    };

    const promoteToAdmin = async (userId) => {
        if (window.confirm('Are you sure you want to promote this user to Admin?')) {
            try {
                await api.updateUserRole(userId, 'admin');
                // alert('User promoted successfully'); // Removed alert
                fetchUsers();
            } catch (err) {
                console.error('Failed to promote user', err); // Replaced alert with console.error
            }
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
                        className={`btn-ghost ${activeTab === 'tasks' ? 'active' : ''}`}
                    >Tasks</button>
                    <button
                        onClick={() => setActiveTab('users')}
                        className={`btn-ghost ${activeTab === 'users' ? 'active' : ''}`}
                    >Users</button>
                    <button
                        onClick={() => setActiveTab('settings')}
                        className={`btn-ghost ${activeTab === 'settings' ? 'active' : ''}`}
                    >Settings</button>
                    <Notifications userId={currentUserId} />
                    <button onClick={() => { localStorage.clear(); window.location.href = '/login'; }} className="btn-ghost" style={{ border: '1px solid var(--border-color)' }}>Logout</button>
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
                                        className="custom-select"
                                        value={currentSelection}
                                        onChange={(e) => {
                                            const val = e.target.value;
                                            setCurrentSelection(val);
                                            if (val && !selectedUsers.includes(val)) {
                                                setSelectedUsers([...selectedUsers, val]);
                                                // Optional: Reset immediately if we want to allow rapid multiple creation? 
                                                // But usually user wants to see what they picked? 
                                                // Actually, if it adds a tag below, the dropdown should probably reset to allow picking another?
                                                // User specific request: "after task is created".
                                                // So I will only reset on submit? 
                                                // If I don't reset here, the dropdown shows the selected user.
                                                // If I reset here (to ""), then it's "blank" (default) immediately.
                                                // The tags show who is selected.
                                                // Let's reset it immediately so they can pick multiple easily.
                                                // setCurrentSelection(''); // User wants it to stay selected until submit
                                            }
                                        }}
                                        style={{ flex: 1 }}
                                    >
                                        <option value="">Assigned To</option>
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
                                    <span className="badge" style={{ background: 'var(--accent-color)', fontSize: '0.7rem' }}>{userRole}</span>
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
                                    await api.updateUserProfile(currentUserId, { username, theme });
                                    // alert('Profile updated successfully!'); // User requested no popup
                                    setActiveTab('tasks');
                                } catch (e) {
                                    console.error("Failed to update profile", e);
                                    // alert('Failed to update profile. Please try again.');
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
                        <table className="table" style={{ width: '100%', borderCollapse: 'collapse', marginTop: '1rem', tableLayout: 'fixed' }}>
                            <thead>
                                <tr style={{ borderBottom: '1px solid var(--border-color)', textAlign: 'left' }}>
                                    <th style={{ padding: '0.75rem', width: '25%' }}>Full Name</th>
                                    <th style={{ padding: '0.75rem', width: '35%' }}>Email</th>
                                    <th style={{ padding: '0.75rem', width: '15%' }}>Role</th>
                                    <th style={{ padding: '0.75rem', width: '25%' }}>Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {users.map(user => (
                                    <tr key={user.id} style={{ borderBottom: '1px solid var(--border-color)', textAlign: 'left' }}>
                                        <td style={{ padding: '0.75rem', width: '25%', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                                            {user.username || 'N/A'}
                                        </td>
                                        <td style={{ padding: '0.75rem', width: '35%', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                                            {user.email}
                                        </td>
                                        <td style={{ padding: '0.75rem', width: '15%' }}>
                                            <span style={{
                                                padding: '0.25rem 0.5rem',
                                                borderRadius: '1rem',
                                                fontSize: '0.85rem',
                                                backgroundColor: user.role === 'admin' ? 'rgba(99, 102, 241, 0.1)' : 'rgba(16, 185, 129, 0.1)',
                                                color: user.role === 'admin' ? '#6366f1' : '#10b981'
                                            }}>
                                                {user.role}
                                            </span>
                                        </td>
                                        <td style={{ padding: '0.75rem', width: '25%' }}>
                                            {user.role !== 'admin' && (
                                                <button
                                                    onClick={() => promoteToAdmin(user.id)}
                                                    style={{
                                                        padding: '0.25rem 0.75rem',
                                                        fontSize: '0.85rem',
                                                        backgroundColor: 'transparent',
                                                        border: '1px solid #6366f1',
                                                        color: '#6366f1',
                                                        borderRadius: '0.25rem',
                                                        cursor: 'pointer'
                                                    }}
                                                >
                                                    Promote to Admin
                                                </button>
                                            )}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
            <Modal
                isOpen={showSuccess}
                title="Task Created!"
                message="The task has been successfully assigned."
                icon="✅"
            />
        </div>
    );
};

export default AdminDashboard;
