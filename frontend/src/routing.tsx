import {BrowserRouter, Navigate, Route, Routes} from "react-router-dom";
import Notification from "./component/notification/notification.tsx";
import LoginPage from "./page/guest-page/login-page.tsx";
import SignupPage from "./page/guest-page/signup-page.tsx";
import DefaultLayout from "./layout/default-layout.tsx";
import CollectionsPage from "./page/collections/collections-page/collections-page.tsx";
import SaveCollectionPage from "./page/collections/save-collection-page/save-collection-page.tsx";
import DeleteCollectionPage from "./page/collections/delete-collection-page/delete-collection-page.tsx";
import SaveActivityPage from "./page/activities/save-activity-page/save-activity-page.tsx";
import DeleteActivityPage from "./page/activities/delete-activity-page/delete-activity-page.tsx";
import SaveRecordPage from "./page/records/save-record-page/save-record-page.tsx";
import RecordsPage from "./page/records/records-page/records-page.tsx";

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
                        <Route path="/collections/:id/delete" element={<DeleteCollectionPage/>}/>
                        <Route path="/collections/:collectionId/activities/new" element={<SaveActivityPage/>}/>
                        <Route path="/collections/:collectionId/activities/:activityId" element={<SaveActivityPage/>}/>
                        <Route path="/collections/:collectionId/activities/:activityId/delete" element={<DeleteActivityPage/>}/>
                        <Route path="/records/new" element={<SaveRecordPage/>}/>
                        <Route path="/records/:id" element={<SaveRecordPage/>}/>
                        <Route path="/records" element={<RecordsPage/>}/>
                    </Route>
                    <Route path="/" element={<Navigate to="/records"/>}/>
                </Routes>
            </BrowserRouter>
        </Notification>
    )
}