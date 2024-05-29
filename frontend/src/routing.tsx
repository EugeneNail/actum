import {BrowserRouter, Navigate, Route, Routes, useLocation} from "react-router-dom";
import Notification from "./component/notification/notification.tsx";
import {CSSTransition, TransitionGroup} from "react-transition-group";
import LoginPage from "./page/guest-page/login-page.tsx";
import SignupPage from "./page/guest-page/signup-page.tsx";

function AnimatedRoutes() {
    const location = useLocation()

    return (
        <div className="animation-wrapper">
            <TransitionGroup>
                <CSSTransition key={location.key} classNames="slide" timeout={300}>
                    <Routes location={location}>
                        <Route path="/Login" element={<LoginPage/>}/>
                        <Route path="/signup" element={<SignupPage/>}/>
                        <Route path="/" element={<Navigate to="/records"/>}/>
                        {/*<Route path="*" element={<Navigate to="/not-found"/>}/>*/}
                    </Routes>
                </CSSTransition>
            </TransitionGroup>
        </div>
    )
}

export default function Routing() {
    return (
        <Notification>
            <BrowserRouter>
                <AnimatedRoutes/>
            </BrowserRouter>
        </Notification>
    )
}