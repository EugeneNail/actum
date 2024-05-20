import "./save-activity-page.sass"
import Icon from "../../../component/icon/icon.tsx";
import classNames from "classnames";
import {useState} from "react";

type Props = {
    currentIcon: string
    label: string
    icons: string[]
    onSelect: (icon: string) => void
}

export default function IconsList({currentIcon, label, icons, onSelect}: Props) {
    const [isVisible, setVisible] = useState(true)

    return (
        <div className="icons-list">
            <h5 className="icons-list__label" onClick={() => setVisible(!isVisible)}>
                {label}
                <Icon className="icons-list__drop" name={isVisible ? "arrow_drop_up" : "arrow_drop_down"} />
            </h5>
            {isVisible && <div className="icons-list__container">
                {icons.map((icon, index) => (
                    <Icon key={index} className={classNames("icons-list__icon", {selected: currentIcon == icon})} name={icon} filled onClick={() => onSelect(icon)}/>
                ))}
            </div>}
        </div>
    )
}