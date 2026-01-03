import React, { useState, useEffect, useRef } from 'react';
import api from '../api';
import websocketService from '../services/websocketService';

const Notifications = ({ userId }) => {
    const [notifications, setNotifications] = useState([]);
    const [isOpen, setIsOpen] = useState(false);

    // Ref for the dropdown container
    const dropdownRef = useRef(null);

    // Close dropdown on click outside
    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
                setIsOpen(false);
            }
        };

        if (isOpen) {
            document.addEventListener("mousedown", handleClickOutside);
        } else {
            document.removeEventListener("mousedown", handleClickOutside);
        }

        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, [isOpen]);

    useEffect(() => {
        if (userId) {
            fetchNotifications();
        }
    }, [userId]);

    // Real-time listener
    useEffect(() => {
        const handleNewNotification = (notification) => {
            console.log("New notification received:", notification);
            setNotifications(prev => [notification, ...prev]);
        };

        websocketService.registerCallback(handleNewNotification);

        return () => {
            websocketService.unregisterCallback(handleNewNotification);
        };
    }, []);

    const fetchNotifications = async () => {
        try {
            const res = await api.getNotifications(userId);
            // Sort by newest first
            const sorted = (res.data || []).sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
            setNotifications(sorted);
        } catch (err) {
            console.error("Failed to fetch notifications", err);
        }
    };

    const handleMarkAsRead = async (id) => {
        try {
            await api.markNotificationAsRead(id);
            setNotifications(notifications.map(n => n.id === id ? { ...n, is_read: true } : n));
        } catch (err) {
            console.error("Failed to mark as read", err);
        }
    };

    const unreadCount = notifications.filter(n => n && !n.is_read).length;

    return (
        <div style={{ position: 'relative' }} ref={dropdownRef}>
            <button
                onClick={() => setIsOpen(!isOpen)}
                style={{ background: 'transparent', border: 'none', color: 'var(--text-primary)', cursor: 'pointer', position: 'relative' }}
            >
                🔔
                {unreadCount > 0 && (
                    <span style={{
                        position: 'absolute',
                        top: '-5px',
                        right: '-5px',
                        background: 'red',
                        color: 'white',
                        borderRadius: '50%',
                        padding: '2px 6px',
                        fontSize: '0.7rem'
                    }}>{unreadCount}</span>
                )}
            </button>

            {isOpen && (
                <div style={{
                    position: 'absolute',
                    top: '30px',
                    right: '0',
                    width: '300px',
                    background: '#1e1e1e',
                    border: '1px solid var(--border-color)',
                    borderRadius: '8px',
                    boxShadow: '0 4px 6px rgba(0,0,0,0.3)',
                    zIndex: 1000,
                    maxHeight: '400px',
                    overflowY: 'auto'
                }}>
                    <div style={{ padding: '0.5rem 1rem', borderBottom: '1px solid var(--border-color)', fontWeight: 'bold' }}>
                        Notifications
                        <button onClick={fetchNotifications} style={{ float: 'right', background: 'transparent', border: 'none', cursor: 'pointer', fontSize: '1rem' }}>🔄</button>
                    </div>
                    {notifications.length === 0 ? (
                        <div style={{ padding: '1rem', color: 'var(--text-secondary)', textAlign: 'center' }}>No notifications</div>
                    ) : (
                        notifications.map(n => n && (
                            <div key={n.id} style={{
                                padding: '0.8rem',
                                borderBottom: '1px solid rgba(255,255,255,0.05)',
                                background: n.is_read ? 'transparent' : 'rgba(255,255,255,0.05)',
                                opacity: n.is_read ? 0.6 : 1
                            }}>
                                <p style={{ margin: '0 0 0.5rem 0', fontSize: '0.9rem' }}>{n.message}</p>
                                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                    <small style={{ color: 'var(--text-secondary)' }}>{new Date(n.created_at).toLocaleString()}</small>
                                    {!n.is_read && (
                                        <button
                                            onClick={() => handleMarkAsRead(n.id)}
                                            style={{ padding: '2px 8px', fontSize: '0.7rem' }}
                                        >Mark Read</button>
                                    )}
                                </div>
                            </div>
                        ))
                    )}
                </div>
            )}
        </div>
    );
};

export default Notifications;
