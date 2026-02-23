import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

export const Results = () => {
    const location = useLocation();
    const [results, setResults] = useState([]);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);
    const [timeTaken, setTimeTaken] = useState('');
    const [totalCount, setTotalCount] = useState(0);

    useEffect(() => {
        const fetchResults = async () => {
            const searchTerm = new URLSearchParams(location.search).get('q');
            if (!searchTerm) {
                setResults([]);
                return;
            }

            setLoading(true);
            setError(null);

            try {
                const params = new URLSearchParams({ q: searchTerm });
                const response = await fetch(`/api/search?${params.toString()}`);

                if (!response.ok) {
                    const errorData = await response.json().catch(() => ({}));
                    throw new Error(errorData.error || `HTTP ${response.status}`);
                }

                const data = await response.json();

                if (data.results && data.results.length > 0) {
                    setResults(data.results);
                    setTimeTaken(data.time_taken || '');
                    setTotalCount(data.total_count || 0);
                } else {
                    setResults([]);
                }
            } catch (err) {
                console.error('Error fetching search results:', err);
                setError('Failed to fetch search results. Please try again later.');
            } finally {
                setLoading(false);
            }
        };

        fetchResults();
    }, [location]);

    if (loading) {
        return (
            <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-green-500"></div>
                <span className="ml-3 text-gray-500">Searching...</span>
            </div>
        );
    }

    if (error) {
        return <div className="p-4 text-red-500">Error: {error}</div>;
    }

    return (
        <div className="sm:px-56 flex flex-col space-y-6 p-4">
            {totalCount > 0 && (
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Found {totalCount} results in {timeTaken}
                </p>
            )}
            {results.length > 0 ? (
                results.map(({ title, url, snippet, source, score }, index) => (
                    <div key={index} className="w-full">
                        <a
                            href={url}
                            target="_blank"
                            rel="noreferrer"
                            className="block p-4 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-300"
                        >
                            <div className="flex items-center gap-2 mb-1">
                                <span className="text-xs px-2 py-0.5 rounded-full bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-300 font-medium uppercase">
                                    {source}
                                </span>
                                <span className="text-xs text-gray-400 dark:text-gray-500 truncate">
                                    {url}
                                </span>
                            </div>
                            <p className="text-lg dark:text-blue-300 text-blue-700 font-medium hover:underline">
                                {title}
                            </p>
                            <p className="text-sm dark:text-gray-300 text-gray-700 mt-2">
                                {snippet}
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
