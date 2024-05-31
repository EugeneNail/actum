import "./header.sass"
import {NavLink} from "react-router-dom";
import Icon from "../icon/icon.tsx";

export default function Header() {
    return (
        <header className="header">
            <NavLink className="header-link" to="/statistics">
                <Icon className="header-link__icon" name="bar_chart" />
            </NavLink>
            <NavLink className="header-link" to="/records">
                <Icon className="header-link__icon" name="post" />
            </NavLink>
            <NavLink className="header-link" to="/collections">
                <Icon className="header-link__icon" name="category" />
            </NavLink>
        </header>
    )
}