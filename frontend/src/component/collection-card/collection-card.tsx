import "./collection-card.sass"
import Collection from "../../model/collection";
import ActivityCard from "../activity-card/activity-card.tsx";
import Icon from "../icon/icon.tsx";
import {useNavigate} from "react-router-dom";
import classNames from "classnames";
import {Color} from "../../model/color.tsx";

type Props = {
    collection: Collection
}

export default function CollectionCard({collection}: Props) {
    const navigate = useNavigate()
    const titleClassName = classNames(
        "collection-card__title",
        {red: collection.color == Color.Red},
        {orange: collection.color == Color.Orange},
        {yellow: collection.color == Color.Yellow},
        {green: collection.color == Color.Green},
        {blue: collection.color == Color.Blue},
        {purple: collection.color == Color.Purple},
    )
    const buttonIconContainerClassName = classNames(
        "collection-card-button__icon-container",
        {red: collection.color == Color.Red},
        {orange: collection.color == Color.Orange},
        {yellow: collection.color == Color.Yellow},
        {green: collection.color == Color.Green},
        {blue: collection.color == Color.Blue},
        {purple: collection.color == Color.Purple},
    )
    const buttonIconClassName = classNames(
        "collection-card-button__icon",
        {red: collection.color == Color.Red},
        {orange: collection.color == Color.Orange},
        {yellow: collection.color == Color.Yellow},
        {green: collection.color == Color.Green},
        {blue: collection.color == Color.Blue},
        {purple: collection.color == Color.Purple},
    )

    return (
        <div className="collection-card">
            <div className="collection-card__title-container" onClick={() => navigate(`./${collection.id}`)}>
                <h6 className={titleClassName}>{collection.name}</h6>
            </div>
            <div className="collection-card__activities">
                {collection.activities && collection.activities.map(activity => (
                    <ActivityCard key={activity.id} activity={activity} collectionId={collection.id}/>
                ))}
                {(collection.activities?.length < 20 || collection.activities == null) && <div className="collection-card-button" onClick={() => navigate(`./${collection.id}/activities/new`)}>
                    <div className={buttonIconContainerClassName} >
                        <Icon className={buttonIconClassName} name="add" bold/>
                    </div>
                    <p className="collection-card-button__label">Add activity</p>
                </div>}
            </div>
        </div>
    )
}