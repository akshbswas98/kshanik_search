import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { Search } from './Search'

export const Navbar = ({ darkTheme, setDarkTheme }) => {
    const location = useLocation();
    const isHome = location.pathname === '/';

    return (
        <header className="sticky top-0 z-50 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800 transition-all duration-300">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex flex-col sm:flex-row items-center justify-between gap-4 py-4 sm:py-0">
                {/* Brand Logo */}
                <div className="flex items-center gap-6">
                    <Link to="/" className="font-serif text-3xl font-bold tracking-tight text-primary dark:text-emerald-400 hover:opacity-80 transition-opacity">
                        Kshanik
                    </Link>
                    <button 
                        type="button" 
                        onClick={() => setDarkTheme(!darkTheme)} 
                        className="p-2 rounded-full border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                        title="Toggle Theme"
                    >
                        {darkTheme ? '💡' : '🌙'}
                    </button>
                </div>

                {/* Navigation Links */}
                <nav className="flex items-center space-x-2 order-3 sm:order-2">
                    <Link 
                        to="/search" 
                        className={`px-4 py-2 rounded-md font-medium transition-colors ${
                            location.pathname === '/search' 
                            ? 'bg-emerald-100/50 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300' 
                            : 'text-slate-600 dark:text-slate-400 hover:text-emerald-800 dark:hover:text-emerald-300'
                        }`}
                    >
                        Search
                    </Link>
                    <Link 
                        to="/about" 
                        className={`px-4 py-2 rounded-md font-medium transition-colors ${
                            location.pathname === '/about' 
                            ? 'bg-emerald-100/50 text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-300' 
                            : 'text-slate-600 dark:text-slate-400 hover:text-emerald-800 dark:hover:text-emerald-300'
                        }`}
                    >
                        About
                    </Link>
                </nav>

                {/* Header Search Bar (Only shown when not on Home) */}
                {!isHome && (
                    <div className="flex-1 max-w-md w-full sm:ml-8 order-2 sm:order-3">
                        <Search variant="header" />
                    </div>
                )}
            </div>
        </header>
    );
};
