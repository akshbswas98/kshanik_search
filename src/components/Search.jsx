import React, { useEffect, useState, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useDebounce } from 'use-debounce';

import { Links } from './Links';

export const Search = () => {
    const navigate = useNavigate();
    const location = useLocation();

    // Initialize from URL query param so the input stays in sync
    const initialQuery = new URLSearchParams(location.search).get('q') || '';
    const [text, setText] = useState(initialQuery);
    const [debouncedValue] = useDebounce(text, 500);
    const prevDebouncedValue = useRef(debouncedValue);

    // Navigate when the debounced value changes (auto-search as you type)
    useEffect(() => {
        // Only navigate if the debounced value ACTUALLY changed (prevents redirect loops on tab switching)
        if (debouncedValue !== prevDebouncedValue.current) {
            prevDebouncedValue.current = debouncedValue;
            if (debouncedValue.trim()) {
                navigate(`/search?q=${encodeURIComponent(debouncedValue.trim())}`, { replace: true });
            }
        }
    }, [debouncedValue, navigate]);

    // Also navigate immediately on Enter key press
    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && text.trim()) {
            navigate(`/search?q=${encodeURIComponent(text.trim())}`);
        }
    };

    return (
        <div className="relative flex flex-col items-center mt-6">
            <input
                value={text}
                type="text"
                className="sm:w-96 w-80 h-12 dark:bg-gray-700 dark:text-white border dark:border-gray-600 rounded-full shadow-md outline-none p-6 text-black hover:shadow-lg transition-all duration-300 focus:shadow-xl focus:ring-2 focus:ring-green-500 dark:focus:ring-gray-500"
                placeholder="🔎 Search on Kshanik Search..."
                onChange={(e) => setText(e.target.value)}
                onKeyDown={handleKeyDown}
            />
            {text !== '' && (
                <button
                    type="button"
                    className="absolute top-3 right-4 text-2xl text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors duration-300"
                    onClick={() => setText('')}
                >
                    ×
                </button>
            )}
            <Links />
        </div>
    );
};
