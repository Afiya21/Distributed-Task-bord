import React from 'react';
import '../index.css'; // Ensure styles are available

const Modal = ({ isOpen, title, message, icon = '🎉', onClose }) => {
    if (!isOpen) return null;

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <span className="modal-icon">{icon}</span>
                <h2 style={{ marginBottom: '0.5rem' }}>{title}</h2>
                <p style={{ color: 'var(--text-secondary)' }}>{message}</p>
                {onClose && (
                    <button
                        onClick={onClose}
                        className="btn btn-primary"
                        style={{ marginTop: '1.5rem', width: '100%', justifyContent: 'center' }}
                    >
                        Close
                    </button>
                )}
            </div>
        </div>
    );
};

export default Modal;
