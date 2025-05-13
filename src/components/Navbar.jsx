import React from 'react'
import { Link } from 'react-router-dom'
import { Search } from './Search'

export const Navbar = ({ darkTheme, setDarkTheme }) => {
    return (
        <div className="p-5 pb-0 flex flex-col sm:flex-row sm:justify-between justify-center items-center border-b dark:border-gray-700 border-gray-200 bg-white dark:bg-gray-800 shadow-md transition-all duration-300">
            <div className="flex justify-center sm:justify-start items-center space-x-5 mb-4 sm:mb-0">
                <p className="text-3xl bg-gradient-to-r from-green-500 to-emerald-600 font-bold text-white py-2 px-4 rounded-lg shadow-md hover:shadow-lg transition-transform duration-300 hover:scale-105 dark:from-gray-600 dark:to-gray-700">
                    Kshanik Search
                </p>
                <button 
                    type="button" 
                    onClick={() => setDarkTheme(!darkTheme)} 
                    className="text-xl bg-white dark:bg-gray-700 border dark:border-gray-600 rounded-full px-4 py-2 hover:shadow-lg transition-transform duration-300 hover:scale-110"
                >
                    {darkTheme ? '💡' : '🌙'}
                </button>
            </div>
            <div className="flex space-x-4">
                <Link to="/search" className="text-lg text-gray-700 dark:text-gray-300 hover:underline transition-transform duration-300 hover:scale-105 hover:text-green-500 dark:hover:text-green-300 px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 shadow-md">
                    Search
                </Link>
                <Link to="/about" className="text-lg text-gray-700 dark:text-gray-300 hover:underline transition-transform duration-300 hover:scale-105 hover:text-green-500 dark:hover:text-green-300 px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 shadow-md">
                    About
                </Link>
            </div>
            <Search />
        </div>
    );
};



