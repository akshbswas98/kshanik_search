import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

export const Results = () => {
    const location = useLocation();
    const [results, setResults] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [timeTaken, setTimeTaken] = useState('');
    const [totalCount, setTotalCount] = useState(0);

    useEffect(() => {
        const controller = new AbortController();

        const fetchResults = async () => {
            const searchTerm = new URLSearchParams(location.search).get('q');
            if (!searchTerm) {
                setResults([]);
                return;
            }

            setLoading(true);
            setError(null);

            try {
                const params = new URLSearchParams({ q: searchTerm });
                const API_BASE = import.meta.env.VITE_BACKEND_URL || '';
                const endpoint = API_BASE ? `${API_BASE}/search?${params.toString()}` : `/api/search?${params.toString()}`;
                
                const response = await fetch(endpoint, { signal: controller.signal });

                if (!response.ok) {
                    const errorData = await response.json().catch(() => ({}));
                    throw new Error(errorData.error || `HTTP ${response.status}`);
                }

                const data = await response.json();

                if (data.results && data.results.length > 0) {
                    setResults(data.results);
                    setTimeTaken(data.time_taken || '');
                    setTotalCount(data.total_count || 0);
                } else {
                    setResults([]);
                    setTimeTaken('');
                    setTotalCount(0);
                }
            } catch (err) {
                if (err.name === 'AbortError') {
                    return;
                }
                console.error('Error fetching search results:', err);
                setError('Failed to fetch search results. Please try again later.');
                setResults([]);
            } finally {
                setLoading(false);
            }
        };

        fetchResults();
        return () => controller.abort();
    }, [location.search]);

    if (loading) {
        return (
            <div className="flex flex-col justify-center items-center py-24">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-emerald-600 dark:border-emerald-400"></div>
                <span className="mt-4 text-slate-500 dark:text-slate-400 font-medium">Searching archives...</span>
            </div>
        );
    }

    if (error) {
        return (
            <div className="max-w-4xl mx-auto p-8 mt-8 bg-red-50 dark:bg-red-900/20 border border-red-100 dark:border-red-800 rounded-xl text-red-600 dark:text-red-400">
                <p className="font-bold mb-2">Search Error</p>
                <p>{error}</p>
            </div>
        );
    }

    const searchTerm = new URLSearchParams(location.search).get('q');

    return (
        <div className="max-w-4xl mx-auto px-4 py-8">
            {searchTerm && (
                <div className="mb-8 border-b border-slate-200 dark:border-slate-800 pb-4 flex justify-between items-end">
                    <div>
                        <h2 className="text-sm text-slate-500 dark:text-slate-400 uppercase tracking-wider font-bold mb-1">Search Results for</h2>
                        <p className="text-2xl font-serif font-bold text-slate-900 dark:text-white">{`"${searchTerm}"`}</p>
                    </div>
                    {totalCount > 0 && (
                        <p className="text-xs text-slate-400 dark:text-slate-500 font-medium italic">
                            Found {totalCount} records in {timeTaken}
                        </p>
                    )}
                </div>
            )}

            <div className="space-y-6">
                {results.length > 0 ? (
                    results.map(({ title, url, snippet, source, score }, index) => (
                        <article key={index} className="group relative bg-white dark:bg-slate-800 p-6 rounded-xl border border-slate-100 dark:border-slate-700 shadow-sm hover:shadow-md hover:border-emerald-200 dark:hover:border-emerald-900/50 transition-all duration-300">
                            <div className="flex items-center gap-3 mb-3">
                                <span className={`text-[10px] px-2 py-0.5 rounded-full font-bold uppercase tracking-tighter ${
                                    source === 'wikipedia' ? 'bg-slate-100 text-slate-700 dark:bg-slate-700 dark:text-slate-300' :
                                    source === 'github' ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-300' :
                                    'bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300'
                                }`}>
                                    {source}
                                </span>
                                <span className="text-xs text-slate-400 dark:text-slate-500 truncate max-w-[200px] sm:max-w-md">
                                    {url}
                                </span>
                            </div>
                            <h3 className="text-xl font-serif font-bold text-emerald-800 dark:text-emerald-400 group-hover:text-emerald-700 dark:group-hover:text-emerald-300 mb-2 transition-colors">
                                <a href={url} target="_blank" rel="noreferrer" className="hover:underline">
                                    {title}
                                </a>
                            </h3>
                            <p className="text-sm text-slate-600 dark:text-slate-400 leading-relaxed mb-4">
                                {snippet}
                            </p>
                            <div className="flex items-center justify-between pt-4 border-t border-slate-50 dark:border-slate-700/50">
                                <a 
                                    href={url} 
                                    target="_blank" 
                                    rel="noreferrer" 
                                    className="text-xs font-bold text-slate-400 hover:text-emerald-600 dark:hover:text-emerald-400 transition-colors flex items-center gap-1"
                                >
                                    Visit Source <span className="material-symbols-outlined text-sm">open_in_new</span>
                                </a>
                                {score && (
                                    <div className="flex items-center gap-1.5">
                                        <div className="w-16 h-1.5 bg-slate-100 dark:bg-slate-700 rounded-full overflow-hidden">
                                            <div 
                                                className="h-full bg-emerald-500" 
                                                style={{ width: `${score * 100}%` }}
                                            ></div>
                                        </div>
                                        <span className="text-[10px] text-slate-400 font-bold uppercase tracking-tighter">Relevance</span>
                                    </div>
                                )}
                            </div>
                        </article>
                    ))
                ) : searchTerm ? (
                    <div className="text-center py-24 bg-slate-50 dark:bg-slate-800/50 rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-700">
                        <span className="material-symbols-outlined text-6xl text-slate-300 dark:text-slate-600 mb-4 italic">find_in_page</span>
                        <p className="text-slate-500 dark:text-slate-400 font-medium">No historical records found for {`"${searchTerm}"`}.</p>
                        <p className="text-xs text-slate-400 mt-2 italic">Try refining your search terms or exploring the archives.</p>
                    </div>
                ) : null}
            </div>
        </div>
    );
};
