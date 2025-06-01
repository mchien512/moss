import React, { useState, useEffect } from 'react';
import { healthCheck } from './api/client';
import './App.css';

const App: React.FC = () => {
    const [status, setStatus] = useState<string>('Connecting...');

    useEffect(() => {
        const checkBackend = async () => {
            try {
                await healthCheck();
                setStatus('Connected to backend ✓');
            } catch (error) {
                setStatus('Backend not available ✗');
            }
        };

        checkBackend();
    }, []);

    return (
        <div className="app">
            <header className="app-header">
                <h1>Moss</h1>
                <p>{status}</p>
                <button className="btn-primary">Get Started</button>
            </header>
        </div>
    );
};

export default App;