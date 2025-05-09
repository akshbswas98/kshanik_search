import React from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter as Router} from 'react-router-dom';
import App from './App.jsx';
import { ResultContextProvider } from './contexts/ResultsContextProvider.jsx';
import './global.css';

const root = createRoot(document.getElementById('root'));
root.render(
    <ResultContextProvider>
        <Router>
            <App />
        </Router>
    </ResultContextProvider>
);
