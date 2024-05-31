import "./icon-select.sass"
import Icon from "../icon/icon.tsx";
import classNames from "classnames";
import {useState} from "react";
import {Icons8} from "../icon8/icons8.ts";
import Icon8 from "../icon8/icon8.tsx";

type Props = {
    label: string
    selectedIconId: number
    group: number
    setIcon: (icon: number) => void
}

export default function IconSelectList({label, selectedIconId, group, setIcon}: Props) {
    const [isVisible, setVisible] = useState(false)
    const ids = Object.keys(Icons8).map(key => Number(key)).filter(id => id >= group && id < group + 100)

    return (
        <div className={classNames("icon-select-list", {invisible: !isVisible})}>
            <div className="icon-select-list__header" onClick={() => setVisible(!isVisible)}>
                <p className="icon-select-list__name">{label}</p>
                <Icon className="icon-select-list__chevron" name={isVisible ? "keyboard_arrow_up" : "keyboard_arrow_down"}/>
            </div>
            <ul className="icon-select-list__list">
                {ids && ids.map(id => (
                    <li key={id} className={classNames("icon-select-list__item", {selected: id == selectedIconId })} onClick={() => setIcon(id)}>
                        <Icon8 className="icon-select-list__icon" id={id}/>
                    </li>
                ))}
            </ul>
        </div>
    )
}