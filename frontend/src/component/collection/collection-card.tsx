import "./collection-card.sass"
import Collection from "../../model/collection.ts";
import ActivityCard from "../activity-card/activity-card.tsx";
import Button from "../button/button.tsx";
import {useNavigate} from "react-router-dom";
import Icon from "../icon/icon.tsx";

type Props = {
    collection: Collection
}

export default function CollectionCard({collection}: Props) {
    const navigate = useNavigate()

    function openDeletePage() {
        navigate(`/settings/collections/${collection.id}/delete`, {
            state: {
                icon: "category",
                name: collection.name,
                message: `Deleting "${collection.name}" will cause this collection and all of its activities to disappear from all of your records. The collection can be edited by simply clicking on the pencil icon.`,
                previousRoute: "/settings/collections"
            }
        })
    }

    function addActivity() {
        navigate("/settings/collections/activities/new", {
            state: {
                collectionId: collection.id,
                collectionName: collection.name
            }
        })
    }

    return (
        <div className="collection-card">
            <div className="collection-card__header">
                <Icon className="collection-card__icon" name="category"/>
                <p className="collection-card__label">{collection.name}</p>
                <Button className="collection-card__edit-button" icon="edit" negative onClick={() => navigate(`/settings/collections/${collection.id}`)}/>
                <Button className="collection-card__delete-button" icon="delete" negative onClick={openDeletePage}/>
            </div>
            <div className="collection-card__activities">
                <Button className="collection-card__new-activity" icon="add" negative onClick={addActivity}/>
            </div>
        </div>
    )
}