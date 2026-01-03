import axios from 'axios';

const API_URLS = {
    AUTH: 'http://localhost:8080',
    TASK: 'http://localhost:8081',
    USER: 'http://localhost:8087',
    NOTIFICATION: 'http://localhost:8083'
};

const api = {
    login: (email, password) => axios.post(`${API_URLS.AUTH}/login`, { email, password }),
    login: (email, password) => axios.post(`${API_URLS.AUTH}/login`, { email, password }),
    register: (email, password, role, username) => axios.post(`${API_URLS.AUTH}/register`, { email, password, role, username }),

    // Auth Service
    updateUserRole: (userId, role) => axios.put(`${API_URLS.AUTH}/users/${userId}/role`, { role }, {
        headers: { Authorization: `Bearer ${localStorage.getItem('token')}` }
    }),

    // Task Service
    createTask: (taskData) => axios.post(`${API_URLS.TASK}/tasks`, taskData),
    getTasks: (params) => axios.get(`${API_URLS.TASK}/tasks`, { params }),
    updateTaskStatus: (taskId, status, updatedBy) => axios.put(`${API_URLS.TASK}/tasks/${taskId}/status`, { status, updatedBy }),

    // User Service
    getUsers: () => axios.get(`${API_URLS.USER}/users`),
    updateUserProfile: (userId, data) => axios.put(`${API_URLS.USER}/users/${userId}`, data),

    // Notification Service
    getNotifications: (userId) => axios.get(`${API_URLS.NOTIFICATION}/notifications/${userId}`),
    markNotificationAsRead: (id) => axios.put(`${API_URLS.NOTIFICATION}/notifications/${id}`),
};

export default api;
