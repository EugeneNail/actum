import "./records-page.sass"
import {useEffect, useState} from "react";
import {useHttp} from "../../../service/use-http.ts";
import RecordCard from "../../../component/record-card/record-card.tsx";
import ShortRecord from "../../../model/short-record.ts";

export default function RecordsPage() {
    const http = useHttp()
    const [records, setRecords] = useState<ShortRecord[]>([])

    useEffect(() => {
        document.title = "Записи"
        if (records.length == 0) {
            fetchRecords()
        }
    }, [])

    async function fetchRecords() {
        const {data, status} = await http.post("/api/records-list", {
            cursor: new Date().toISOString().split("T")[0]
        })

        if (status == 200) {
            setRecords(data)
        }
    }

    return (
        <div className="records-page page">
            <div className="records-page__records">
                {records && records.map(record => (
                    <RecordCard key={record.id} record={record}/>
                ))}
            </div>
        </div>
    )
}