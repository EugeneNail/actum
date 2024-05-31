import "./activity-card.sass"
import Activity from "../../model/activity.ts";
import Icon from "../icon/icon.tsx";
import {useNavigate} from "react-router-dom";

type Props ={
    collectionId: number
    activity: Activity
}

export default function ActivityCard({activity, collectionId}: Props) {
    const navigate = useNavigate()

    function formatName(): string {
        const name = activity.name
        if (name.length > 16) {
            return name.substring(0, 14).trim() + "..."
        }
        return name
    }

    return (
        <div className="activity-card" onClick={() => navigate(`./${collectionId}/activities/${activity.id}`)}>
            <div className="activity-card__icon-container">
                <Icon name={activity.icon} className="activity-card__icon"/>
            </div>
            <p className="activity-card__name">{formatName()}</p>
        </div>
    )
}