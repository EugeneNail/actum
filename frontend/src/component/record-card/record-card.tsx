import "./record-card.sass"
import ShortRecord from "../../model/short-record.ts";
import Icon from "../icon/icon.tsx";
import {Mood, MoodIcons} from "../../model/mood.ts";
import classNames from "classnames";
import RecordCardCollection from "./record-card-collection.tsx";
import {useNavigate} from "react-router-dom";

type Props = {
    record: ShortRecord
}

const months = ["Января", "Февраля", "Марта", "Апреля", "Мая", "Июня", "Июля", "Августа", "Сентября", "Октября", "Ноября", "Декабря"]
const weekdays = ["Воскресенье", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"]

export default function RecordCard({record}: Props) {
    const navigate = useNavigate()
    const moodClassName = classNames(
        "record-card__mood",
        {radiating: record.mood == Mood.Radiating},
        {happy: record.mood == Mood.Happy},
        {neutral: record.mood == Mood.Neutral},
        {bad: record.mood == Mood.Bad},
        {awful: record.mood == Mood.Awful},
    )


    function formatDate(): string {
        const date = new Date(record.date)
        return `${weekdays[date.getDay()]}, ${date.getDate()} ${months[date.getMonth()]}`
    }


    return (
        <div className="record-card">
            <div className="record-card__header" onClick={() => navigate(`./${record.id}`)}>
                <Icon className={moodClassName} name={MoodIcons[record.mood]}/>
                <p className="record-card__date">{formatDate()}</p>
            </div>
            <div className="record-card__collections">
                {record.collections && record.collections.map(collection =>
                    collection.activities?.length > 0 && <RecordCardCollection key={Math.random()} collection={collection}/>
                )}
            </div>
            {record.notes.length > 0 && <p className="record-card__notes">{record.notes}</p>}
        </div>
    )
}