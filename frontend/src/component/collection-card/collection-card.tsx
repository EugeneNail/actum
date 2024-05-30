import "./collection-card.sass"
import Collection from "../../model/collection";
import ActivityCard from "../activity-card/activity-card.tsx";
import Icon from "../icon/icon.tsx";
import {useNavigate} from "react-router-dom";

type Props = {
    collection: Collection
}

export default function CollectionCard({collection}: Props) {
    const navigate = useNavigate()

    return (
        <div className="collection-card">
            <div className="collection-card__title-container" onClick={() => navigate(`./${collection.id}`)}>
                <h6 className="collection-card__title">{collection.name}</h6>
            </div>
            <div className="collection-card__activities">
                {collection.activities && collection.activities.map(activity => (
                    <ActivityCard key={activity.id} activity={activity}/>
                ))}
                {(collection.activities?.length < 20 || collection.activities == null) && <div className="collection-card-button">
                    <div className="collection-card-button__icon-container">
                        <Icon className="collection-card-button__icon" name="add" bold/>
                    </div>
                    <p className="collection-card-button__label">Add activity</p>
                </div>}
            </div>
        </div>
    )
}