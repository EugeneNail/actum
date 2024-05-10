import "./activity-card.sass"
import Icon from "../icon/icon.tsx";

export default function ActivityCard() {
    return (
        <div className="activity-card">
            <div className="activity-card__icon-container">
                <Icon className="activity-card__icon" name="edit_document" filled onClick={() => {}}/>
            </div>
            <p className="activity-card__label">Discord</p>
        </div>
    )
}