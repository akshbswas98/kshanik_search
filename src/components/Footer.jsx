import React from 'react';

export const Footer = () => (
  <footer className="bg-white dark:bg-slate-900 border-t border-slate-200 dark:border-slate-800 py-12 mt-auto transition-colors duration-300">
    <div className="max-w-7xl mx-auto px-4 text-center">
      <p className="text-slate-800 dark:text-slate-200 font-medium mb-2">
        © Kshanik Search ({new Date().getFullYear()} - Present)
      </p>
      <p className="text-slate-500 dark:text-slate-400 text-sm flex items-center justify-center gap-1.5">
        Made with <span className="text-red-500">❤️</span> in memory of Kshanik Kumar Biswas (1943-2021)
      </p>
    </div>
  </footer>
);
