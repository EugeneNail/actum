import React from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'
import GuestLayout from "./layout/guest-layout.tsx";
import SignupPage from "./page/guest/signup-page.tsx";
import "./shared.sass"
import LoginPage from "./page/guest/login-page.tsx";
import NotFoundPage from "./page/not-found-page/not-found-page.tsx";

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
        <Routes>
            <Route path="/not-found" element={<NotFoundPage/>}/>
            <Route element={<GuestLayout/>}>
                <Route path="/signup" element={<SignupPage/>}/>
                <Route path="/login" element={<LoginPage/>}/>
            </Route>
            <Route path="*" element={<Navigate to="/not-found"/>}/>
        </Routes>
    </BrowserRouter>
  </React.StrictMode>
)
