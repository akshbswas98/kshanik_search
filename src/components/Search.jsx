import React, { useEffect, useState, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useDebounce } from 'use-debounce';

import { Links } from './Links';

export const Search = ({ variant = 'hero' }) => {
    const navigate = useNavigate();
    const location = useLocation();

    // Initialize from URL query param
    const initialQuery = new URLSearchParams(location.search).get('q') || '';
    const [text, setText] = useState(initialQuery);
    const [debouncedValue] = useDebounce(text, 500);
    const prevDebouncedValue = useRef(debouncedValue);

    useEffect(() => {
        if (debouncedValue !== prevDebouncedValue.current) {
            prevDebouncedValue.current = debouncedValue;
            if (debouncedValue.trim()) {
                navigate(`/search?q=${encodeURIComponent(debouncedValue.trim())}`, { replace: true });
            }
        }
    }, [debouncedValue, navigate]);

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && text.trim()) {
            navigate(`/search?q=${encodeURIComponent(text.trim())}`);
        }
    };

    if (variant === 'header') {
        return (
            <div className="relative flex items-center w-full">
                <span className="material-symbols-outlined absolute left-3 text-slate-400 text-xl">search</span>
                <input
                    value={text}
                    type="text"
                    className="w-full pl-10 pr-10 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-md focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none text-sm transition-all"
                    placeholder="Search archives..."
                    onChange={(e) => setText(e.target.value)}
                    onKeyDown={handleKeyDown}
                />
                {text !== '' && (
                    <button
                        type="button"
                        className="absolute right-3 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
                        onClick={() => setText('')}
                    >
                        ×
                    </button>
                )}
            </div>
        );
    }

    return (
        <div className="w-full max-w-2xl px-4 flex flex-col items-center">
            <div className="relative w-full group search-glow bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 transition-all duration-300">
                <div className="flex items-center px-6 py-5">
                    <span className="material-symbols-outlined text-slate-400 mr-4 text-2xl">search</span>
                    <input
                        value={text}
                        type="text"
                        className="w-full bg-transparent border-none focus:ring-0 text-slate-900 dark:text-slate-100 font-sans text-lg placeholder:text-slate-400/60 outline-none"
                        placeholder="Search records, memorials, or archives..."
                        onChange={(e) => setText(e.target.value)}
                        onKeyDown={handleKeyDown}
                    />
                    <button className="ml-4 flex items-center justify-center w-10 h-10 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors">
                        <span className="material-symbols-outlined text-primary dark:text-emerald-400">settings_voice</span>
                    </button>
                </div>
            </div>
            {/* Quick Suggestions / Links */}
            <div className="mt-6 flex flex-wrap justify-center gap-3">
                <span className="text-sm text-slate-500 py-1">Common searches:</span>
                {['Genealogy', 'Digital Obituaries', 'Historical Records'].map((tag) => (
                    <button
                        key={tag}
                        onClick={() => {
                            setText(tag);
                            navigate(`/search?q=${encodeURIComponent(tag)}`);
                        }}
                        className="text-sm text-slate-600 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 hover:bg-emerald-100 dark:hover:bg-emerald-900/30 hover:text-emerald-800 dark:hover:text-emerald-300 px-3 py-1 rounded transition-colors"
                    >
                        {tag}
                    </button>
                ))}
            </div>
            <div className="mt-8">
                <Links />
            </div>
        </div>
    );
};
