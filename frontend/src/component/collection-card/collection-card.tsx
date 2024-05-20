import "./collection-card.sass"
import Collection from "../../model/collection.ts";
import Button from "../button/button.tsx";
import {useNavigate} from "react-router-dom";
import Icon from "../icon/icon.tsx";
import ActivityCard from "../activity-card/activity-card.tsx";

type Props = {
    collection: Collection
}

export default function CollectionCard({collection}: Props) {
    const navigate = useNavigate()

    return (
        <div className="collection-card">
            <div className="collection-card__header">
                <Icon className="collection-card__icon" name="category"/>
                <p className="collection-card__label">{collection.name}</p>
                <Button className="collection-card__edit-button" icon="edit" negative onClick={() => navigate(`./${collection.id}`)}/>
                <Button className="collection-card__delete-button" icon="delete" negative onClick={() => navigate(`./${collection.id}/delete`)}/>
            </div>
            <div className="collection-card__activities">
                {collection.activities && collection.activities.map(activity => (
                    <ActivityCard activity={activity}/>
                ))}
                <Button className="collection-card__new-activity" icon="add" negative onClick={() => navigate(`./${collection.id}/activities/new`)}/>
            </div>
        </div>
    )
}