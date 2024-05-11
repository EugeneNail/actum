import "./header.sass"
import {NavLink} from "react-router-dom";
import Icon from "../icon/icon.tsx";

type Props = {
    icon: string
    label: string
    to: string
}

export default function HeaderLink({icon, label, to}:Props) {
    return (
        <NavLink className="header-link" to={to}>
            <Icon className="header-link__icon" name={icon}/>
            <p className="header-link__label">{label}</p>
        </NavLink>
    )
}