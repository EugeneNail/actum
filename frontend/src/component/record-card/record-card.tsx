import "./record-card.sass"
import ShortRecord from "../../model/short-record.ts";
import Icon from "../icon/icon.tsx";
import {Mood, MoodIcons} from "../../model/mood.ts";
import classNames from "classnames";
import RecordCardCollection from "./record-card-collection.tsx";
import {useNavigate} from "react-router-dom";
import {Weather, WeatherIcons, WeatherNames} from "../../model/weather.ts";

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


    const weatherClassName = classNames(
        "record-card__weather",
        {heat: record.weather == Weather.Heat},
        {sunny: record.weather == Weather.Sunny},
        {cloudy: record.weather == Weather.Cloudy},
        {windy: record.weather == Weather.Windy},
        {rainy: record.weather == Weather.Rainy},
        {thunder: record.weather == Weather.Thunder},
        {foggy: record.weather == Weather.Foggy},
        {snowy: record.weather == Weather.Snowy},
        {cold: record.weather == Weather.Cold},
    )


    function formatDate(): string {
        const date = new Date(record.date)
        return `${weekdays[date.getDay()]}, ${date.getDate()} ${months[date.getMonth()]}`
    }


    return (
        <div className="record-card" onClick={() => navigate(`./${record.id}`)}>
            <div className="record-card__header">
                <Icon className={moodClassName} name={MoodIcons[record.mood]}/>
                <Icon className={weatherClassName} name={WeatherIcons[record.weather]}/>
                <div className="record-card__label">
                    <p className="record-card__date">{formatDate()}</p>
                    <p className="record-card__weather-name">{WeatherNames[record.weather]}</p>
                </div>
            </div>
            <div className="record-card__collections">
                {record.collections && record.collections.map(collection =>
                    collection.activities?.length > 0 && <RecordCardCollection key={Math.random()} collection={collection}/>
                )}
            </div>
            {record.notes.length > 0 && <p className="record-card__notes">{record.notes}</p>}
            {record.photos && <div className="record-card__photos">
                {record.photos.map(photo =>
                    <div className="record-card__photo-container" key={photo}>
                        <img className="record-card__photo" src={`${import.meta.env.VITE_API_DOMAIN}/api/photos/${photo}`} alt={photo} />
                    </div>
                )}
            </div>}
        </div>
    )
}