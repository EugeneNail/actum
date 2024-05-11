import "./delete-page.sass"
import {useLocation, useNavigate, useParams} from "react-router-dom";
import Button from "../../component/button/button.tsx";
import Icon from "../../component/icon/icon.tsx";
import {useHttp} from "../../service/use-http.ts";
import {FormEvent} from "react";

export function DeletePage() {
    const {state} = useLocation()
    const {id} = useParams<string>()
    const {pathname} = window.location
    const http = useHttp()
    const navigate = useNavigate()
    let endpoint = ""


    if (pathname.includes("collections")) {
        endpoint = "/collections/" + id
    }

    if (pathname.includes("activities")) {
        endpoint = "/activities/" + id
    }

    if (state == null) {
        navigate(state?.previousRoute)
    }

    async function confirm(event: FormEvent) {
        event.preventDefault()
        await http.delete(endpoint)
        navigate(state?.previousRoute)
    }

    return (
        <div className="delete-page">
            <form className="delete-page__form">
                <div className="delete-page__cover">
                    <div className="delete-page__icon-container">
                        <Icon name={state?.icon} className="delete-page__icon"/>
                    </div>
                </div>
                <h3 className="delete-page__disclaimer">Delete "{state?.name}"?</h3>
                <p className="delete-page__message">{state?.message}</p>
                <div className="delete-page__button-container">
                    <Button className="delete-page__delete-button" label="Delete" pill onClick={confirm}/>
                    <Button className="delete-page__cancel-button" label="Cancel" pill onClick={() => navigate(state?.previousRoute)}/>
                </div>
            </form>
        </div>
    )
}