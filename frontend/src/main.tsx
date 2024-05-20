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
import {DeleteCollectionPage} from "./page/collection/delete-collection-page/delete-collection-page.tsx";
import Notification from "./component/notification/notification.tsx";
import SaveActivityPage from "./page/activity/save-activity-page/save-activity-page.tsx";
import {DeleteActivityPage} from "./page/activity/delete-activity-page/delete-activity-page.tsx";

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Notification>
        <BrowserRouter>
            <Routes>
                <Route path="/not-found" element={<NotFoundPage/>}/>
                <Route element={<GuestLayout/>}>
                    <Route path="/signup" element={<SignupPage/>}/>
                    <Route path="/login" element={<LoginPage/>}/>
                </Route>
                <Route element={<DefaultLayout/>}>
                    <Route path="/collections" element={<CollectionsPage/>}>
                        <Route path="/collections/new" element={<SaveCollectionPage/>}/>
                        <Route path="/collections/:id" element={<SaveCollectionPage/>}/>
                        <Route path="/collections/:id/delete" element={<DeleteCollectionPage/>}/>
                        <Route path="/collections/:collectionId/activities/new" element={<SaveActivityPage/>}/>
                        <Route path="/collections/:collectionId/activities/:id" element={<SaveActivityPage/>}/>
                        <Route path="/collections/:collectionId/activities/:id/delete" element={<DeleteActivityPage/>}/>
                    </Route>
                </Route>
                <Route path="*" element={<Navigate to="/not-found"/>}/>
            </Routes>
        </BrowserRouter>
    </Notification>
  </React.StrictMode>
)
