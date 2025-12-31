import React, { useEffect, useState } from 'react';
import { getNotifications, markAsRead } from '../services/notificationService';

const Notification = ({ userId }) => {
    const [notifications, setNotifications] = useState([]);

    useEffect(() => {
        async function fetchNotifications() {
            const fetchedNotifications = await getNotifications(userId);
            setNotifications(fetchedNotifications);
        }
        fetchNotifications();
    }, [userId]);

    return (
        <div>
            <h2>Notifications</h2>
            <ul>
                {notifications.map((notification) => (
                    <li key={notification.id}>
                        <p>{notification.message}</p>
                        <button onClick={() => markAsRead(notification.id)}>Mark as read</button>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Notification;
