import "./delete-collection-page.sass"
import {useNavigate, useParams} from "react-router-dom";
import Button from "../../../component/button/button.tsx";
import {useHttp} from "../../../service/use-http.ts";
import {useEffect, useState} from "react";
import FormHeader from "../../../component/form-header/form-header.tsx";

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
            navigate("/collections")
            return
        }

        setState({
            name: data.name,
            message: `Deleting "${data.name}" will cause this collection and all of its activities to disappear from all of your records. The collection can be edited by simply clicking on the pencil icon.`
        })
    }

    async function confirm() {
        await http.delete(`/collections/${id}`)
        navigate("/collections")
    }

    return (
        <div className="delete-collection-page" onSubmit={e => e.preventDefault()}>
            <div className="cover" onClick={() => navigate("/collections")}/>
            <form className="delete-collection-page__form">
                <FormHeader icon="category" title={`Delete "${state.name}"?`}/>
                <p className="delete-collection-page__message">{state?.message}</p>
                <div className="delete-collection-page__button-container">
                    <Button className="delete-collection-page__delete-button" label="Delete" pill onClick={confirm}/>
                    <Button className="delete-collection-page__cancel-button" label="Cancel" pill onClick={() => navigate("/collections")}/>
                </div>
            </form>
        </div>
    )
}