import React from 'react'
import { Link } from 'react-router-dom'
import { Search } from './Search'

export const Navbar = ({ darkTheme, setDarkTheme }) => {
    return (
        <div className="p-5 pb-0 flex flex-wrap sm:justify-between justify-center items-center border-b dark:border-gray-700 border-gray-200 bg-white dark:bg-gray-800 shadow-sm">
            <div className="flex justify-between items-center space-x-5 w-screen">
                <p className="text-2xl bg-gradient-to-r from-green-500 to-emerald-600 font-bold text-white py-2 px-4 rounded-lg shadow-md hover:shadow-lg transition-all duration-300 dark:from-gray-600 dark:to-gray-700">
                    Kshanik Search
                </p>
                <button 
                    type="button" 
                    onClick={() => setDarkTheme(!darkTheme)} 
                    className="text-xl bg-white dark:bg-gray-700 border dark:border-gray-600 rounded-full px-4 py-2 hover:shadow-lg transition-all duration-300 hover:scale-105"
                >
                    {darkTheme ? '💡' : '🌙'}
                </button>
            </div>
            <div className="flex space-x-4 mt-4">
                <Link to="/search" className="text-lg text-gray-700 dark:text-gray-300 hover:underline">Search</Link>
                <Link to="/about" className="text-lg text-gray-700 dark:text-gray-300 hover:underline">About</Link>
            </div>
            <Search />
        </div>
    );
};



