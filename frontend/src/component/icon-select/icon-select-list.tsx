import "./icon-select.sass"
import Icon from "../icon/icon.tsx";
import classNames from "classnames";
import {useState} from "react";

type Props = {
    label: string
    selectedIcon: string
    icons: string[]
    setIcon: (icon: string) => void
}

export default function IconSelectList({label, selectedIcon, icons, setIcon}: Props) {
    const [isVisible, setVisible] = useState(true)

    return (
        <div className={classNames("icon-select-list", {invisible: !isVisible})}>
            <div className="icon-select-list__header" onClick={() => setVisible(!isVisible)}>
                <p className="icon-select-list__name">{label}</p>
                <Icon className="icon-select-list__chevron" name={isVisible ? "keyboard_arrow_up" : "keyboard_arrow_down"}/>
            </div>
            <ul className="icon-select-list__list">
                {icons && icons.map(icon => (
                    <li className={classNames("icon-select-list__item", {selected: icon == selectedIcon })} onClick={() => setIcon(icon)}>
                        <Icon className="icon-select-list__icon" name={icon}/>
                    </li>
                ))}
            </ul>
        </div>
    )
}