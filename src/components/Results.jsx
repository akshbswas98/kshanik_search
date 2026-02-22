import React, { useEffect, useMemo, useState } from 'react';
import { useLocation } from 'react-router-dom';

const getApiBase = () => {
    const configured = import.meta.env.VITE_SEARCH_API_BASE_URL;
    return configured && configured.trim() ? configured.trim().replace(/\/$/, '') : '/api';
};

export const Results = () => {
    const location = useLocation();
    const [results, setResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const apiBase = useMemo(() => getApiBase(), []);

    useEffect(() => {
        const controller = new AbortController();

        const fetchResults = async () => {
            try {
                const searchTerm = new URLSearchParams(location.search).get('q')?.trim() || '';
                if (!searchTerm) {
                    setResults([]);
                    setError(null);
                    return;
                }

                setLoading(true);
                setError(null);

                const response = await fetch(`${apiBase}/search?q=${encodeURIComponent(searchTerm)}`, {
                    method: 'GET',
                    headers: { Accept: 'application/json' },
                    signal: controller.signal
                });

                if (!response.ok) {
                    const errorData = await response.text();
                    throw new Error(`HTTP ${response.status}: ${errorData}`);
                }

                const data = await response.json();
                if (!Array.isArray(data)) {
                    throw new Error('Invalid API response format. Expected an array.');
                }

                setResults(data);
            } catch (err) {
                if (err.name === 'AbortError') {
                    return;
                }
                setError('Failed to fetch search results. Please try again later.');
                setResults([]);
            } finally {
                setLoading(false);
            }
        };

        fetchResults();
        return () => controller.abort();
    }, [apiBase, location.search]);

    if (error) {
        return <div className="p-4 text-red-500">Error: {error}</div>;
    }

    if (loading) {
        return <p className="text-center text-gray-500 p-4">Loading results...</p>;
    }

    if (!new URLSearchParams(location.search).get('q')) {
        return <p className="text-center text-gray-500 p-4">Type a query above to start searching.</p>;
    }

    return (
        <div className="sm:px-56 flex flex-wrap justify-between space-y-6 p-4">
            {results.length > 0 ? (
                results.map(({ title, url, snippet, source }, index) => (
                    <div key={`${url}-${index}`} className="md:w-2/5 w-full">
                        <a
                            href={url}
                            target="_blank"
                            rel="noreferrer"
                            className="block p-4 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-300"
                        >
                            <p className="text-lg dark:text-blue-300 text-blue-700 font-medium hover:underline">
                                {title}
                            </p>
                            <p className="text-xs uppercase tracking-wide text-gray-500 mt-1">{source}</p>
                            <p className="text-sm dark:text-gray-300 text-gray-700 mt-2">
                                {snippet || 'No description available.'}
                            </p>
                        </a>
                    </div>
                ))
            ) : (
                <p className="text-center text-gray-500">No results found. Try searching for something else.</p>
            )}
        </div>
    );
};
