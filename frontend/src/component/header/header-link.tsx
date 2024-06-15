import "./header.sass"
import Icon from "../icon/icon.tsx";
import {NavLink} from "react-router-dom";
import classNames from "classnames";
import {Color} from "../../model/color.tsx";

type Props = {
    label: string
    color: Color,
    icon: string,
    to: string
}

export default function HeaderLink({label, color, icon, to}:Props) {
    const className = classNames(
        "header-link",
        {red: color == Color.Red},
        {orange: color == Color.Orange},
        {yellow: color == Color.Yellow},
        {green: color == Color.Green},
        {blue: color == Color.Blue},
        {purple: color == Color.Purple},
        {accent: color == Color.Accent}
    )

    return (
        <NavLink className={className} to={to}>
            <Icon className="header-link__icon" name={icon}/>
            <p className="header-link__label">{label}</p>
        </NavLink>
    )
}