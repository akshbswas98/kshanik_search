import React, { useEffect, useState } from 'react';
import { useDebounce } from 'use-debounce';
import { useLocation, useNavigate } from 'react-router-dom';

import { Links } from './Links';

export const Search = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [text, setText] = useState('');
    const [debouncedValue] = useDebounce(text, 300);

    useEffect(() => {
        const query = new URLSearchParams(location.search).get('q') || '';
        setText(query);
    }, [location.search]);

    useEffect(() => {
        const query = debouncedValue.trim();
        if (location.pathname !== '/search') {
            return;
        }

        if (!query) {
            navigate('/search', { replace: true });
            return;
        }

        const currentQuery = new URLSearchParams(location.search).get('q') || '';
        if (currentQuery !== query) {
            navigate(`/search?q=${encodeURIComponent(query)}`, { replace: true });
        }
    }, [debouncedValue, location.pathname, location.search, navigate]);

    return (
        <div className="relative flex flex-col items-center mt-6">
            <input
                value={text}
                type="text"
                className="sm:w-96 w-80 h-12 dark:bg-gray-700 dark:text-white border dark:border-gray-600 rounded-full shadow-md outline-none p-6 text-black hover:shadow-lg transition-all duration-300 focus:shadow-xl focus:ring-2 focus:ring-green-500 dark:focus:ring-gray-500"
                placeholder="🔎 Search on Kshanik Search..."
                onChange={(e) => setText(e.target.value)}
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
