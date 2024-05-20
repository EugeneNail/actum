import "./activity-card.sass"
import Icon from "../icon/icon.tsx";
import Activity from "../../model/activity.ts";
import {useNavigate} from "react-router-dom";

type Props = {
    activity: Activity
}

export default function ActivityCard({activity}: Props) {
    const navigate = useNavigate()

    return (
        <div className="activity-card" onClick={() => navigate(`./${activity.collectionId}/activities/${activity.id}`)}>
            <div className="activity-card__icon-container">
                <Icon className="activity-card__icon" name={activity.icon} filled />
            </div>
            <p className="activity-card__label">{activity.name}</p>
        </div>
    )
}