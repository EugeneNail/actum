import "./records-page.sass"
import {useEffect, useState} from "react";
import {useApi} from "../../../service/use-api.ts";
import RecordCard from "../../../component/record-card/record-card.tsx";
import ShortRecord from "../../../model/short-record.ts";
import Throbber from "../../../component/throbber/throbber.tsx";

export default function RecordsPage() {
    const [isLoading, setLoading] = useState(true)
    const api = useApi()
    const [records, setRecords] = useState<ShortRecord[]>([])


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
    }

    return (
        <div className="records-page page">
            {isLoading && <Throbber/>}
            {!isLoading &&
                <div className="records-page__records">
                    {records && records.map(record => (
                        <RecordCard key={record.id} record={record}/>
                    ))}
                </div>
            }
        </div>
    )
}