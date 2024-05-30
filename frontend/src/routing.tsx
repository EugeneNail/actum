import {BrowserRouter, Navigate, Route, Routes} from "react-router-dom";
import Notification from "./component/notification/notification.tsx";
import LoginPage from "./page/guest-page/login-page.tsx";
import SignupPage from "./page/guest-page/signup-page.tsx";
import DefaultLayout from "./layout/default-layout.tsx";
import CollectionsPage from "./page/collections/collections-page/collections-page.tsx";
import SaveCollectionPage from "./page/collections/save-collection-page/save-collections-page.tsx";

export default function Routing() {
    return (
        <Notification>
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={<LoginPage/>}/>
                    <Route path="/signup" element={<SignupPage/>}/>
                    <Route element={<DefaultLayout/>}>
                        <Route path="/collections" element={<CollectionsPage/>}/>
                        <Route path="/collections/new" element={<SaveCollectionPage/>}/>
                        <Route path="/collections/:id" element={<SaveCollectionPage/>}/>
                    </Route>
                    <Route path="/" element={<Navigate to="/records"/>}/>
                </Routes>
            </BrowserRouter>
        </Notification>
    )
}