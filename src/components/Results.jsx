import React, { useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';

const mockResults = [
    {
        title: 'Example Result 1',
        link: 'https://example.com/1',
        snippet: 'This is a snippet for example result 1.'
    },
    {
        title: 'Example Result 2',
        link: 'https://example.com/2',
        snippet: 'This is a snippet for example result 2.'
    },
    {
        title: 'Example Result 3',
        link: 'https://example.com/3',
        snippet: 'This is a snippet for example result 3.'
    }
];

export const Results = () => {
    const location = useLocation();
    const searchMadeRef = useRef(false);

    useEffect(() => {
        // Placeholder for search logic
        console.log('Search triggered for location:', location);
    }, [location]);

    return (
        <div className="sm:px-56 flex flex-wrap justify-between space-y-6 p-4">
            {mockResults.map(({ title, link, snippet }, index) => (
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
            ))}
        </div>
    );
};

