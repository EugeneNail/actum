import "./records-page.sass"
import {useEffect, useState} from "react";
import {useApi} from "../../../service/use-api.ts";
import RecordCard from "../../../component/record-card/record-card.tsx";
import ShortRecord from "../../../model/short-record.ts";
import Throbber from "../../../component/throbber/throbber.tsx";
import Icon from "../../../component/icon/icon.tsx";
import {useNavigate} from "react-router-dom";

export default function RecordsPage() {
    const navigate = useNavigate()
    const [isLoading, setLoading] = useState(true)
    const api = useApi()
    const [records, setRecords] = useState<ShortRecord[]>([])
    const messages = [
        "Давайте не будем оставлять эту страницу дневника пустой? ✌",
        "Давайте продолжим с того места, где вы остановились. 🙌",
        "Что ни день, то новая история. 👏",
        "Сделайте перерыв и добавьте запись на сегодня. ✍"
    ]


    useEffect(() => {
        document.title = "Записи"
        if (records?.length == 0) {
            setLoading(true)
            fetchRecords().then()
        }
    }, [])

    async function fetchRecords() {
        const {data, status} = await api.post("/api/records-list", {
            cursor: new Date().toISOString().split("T")[0]
        })

        if (status == 200) {
            setRecords(data)
            setLoading(false)
        }

        // data.forEach(r => console.log(r.date))
    }


    function checkIfToday(record: ShortRecord): boolean {
        const today = new Date().toISOString().split("T")[0] + "T00:00:00Z"
        return record.date == today
    }


    return (
        <div className="records-page page">
            {isLoading && <Throbber/>}
            {!isLoading &&
                <div className="records-page__records">
                    {records && !records.some(checkIfToday) &&
                        <div className="records-page-button" onClick={() => navigate("/records/new")}>
                        <div className="records-page-button__icon-container">
                            <Icon className="records-page-button__icon" name="add"/>
                        </div>
                        <p className="records-page-button__label">{messages[Math.floor(Math.random() * messages.length)]}</p>
                    </div>}
                    {records && records.map(record => (
                        <RecordCard key={record.id} record={record}/>
                    ))}
                </div>
            }
        </div>
    )
}