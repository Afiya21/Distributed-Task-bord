import axios from 'axios';

const API_URL = 'http://localhost:8081'; // Task Service URL

export const createTask = async (task) => {
    try {
        const response = await axios.post(`${API_URL}/tasks`, task, {
            headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        });
        return response.data;
    } catch (error) {
        console.error('Error creating task:', error);
    }
};

export const getTasks = async () => {
    try {
        const response = await axios.get(`${API_URL}/tasks`, {
            headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        });
        return response.data;
    } catch (error) {
        console.error('Error fetching tasks:', error);
    }
};

export const updateTask = async (taskId, updatedTask) => {
    try {
        const response = await axios.put(`${API_URL}/tasks/${taskId}`, updatedTask, {
            headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        });
        return response.data;
    } catch (error) {
        console.error('Error updating task:', error);
    }
};

export const deleteTask = async (taskId) => {
    try {
        const response = await axios.delete(`${API_URL}/tasks/${taskId}`, {
            headers: { Authorization: `Bearer ${localStorage.getItem('token')}` },
        });
        return response.data;
    } catch (error) {
        console.error('Error deleting task:', error);
    }
};
