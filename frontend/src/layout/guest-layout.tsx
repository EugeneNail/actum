import {Outlet} from "react-router-dom";

export default function GuestLayout() {
    return (
        <div className="guest-layout">
            <Outlet/>
        </div>
    )
}
