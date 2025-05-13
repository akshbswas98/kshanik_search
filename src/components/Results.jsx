import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

export const Results = () => {
    const location = useLocation();
    const [results, setResults] = useState([]);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchResults = async () => {
            try {
                const searchTerm = new URLSearchParams(location.search).get('q') || 'default search';
                const params = new URLSearchParams({
                    key: import.meta.env.VITE_GOOGLE_API_KEY, // Google API Key
                    cx: import.meta.env.VITE_GOOGLE_CX,       // Custom Search Engine ID
                    q: searchTerm
                });

                const response = await fetch(`https://www.googleapis.com/customsearch/v1?${params.toString()}`);

                if (!response.ok) {
                    const errorData = await response.text();
                    throw new Error(`HTTP error! status: ${response.status}, message: ${errorData}`);
                }

                const data = await response.json();
                console.log('API Response:', data); // Debugging API response
                setResults(data.items || []);
            } catch (err) {
                console.error('Error fetching data from Google Programmable Search API:', err);
                setError(err.message);
            }
        };

        fetchResults();
    }, [location]);

    if (error) {
        return <div className="p-4 text-red-500">Error: {error}</div>;
    }

    return (
        <div className="sm:px-56 flex flex-wrap justify-between space-y-6 p-4">
            {results.length > 0 ? (
                results.map(({ title, link, snippet }, index) => (
                    <div key={index} className="md:w-2/5 w-full">
                        <a 
                            href={link} 
                            target="_blank" 
                            rel="noreferrer"
                            className="block p-4 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors duration-300"
                        >
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

