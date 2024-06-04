import "./activity-picker.sass"
import Collection from "../../model/collection.ts";
import Icon from "../icon/icon.tsx";
import classNames from "classnames";
import {useState} from "react";
import PickerActivity from "./picker-activity.tsx";

type Props = {
    collection: Collection
    value: number[]
    toggleActivity: (id: number) => void
}

export default function PickerCollection({collection, value, toggleActivity}: Props) {
    const [isVisible, setVisible] = useState(true)

    return (
        <div className={classNames("picker-collection", {invisible: !isVisible})}>
            <div className="picker-collection__header" onClick={() => setVisible(!isVisible)}>
                <p className="picker-collection__name">{collection.name}</p>
                <Icon className="picker-collection__chevron" name={isVisible ? "keyboard_arrow_up" : "keyboard_arrow_down"}/>
            </div>
            <ul className="picker-collection__activities">
                {collection.activities && collection.activities.map(activity => (
                    <PickerActivity key={activity.id} activity={activity} color={collection.color} toggled={value.includes(activity.id)} toggle={toggleActivity}/>
                ))}
            </ul>
        </div>
    )
}