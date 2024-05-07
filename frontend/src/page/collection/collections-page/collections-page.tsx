import {Outlet} from "react-router-dom";
import "./collections-page.sass"

export default function CollectionsPage() {
    return (
        <div className="collections-page">
            <Outlet/>
        </div>
    )
}