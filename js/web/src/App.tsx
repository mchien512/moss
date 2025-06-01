import React, { useState, useEffect } from 'react';
import EntryForm from './components/EntryForm';
import './App.css';

const App: React.FC = () => {
    const [status, setStatus] = useState<string>('Ready to create entries');
    const [lastCreatedEntry, setLastCreatedEntry] = useState<string>('');

    const handleEntryCreated = (entryId: string) => {
        setLastCreatedEntry(entryId);
        setStatus(`✅ Last created entry: ${entryId}`);
    };

    return (
        <div className="app">
            <header className="app-header">
                <h1>Moss Entry Manager</h1>
                <p className="status">{status}</p>

                {lastCreatedEntry && (
                    <div className="last-entry">
                        <p>Last created entry ID: <code>{lastCreatedEntry}</code></p>
                    </div>
                )}
            </header>

            <main className="app-main">
                <EntryForm onEntryCreated={handleEntryCreated} />
            </main>
        </div>
    );
};

export default App;