import React from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'
import GuestLayout from "./layout/guest-layout.tsx";
import SignupPage from "./page/guest/signup-page.tsx";
import "./shared.sass"
import LoginPage from "./page/guest/login-page.tsx";
import NotFoundPage from "./page/not-found-page/not-found-page.tsx";
import CollectionsPage from "./page/collection/collections-page/collections-page.tsx";
import SaveCollectionPage from "./page/collection/save-collection-page/save-collection-page.tsx";
import DefaultLayout from "./layout/default-layout.tsx";
import {DeletePage} from "./page/delete-page/delete-page.tsx";

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
        <Routes>
            <Route path="/not-found" element={<NotFoundPage/>}/>
            <Route element={<GuestLayout/>}>
                <Route path="/signup" element={<SignupPage/>}/>
                <Route path="/login" element={<LoginPage/>}/>
            </Route>
            <Route element={<DefaultLayout/>}>
                <Route path="/settings/collections" element={<CollectionsPage/>}>
                    <Route path="/settings/collections/new" element={<SaveCollectionPage/>}/>
                    <Route path="/settings/collections/:id" element={<SaveCollectionPage/>}/>
                    <Route path="/settings/collections/:id/delete" element={<DeletePage/>}/>
                </Route>
            </Route>
            <Route path="*" element={<Navigate to="/not-found"/>}/>
        </Routes>
    </BrowserRouter>
  </React.StrictMode>
)
