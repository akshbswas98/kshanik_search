import React from 'react'
import { Routes as RouterRoutes, Route, Navigate } from 'react-router-dom';
import { Results } from './Results';

export const Routes = () => {
    return (
        <div className="p-4">
             <RouterRoutes>
                <Route path="/" element={<Navigate to="/search" replace />} />
                <Route path="/search" element={<Results />} />
            </RouterRoutes>
        </div>
    );
}


   