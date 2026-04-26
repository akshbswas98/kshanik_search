import React from 'react'
import { Routes as RouterRoutes, Route } from 'react-router-dom';
import { Results } from './Results';
import { About } from './About';
import { Home } from './Home';

export const Routes = () => {
    return (
        <div className="py-8">
            <RouterRoutes>
                <Route path="/" element={<Home />} />
                <Route path="/search" element={<Results />} />
                <Route path="/about" element={<About />} />
            </RouterRoutes>
        </div>
    );
}
