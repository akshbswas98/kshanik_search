import  React, {createContext, useContext, useState } from 'react';
const ResultContext = createContext();
const baseUrl = 'https://google-search3.p.rapidapi.com/api/v1';

export const ResultContextProvider = ({ children }) => {
    const [results, setResults] = useState([]);
    const [isLoading, setLoading] = useState(false);
    const [searchTerm, setSearchTerm] = useState('');

    const getResults = async (url) => {

        setLoading(true);
        const res = await fetch(`${baseUrl}${url}`, {
            method: 'GET',
            headers: {
                'x-user-agent': 'desktop',
                'x-proxy-location': 'US',
                'x-rapidapi-host': 'google-search3.p.rapidapi.com',
                'x-rapidapi-key': '8c18d905d5msh7e052759d7d49e0p15af1ejsn982175305efd',
            }
        });

        const data = await res.json();
        setResults(data);
    setLoading(false);
    }
    return (
        <ResultContext.Provider value={{ getResults, results, searchTerm , setSearchTerm, isLoading }}>
            {children}
        </ResultContext.Provider>
        
    )
    }
export const useResultContext = () => useContext(ResultContext);
