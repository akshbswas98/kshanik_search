import React, { createContext, useContext, useState } from 'react';

const ResultContext = createContext();

export const ResultContextProvider = ({ children }) => {
    const [results, setResults] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');

    const getResults = async (query) => {
        try {
            const params = new URLSearchParams({
                key: import.meta.env.VITE_GOOGLE_API_KEY, // Google API Key
                cx: import.meta.env.VITE_GOOGLE_CX,       // Custom Search Engine ID
                q: query
            });

            const response = await fetch(`https://www.googleapis.com/customsearch/v1?${params.toString()}`, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                }
            });

            if (!response.ok) {
                const errorData = await response.text();
                console.error('Error response:', errorData);
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorData}`);
            }

            const data = await response.json();

            console.log('Response:', data);

            if (data.error) {
                throw new Error(data.error.message);
            }

            setResults(data.items?.map((item) => ({
                title: item.title,
                link: item.link,
                snippet: item.snippet,
            })) || []);
        } catch (error) {
            console.error('Error fetching data from Google Programmable Search:', error);
            setResults([]);
        }
    };

    return (
        <ResultContext.Provider value={{ getResults, results, searchTerm, setSearchTerm }}>
            {children}
        </ResultContext.Provider>
    );
};

export const useResultContext = () => useContext(ResultContext);