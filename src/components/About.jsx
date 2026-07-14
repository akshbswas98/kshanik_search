import React from 'react';

export const About = () => {
    return (
        <div className="bg-slate-50 dark:bg-slate-900 min-h-screen font-sans text-slate-800 dark:text-slate-200 leading-relaxed transition-colors duration-300">
            {/* Hero Section */}
            <main className="relative">
                <div className="absolute top-0 left-0 right-0 h-48 bg-gradient-to-b from-emerald-50 dark:from-emerald-950/20 to-transparent -z-10"></div>
                <section className="max-w-4xl mx-auto px-4 pt-16 pb-12 text-center">
                    {/* Portrait Image Container */}
                    <div className="relative inline-block mb-8">
                        <div className="w-48 h-48 rounded-full border-4 border-white dark:border-slate-800 shadow-xl overflow-hidden bg-slate-200 dark:bg-slate-800">
                            <img 
                                src="/kshanik_kumar_biswas.jpg" 
                                alt="Shri Kshanik Kumar Biswas" 
                                className="w-full h-full object-cover"
                            />
                        </div>
                    </div>
                    {/* Header Titles */}
                    <h1 className="text-4xl md:text-5xl font-serif font-bold text-slate-900 dark:text-white mb-2">Shri Kshanik Kumar Biswas</h1>
                    <p className="text-emerald-700 dark:text-emerald-400 font-medium text-lg mb-8 italic">A Great Police Officer, A Beloved Grandfather</p>
                    {/* Tribute Text */}
                    <div className="max-w-3xl mx-auto">
                        <p className="text-lg text-slate-600 dark:text-slate-400 leading-relaxed">
                            This meta-search engine, <strong>Kshanik Search</strong>, is built in loving memory of my late grandfather,
                            <strong> Shri Kshanik Kumar Biswas</strong>. He was a dedicated and esteemed Police Officer who spent
                            his life serving with honor, courage, and integrity. He passed away peacefully on May 9, 2021,
                            on the highly auspicious day of Rabindra Jayanti. His legacy of righteousness, wisdom, and
                            unconditional love continues to inspire us. This platform blends his enduring spirit with modern web performance.
                        </p>
                    </div>
                    {/* Divider Ornament */}
                    <div className="flex justify-center mt-12 opacity-30 dark:opacity-50">
                        <svg fill="none" height="24" viewBox="0 0 200 24" width="200" xmlns="http://www.w3.org/2000/svg">
                            <path d="M0 12H90M200 12H110M95 12C95 9.23858 97.2386 7 100 7C102.761 7 105 9.23858 105 12C105 14.7614 102.761 17 100 17C97.2386 17 95 14.7614 95 12Z" stroke="currentColor" strokeWidth="1.5"></path>
                        </svg>
                    </div>
                </section>
            </main>

            {/* Content Grids */}
            <section className="max-w-7xl mx-auto px-4 pb-24 grid grid-cols-1 lg:grid-cols-2 gap-12">
                {/* Our Mission Block */}
                <div>
                    <h2 className="text-3xl font-serif font-bold text-emerald-800 dark:text-emerald-400 mb-8 flex items-center justify-center lg:justify-start gap-4">
                        Our Mission
                    </h2>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-sm border border-slate-100 dark:border-slate-700 hover:shadow-md transition-all">
                            <h3 className="font-bold text-lg text-slate-900 dark:text-white mb-3">Legacy in Code</h3>
                            <p className="text-sm text-slate-600 dark:text-slate-400">{"Honoring Kshanik Kumar Biswas's values through a robust, reliable search tool that embodies trust and clarity."}</p>
                        </div>
                        <div className="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-sm border border-slate-100 dark:border-slate-700 hover:shadow-md transition-all">
                            <h3 className="font-bold text-lg text-slate-900 dark:text-white mb-3">Information Access</h3>
                            <p className="text-sm text-slate-600 dark:text-slate-400">Providing a unified, efficient search experience across diverse sources, empowering users with knowledge.</p>
                        </div>
                    </div>
                </div>

                {/* Tech Stack Block */}
                <div>
                    <h2 className="text-3xl font-serif font-bold text-emerald-800 dark:text-emerald-400 mb-8 flex items-center justify-center lg:justify-start gap-4">
                        Tech Stack
                    </h2>
                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
                        {/* Go Card */}
                        <div className="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-sm border border-slate-100 dark:border-slate-700 hover:shadow-md transition-all flex flex-col items-start">
                            <div className="mb-4 text-[#00ADD8]">
                                <span className="material-symbols-outlined text-4xl">terminal</span>
                            </div>
                            <h3 className="font-bold text-lg text-slate-900 dark:text-white mb-2">Go</h3>
                            <p className="text-xs text-slate-600 dark:text-slate-400">High-performance, scalable backend for fast, secure data processing.</p>
                        </div>
                        {/* React Card */}
                        <div className="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-sm border border-slate-100 dark:border-slate-700 hover:shadow-md transition-all flex flex-col items-start">
                            <div className="mb-4 text-[#61DAFB]">
                                <span className="material-symbols-outlined text-4xl">rebase_edit</span>
                            </div>
                            <h3 className="font-bold text-lg text-slate-900 dark:text-white mb-2">React</h3>
                            <p className="text-xs text-slate-600 dark:text-slate-400">Modern, dynamic front-end library for building an interactive UI.</p>
                        </div>
                        {/* Vite Card */}
                        <div className="bg-white dark:bg-slate-800 p-6 rounded-xl shadow-sm border border-slate-100 dark:border-slate-700 hover:shadow-md transition-all flex flex-col items-start">
                            <div className="mb-4 text-[#FFD62E]">
                                <span className="material-symbols-outlined text-4xl">bolt</span>
                            </div>
                            <h3 className="font-bold text-lg text-slate-900 dark:text-white mb-2">Vite</h3>
                            <p className="text-xs text-slate-600 dark:text-slate-400">Next-generation frontend tooling for rapid development.</p>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    );
};
