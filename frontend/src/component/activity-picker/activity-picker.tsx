import "./activity-picker.sass"
import Collection from "../../model/collection.ts";
import PickerCollection from "./picker-collection.tsx";

type Props = {
    collections: Collection[]
    value: number[]
    toggleActivity: (id: number) => void
}

export default function ActivityPicker({collections, value, toggleActivity}: Props) {
    return (
        <div className="activity-picker">
            <p className="activity-picker__label">Чем вы занимались?</p>
            {collections && collections.map(collection =>
                collection.activities && collection.activities.length > 0 && (
                    <PickerCollection key={collection.id} collection={collection} value={value} toggleActivity={toggleActivity}/>
                )
            )}
        </div>
    )
}