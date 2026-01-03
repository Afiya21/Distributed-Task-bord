import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../api';
import Modal from '../components/Modal';

const RegisterPage = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [username, setUsername] = useState('');
    const [role, setRole] = useState('user'); // Default to user
    const [error, setError] = useState('');
    const [showSuccess, setShowSuccess] = useState(false);
    const navigate = useNavigate();

    const validatePassword = (pwd) => {
        // At least 8 chars, 1 uppercase, 1 number, 1 special char
        const regex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/;
        return regex.test(pwd);
    };

    const handleRegister = async (e) => {
        e.preventDefault();
        setError('');

        if (!validatePassword(password)) {
            setError('Password must be at least 8 characters long and include an uppercase letter, a number, and a special character (!@#$%^&*).');
            return;
        }

        try {
            await api.register(email, password, role, username);
            setShowSuccess(true);
            setTimeout(() => {
                navigate('/login');
            }, 2000);
        } catch (err) {
            setError(err.response?.data?.error || 'Registration failed');
        }
    };

    return (
        <div className="auth-container">
            <Modal
                isOpen={showSuccess}
                title="Welcome aboard!"
                message="Your account has been created successfully. Redirecting to login..."
            />

            <div className="card auth-form">
                <center>
                    <h2 className="logo">TaskBoard</h2>
                    <p style={{ color: 'var(--text-secondary)' }}>Create your account</p>
                </center>

                {error && (
                    <div style={{
                        backgroundColor: 'rgba(239, 68, 68, 0.2)',
                        color: '#f87171',
                        padding: '0.75rem',
                        borderRadius: '0.5rem',
                        marginBottom: '1rem',
                        textAlign: 'center',
                        fontSize: '0.9rem'
                    }}>
                        {error}
                    </div>
                )}

                <form onSubmit={handleRegister} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <div>
                        <input
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            placeholder="Full Name"
                            required
                        />
                    </div>
                    <div>
                        <input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            placeholder="Email Address"
                            required
                        />
                    </div>
                    <div>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Password"
                            required
                        />
                        <p style={{ fontSize: '0.75rem', color: 'var(--text-secondary)', marginTop: '0.25rem', paddingLeft: '0.25rem' }}>
                            Must be 8+ chars with uppercase, number & special char.
                        </p>
                    </div>
                    {/* Role selection removed: Registration is only for Users */}

                    <button type="submit" className="btn btn-primary" style={{ justifyContent: 'center' }}>Create Account</button>

                    <p style={{ textAlign: 'center', fontSize: '0.9rem', color: 'var(--text-secondary)', marginTop: '1rem' }}>
                        Already have an account? <Link to="/login" style={{ color: 'var(--accent-color)', textDecoration: 'none' }}>Sign In</Link>
                    </p>
                </form>
            </div>
        </div>
    );
};

export default RegisterPage;
