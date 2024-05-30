import "./activity-card.sass"
import Activity from "../../model/activity.ts";
import Icon from "../icon/icon.tsx";

type Props ={
    activity: Activity
}

export default function ActivityCard({activity}: Props) {
    function formatName(): string {
        const name = activity.name
        if (name.length > 16) {
            return name.substring(0, 14).trim() + "..."
        }
        return name
    }

    return (
        <div className="activity-card">
            <div className="activity-card__icon-container">
                <Icon name={activity.icon} className="activity-card__icon"/>
            </div>
            <p className="activity-card__name">{formatName()}</p>
        </div>
    )
}