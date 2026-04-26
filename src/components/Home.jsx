import React from 'react';
import { Search } from './Search';

export const Home = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-[70vh] hero-gradient dark:bg-[#0a121e] transition-colors duration-300">
      {/* Tribute Card */}
      <div className="mb-12 flex flex-col items-center">
        <div className="relative group">
          <div className="absolute inset-0 rounded-full bg-primary/5 dark:bg-emerald-500/10 group-hover:bg-primary/10 transition-colors duration-500 scale-110"></div>
          <div className="w-32 h-32 rounded-full border-2 border-primary dark:border-emerald-500 p-1 bg-surface dark:bg-[#111c2e] relative">
            <img 
              className="w-full h-full rounded-full object-cover grayscale opacity-90" 
              src="/kshanik_kumar_biswas.jpg"
              alt="Shri Kshanik Kumar Biswas"
            />
          </div>
        </div>
        <div className="mt-6 text-center">
          <p className="text-xs text-secondary dark:text-slate-400 uppercase tracking-[0.15em] mb-2 font-bold">In Dedicated Memory</p>
          <h2 className="font-serif text-2xl text-primary dark:text-emerald-400 text-emerald-shadow">Shri Kshanik Kumar Biswas</h2>
          <p className="text-sm text-slate-500 dark:text-slate-400 mt-1 font-medium">1943 — 2021</p>
        </div>
      </div>

      {/* Hero Text */}
      <div className="text-center max-w-3xl mb-12 px-4">
        <h1 className="font-serif text-4xl md:text-5xl text-primary dark:text-white mb-6 leading-tight font-bold">
          In Loving Memory of <br /> Shri Kshanik Kumar Biswas
        </h1>
        <p className="text-lg text-secondary dark:text-slate-400 leading-relaxed max-w-2xl mx-auto">
          A memorial meta-search engine blending legacy with modern web performance. Preserve, discover, and celebrate the stories that define our history.
        </p>
      </div>

      {/* Search Bar Component */}
      <Search variant="hero" />

      {/* Feature Grid */}
      <div className="mt-24 grid grid-cols-1 md:grid-cols-3 gap-8 w-full max-w-6xl px-4">
        <div className="p-8 rounded-xl bg-white dark:bg-[#111c2e] border border-slate-100 dark:border-slate-800 hover:border-primary/20 transition-all duration-300 shadow-sm">
          <span className="material-symbols-outlined text-primary dark:text-emerald-400 mb-4 text-3xl">history_edu</span>
          <h3 className="font-serif text-xl text-primary dark:text-emerald-400 mb-3 font-bold">Living Archives</h3>
          <p className="text-secondary dark:text-slate-400 text-sm leading-relaxed">Dynamic records that grow with contributions from family and friends across generations.</p>
        </div>
        <div className="p-8 rounded-xl bg-white dark:bg-[#111c2e] border border-slate-100 dark:border-slate-800 hover:border-primary/20 transition-all duration-300 shadow-sm">
          <span className="material-symbols-outlined text-primary dark:text-emerald-400 mb-4 text-3xl">verified_user</span>
          <h3 className="font-serif text-xl text-primary dark:text-emerald-400 mb-3 font-bold">Verified Legacy</h3>
          <p className="text-secondary dark:text-slate-400 text-sm leading-relaxed">Secure, blockchain-verified digital footprints to ensure the authenticity of personal histories.</p>
        </div>
        <div className="p-8 rounded-xl bg-white dark:bg-[#111c2e] border border-slate-100 dark:border-slate-800 hover:border-primary/20 transition-all duration-300 shadow-sm">
          <span className="material-symbols-outlined text-primary dark:text-emerald-400 mb-4 text-3xl">auto_awesome</span>
          <h3 className="font-serif text-xl text-primary dark:text-emerald-400 mb-3 font-bold">AI Curation</h3>
          <p className="text-secondary dark:text-slate-400 text-sm leading-relaxed">Intelligent indexing that connects disparate historical records into a cohesive narrative.</p>
        </div>
      </div>
    </div>
  );
};
