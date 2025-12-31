import axios from 'axios';

const API_URL = 'http://localhost:8083'; // Notification Service URL

export const getNotifications = async (userId) => {
    try {
        const response = await axios.get(`${API_URL}/notifications/${userId}`);
        return response.data;
    } catch (error) {
        console.error('Error fetching notifications:', error);
    }
};

export const markAsRead = async (notificationId) => {
    try {
        const response = await axios.put(`${API_URL}/notifications/${notificationId}`);
        return response.data;
    } catch (error) {
        console.error('Error marking notification as read:', error);
    }
};
