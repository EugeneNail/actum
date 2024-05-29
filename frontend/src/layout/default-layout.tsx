import "./default-layout.sass"
import {Outlet} from "react-router-dom";

export default function DefaultLayout() {
    return (
        <div className="default-layout">
            <Outlet/>
        </div>
    )
}