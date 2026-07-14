import React, { createContext, useContext, useState } from 'react';

const ResultContext = createContext();

export const ResultContextProvider = ({ children }) => {
    const [results, setResults] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const getResults = async (query) => {
        setIsLoading(true);
        try {
            const params = new URLSearchParams({ q: query });
            const API_BASE = import.meta.env.VITE_BACKEND_URL || '';
            const endpoint = API_BASE ? `${API_BASE}/search?${params.toString()}` : `/api/search?${params.toString()}`;
            
            const response = await fetch(endpoint);

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            setResults(data.results || []);
        } catch (error) {
            console.error('Error fetching search results:', error);
            setResults([]);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <ResultContext.Provider value={{ getResults, results, searchTerm, setSearchTerm, isLoading }}>
            {children}
        </ResultContext.Provider>
    );
};

export const useResultContext = () => useContext(ResultContext);
