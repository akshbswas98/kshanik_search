import React, { useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';

export const Results = () => {
    const location = useLocation();
    const searchMadeRef = useRef(false);

    useEffect(() => {
        // Placeholder for search logic
        console.log('Search triggered for location:', location);
    }, [location]);

    return (
        <div className="sm:px-56 flex flex-wrap justify-between space-y-6 p-4">
            {/* Placeholder for rendering search results */}
            <p>No results to display.</p>
        </div>
    );
};

