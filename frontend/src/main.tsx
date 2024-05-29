import React from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'
import "./shared.sass"
import Notification from "./component/notification/notification.tsx";
import SignupPage from "./page/guest-page/signup-page.tsx";
import LoginPage from "./page/guest-page/login-page.tsx";

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Notification>
        <BrowserRouter>
            <Routes>
                <Route path="/Login" element={<LoginPage/>}/>
                <Route path="/signup" element={<SignupPage/>}/>
                <Route path="/" element={<Navigate to="/records"/>}/>
                {/*<Route path="*" element={<Navigate to="/not-found"/>}/>*/}
            </Routes>
        </BrowserRouter>
    </Notification>
  </React.StrictMode>
)
