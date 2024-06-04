import "./activity-picker.sass"
import {Color} from "../../model/color.tsx";
import Activity from "../../model/activity.ts";
import classNames from "classnames";
import Icon8 from "../icon8/icon8.tsx";

type Props = {
    activity: Activity
    color: Color
    toggled: boolean
    toggle: (id: number) => void
}

export default function PickerActivity({activity, color, toggled, toggle}: Props) {
    const className = classNames(
        "picker-activity",
        {toggled: toggled},
        {red: color == Color.Red},
        {orange: color == Color.Orange},
        {yellow: color == Color.Yellow},
        {green: color == Color.Green},
        {blue: color == Color.Blue},
        {purple: color == Color.Purple},
    )

    return (
        <div className={className} onClick={() => toggle(activity.id)}>
            <div className="picker-activity__icon-container">
                <Icon8 id={activity.icon} className="picker-activity__icon"/>
            </div>
            <p className="picker-activity__name">{activity.name}</p>
        </div>
    )
}