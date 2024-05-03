import React from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Route, Routes} from 'react-router-dom'
import GuestLayout from "./layout/guest-layout.tsx";
import SignupPage from "./page/guest/signup-page.tsx";
import "./shared.sass"
import LoginPage from "./page/guest/login-page.tsx";

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
        <Routes>
            <Route element={<GuestLayout/>}>
                <Route path="/signup" element={<SignupPage/>}/>
                <Route path="/login" element={<LoginPage/>}/>
            </Route>
        </Routes>
    </BrowserRouter>
  </React.StrictMode>
)
