const NOTIFICATION_SERVICE_URL = 'ws://localhost:8083/ws';

class WebSocketService {
    constructor() {
        this.ws = null;
        this.listeners = [];
    }

    connect(userId) {
        if (this.ws) {
            this.ws.close();
        }

        this.ws = new WebSocket(`${NOTIFICATION_SERVICE_URL}?userId=${userId}`);

        this.ws.onopen = () => {
            console.log('Connected to Notification Service');
        };

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                // Notify all listeners
                this.listeners.forEach(listener => listener(data.content));

                // Keep the alert as a fallback for now
                const message = data.content.message || "New Notification";
                alert(`New Notification: ${message}`);
            } catch (error) {
                console.error('Error parsing notification:', error);
            }
        };

        this.ws.onclose = () => {
            console.log('Disconnected from Notification Service');
        };
    }

    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }

    registerCallback(callback) {
        this.listeners.push(callback);
    }

    unregisterCallback(callback) {
        this.listeners = this.listeners.filter(cb => cb !== callback);
    }
}

export default new WebSocketService();
