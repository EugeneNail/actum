import "./delete-activity-page.sass"
import {useNavigate, useParams} from "react-router-dom";
import Button from "../../../component/button/button.tsx";
import {useHttp} from "../../../service/use-http.ts";
import {useEffect, useState} from "react";
import FormHeader from "../../../component/form-header/form-header.tsx";

export function DeleteActivityPage() {
    const {id} = useParams<string>()
    const http = useHttp()
    const navigate = useNavigate()
    const [state, setState] = useState({name: "", icon:"", message: ""})

    useEffect(() => {
        fetchActivity()
    }, [])

    async function fetchActivity() {
        const {data, status} = await http.get(`/activities/${id}`)

        if (status == 403) {
            navigate("/collections")
            return
        }

        setState({
            name: data.name,
            icon: data.icon,
            message: `Deleting "${data.name}" will cause this activity to disappear from all of your records.`
        })
    }

    async function confirm() {
        await http.delete(`/activities/${id}`)
        navigate("/collections")
    }

    return (
        <div className="delete-activity-page" onSubmit={e => e.preventDefault()}>
            <div className="cover" onClick={() => navigate("/collections")}/>
            <form className="delete-activity-page__form">
                <FormHeader icon={state.icon} title={`Delete "${state.name}"?`}/>
                <p className="delete-activity-page__message">{state?.message}</p>
                <div className="delete-activity-page__button-container">
                    <Button className="delete-activity-page__delete-button" label="Delete" pill onClick={confirm}/>
                    <Button className="delete-activity-page__cancel-button" label="Cancel" pill onClick={() => navigate("/collections")}/>
                </div>
            </form>
        </div>
    )
}