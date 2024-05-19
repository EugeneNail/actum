import "./delete-collection-page.sass"
import {useNavigate, useParams} from "react-router-dom";
import Button from "../../../component/button/button.tsx";
import Icon from "../../../component/icon/icon.tsx";
import {useHttp} from "../../../service/use-http.ts";
import {FormEvent, useEffect, useState} from "react";

export function DeleteCollectionPage() {
    const {id} = useParams<string>()
    const http = useHttp()
    const navigate = useNavigate()
    const [state, setState] = useState({name: "", message: ""})

    useEffect(() => {
        fetchCollection()
    }, [])

    async function fetchCollection() {
        const {data, status} = await http.get(`/collections/${id}`)

        if (status == 403) {
            navigate("/settings/collections")
            return
        }

        setState({
            name: data.name,
            message: `Deleting "${data.name}" will cause this collection and all of its activities to disappear from all of your records. The collection can be edited by simply clicking on the pencil icon.`
        })
    }

    async function confirm(event: FormEvent) {
        event.preventDefault()
        await http.delete(`/collections/${id}`)
        navigate("/settings/collections")
    }

    return (
        <div className="delete-collection-page">
            <form className="delete-collection-page__form">
                <div className="delete-collection-page__cover">
                    <div className="delete-collection-page__icon-container">
                        <Icon name="category" className="delete-collection-page__icon"/>
                    </div>
                </div>
                <h3 className="delete-collection-page__disclaimer">Delete "{state?.name}"?</h3>
                <p className="delete-collection-page__message">{state?.message}</p>
                <div className="delete-collection-page__button-container">
                    <Button className="delete-collection-page__delete-button" label="Delete" pill onClick={confirm}/>
                    <Button className="delete-collection-page__cancel-button" label="Cancel" pill onClick={() => navigate("/settings/collections")}/>
                </div>
            </form>
        </div>
    )
}