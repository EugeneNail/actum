import "./default-layout.sass"
import {Outlet} from "react-router-dom";
import Header from "../component/header/header.tsx";

export default function DefaultLayout() {
    return (
        <div className="default-layout">
            <Header/>
            <Outlet/>
        </div>
    )
}