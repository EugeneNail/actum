import "./record-card.sass"
import ShortCollection from "../../model/short-collection.ts";
import RecordCardActivity from "./record-card-activity.tsx";
import classNames from "classnames";
import {Color} from "../../model/color.tsx";
type Props = {
    collection: ShortCollection
}

export default function RecordCardCollection({collection}: Props) {
    const nameClassName = classNames(
        "record-card-collection__name",
        {red: collection.color == Color.Red},
        {orange: collection.color == Color.Orange},
        {yellow: collection.color == Color.Yellow},
        {green: collection.color == Color.Green},
        {blue: collection.color == Color.Blue},
        {purple: collection.color == Color.Purple},
    )

    return (
        <div className="record-card-collection">
            <p className={nameClassName}>{collection.name}</p>
            <div className="record-card-collection__activities">
                {collection.activities && collection.activities.map(activity => (
                    <RecordCardActivity key={Math.random()} activity={activity}/>
                ))}
            </div>
        </div>
    )
}