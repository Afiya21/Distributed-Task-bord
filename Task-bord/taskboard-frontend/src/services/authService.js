import axios from 'axios';

const API_URL = 'http://localhost:8080'; // Auth Service URL

export const registerUser = async (email, password, role) => {
    try {
        const response = await axios.post(`${API_URL}/register`, {
            email,
            password,
            role,
        });
        return response.data;
    } catch (error) {
        console.error('Error during registration:', error);
    }
};

export const loginUser = async (email, password) => {
    try {
        const response = await axios.post(`${API_URL}/login`, {
            email,
            password,
        });
        localStorage.setItem('token', response.data.token); // Store JWT token in local storage
        return response.data;
    } catch (error) {
        console.error('Error during login:', error);
    }
};

export const logoutUser = () => {
    localStorage.removeItem('token'); // Remove token on logout
};
