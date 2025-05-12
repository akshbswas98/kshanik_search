import React, { createContext, useContext, useState } from 'react';

const ResultContext = createContext();

export const ResultContextProvider = ({ children }) => {
    const [results, setResults] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');

    const getResults = async (query) => {
        try {
            const params = new URLSearchParams({
                q: query,
                api_key: import.meta.env.VITE_SERPAPI_KEY,
                engine: 'google',
                google_domain: 'google.com',
                gl: 'us',
                hl: 'en',
                num: '10'
            });

            const response = await fetch(`https://serpapi.com/search.json?${params.toString()}`, {
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

            const contentType = response.headers.get('content-type');
            if (!contentType || !contentType.includes('application/json')) {
                const errorData = await response.text();
                console.error('Non-JSON response:', errorData);
                throw new Error('Received non-JSON response');
            }

            const data = await response.json();

            console.log('Response:', data);

            if (data.error) {
                throw new Error(data.error);
            }

            setResults(data.organic_results?.map((item) => ({
                title: item.title,
                link: item.link,
                snippet: item.snippet,
            })) || []);
        } catch (error) {
            console.error('Error fetching data from SerpApi:', error);
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