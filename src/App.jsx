import React from 'react';
import { Navbar } from './components/Navbar';
import { Footer } from './components/Footer';
import { Routes } from './components/Routes';

const App = () => {
    return (
        <div className="bg-gray-50 dark:bg-gray-900 dark:text-gray-200 min-h-screen transition-colors duration-300">
            <Navbar />
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <Routes />
            </main>
            <Footer />
        </div>
    );
};

export default App;
