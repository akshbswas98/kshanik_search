import React, { useEffect, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import { useResultContext } from '../contexts/ResultsContextProvider.jsx';

export const Results = () => {
  const { results, getResults, searchTerm } = useResultContext();
  const location = useLocation();
  const searchMadeRef = useRef(false);

  useEffect(() => {
    if (searchTerm !== '' && !searchMadeRef.current) {
      getResults(searchTerm);
      searchMadeRef.current = true;
    }
  }, [searchTerm]);

  if (!results || results.length === 0) {
    return null;
  }

  return (
    <div className="sm:px-56 flex flex-wrap justify-between space-y-6 p-4">
      {results.map(({ title, link, snippet }, index) => (
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

