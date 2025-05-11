import React from 'react';

export const About = () => {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900 dark:text-gray-200 p-6">
            <h1 className="text-4xl font-bold mb-6">In Loving Memory</h1>
            <img 
                src="/public/kshanik_kumar_biswas.jpg" 
                alt="Kshanik Kumar Biswas" 
                className="w-64 h-64 rounded-full shadow-lg mb-6"
            />
            <p className="text-lg text-center max-w-2xl">
                This website is dedicated to the loving memory of my grandfather, Kshanik Kumar Biswas (1943-2021). 
                His wisdom, kindness, and unwavering support continue to inspire us every day. 
                May his soul rest in peace.
            </p>
        </div>
    );
};
