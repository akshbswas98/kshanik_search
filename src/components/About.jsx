import React from 'react';

export const About = () => {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900 dark:text-gray-200 p-6">
            <h1 className="text-4xl font-bold mb-6 text-center text-green-600 dark:text-green-400">In Loving Memory</h1>

            <img
                src="/kshanik_kumar_biswas.jpg"
                alt="Shri Kshanik Kumar Biswas"
                className="w-64 h-64 rounded-full shadow-2xl mb-8 border-4 border-white dark:border-gray-800 object-cover"
            />

            <div className="max-w-3xl text-center bg-white dark:bg-gray-800 p-8 rounded-2xl shadow-lg border border-gray-100 dark:border-gray-700">
                <h2 className="text-2xl font-semibold mb-2">Shri Kshanik Kumar Biswas</h2>
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-6 italic">A Great Police Officer, A Beloved Grandfather</p>

                <p className="text-lg leading-relaxed mb-6 text-gray-700 dark:text-gray-300">
                    This search engine project is started in the loving memory of my late grandfather,
                    <strong> Shri Kshanik Kumar Biswas</strong>. He was a dedicated and esteemed Police Officer who spent his life serving with honor, courage, and integrity.
                </p>

                <p className="text-lg leading-relaxed mb-6 text-gray-700 dark:text-gray-300">
                    He passed away peacefully on <strong>May 9, 2021</strong>, on the highly auspicious day of <strong>Rabindra Jayanti</strong>.
                </p>

                <div className="border-t border-gray-200 dark:border-gray-700 pt-6 mt-6">
                    <p className="text-md italic text-gray-600 dark:text-gray-400">
                        "Your legacy of righteousness, wisdom, and unconditional love continues to inspire us. You will forever live on in our hearts."
                    </p>
                </div>
            </div>
        </div>
    );
};
