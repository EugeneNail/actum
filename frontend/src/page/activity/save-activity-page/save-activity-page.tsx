import "./save-activity-page.sass"
import {useFormState} from "../../../service/use-form-state.ts";
import { useNavigate, useParams} from "react-router-dom";
import {useHttp} from "../../../service/use-http.ts";
import {useEffect, useState} from "react";
import Field from "../../../component/field/field.tsx";
import Button from "../../../component/button/button.tsx";
import IconsList from "./icons-list.tsx";
import {icons} from "../../../assets/icons.ts"
import Icon from "../../../component/icon/icon.tsx";
import {useNotificationContext} from "../../../component/notification/notification.tsx";

class Payload {
    name: string = ""
    icon: string = ""
    collectionId: number = 0
}

class Errors {
    name: string = ""
    icon: string = ""
    collectionId: string = ""
}

export default function SaveActivityPage() {
    const {state, setField, setState, errors, setErrors} = useFormState(new Payload(), new Errors())
    const http = useHttp()
    const navigate = useNavigate()
    const {collectionId, id} = useParams<string>()
    const notification = useNotificationContext()
    const willCreate = window.location.pathname.includes("/new")
    const [collectionName, setCollectionName] =useState<string>()

    useEffect(() => {
        setState({
            ...state,
            collectionId: parseInt(collectionId ?? "0"),
            icon: "Man"
        })
        fetchCollection()

        if (!willCreate) {
            fetchActivity()
        }
    }, [])

    async function fetchCollection() {
        const {data, status} = await http.get(`/collections/${collectionId}`)
        if (status == 403 || status == 404) {
            notification.pop(data)
            navigate("/collections")
        }
        setCollectionName(data.name)
    }

    async function fetchActivity() {
        const {data, status} = await http.get(`/activities/${id}`)
        if (status == 403 || status == 404) {
            notification.pop(data)
            navigate("/collections")
            return
        }
        setState({
            ...state,
            name: data.name,
            icon: data.icon ?? "Man"
        })
    }

    function setIcon(icon: string) {
        setState({
            ...state,
            icon: icon
        })
    }

    async function create() {
        const {data, status} = await http.post("/activities", state)

        console.log(state)

        if (status == 422) {
            setErrors(data)
            return
        }

        if (status == 201) {
            navigate("/collections")
        }
    }

    async function edit() {
        const {data, status} = await http.put(`/activities/${id}`, {
            name: state.name,
            icon: state.icon
        })

        if (status == 422) {
            setErrors(data)
            return
        }

        if (status == 204) {
            navigate("/collections")
        }
    }

    return (
        <div className="save-activity-page">
            <div className="cover" onClick={() => navigate("/collections")}/>
            <form className="form" onSubmit={e => e.preventDefault()}>
                <Button className="form__delete-button" icon="delete" negative onClick={() => navigate(`./delete`)}/>
                <div className="form__cover">
                    <div className="form__icon-container">
                        <Icon name={state?.icon} className="form__icon" filled/>
                    </div>
                </div>
                <h3 className="form__title">{willCreate ? "Add" : "Edit"} activity {willCreate ? "to" : "of"} collection "{collectionName}"</h3>
                <div className="form__content">
                    <Field value={state.name} name="name" label="Name" icon="label" error={errors.name} onChange={setField}/>
                        <IconsList currentIcon={state.icon} label="People" icons={icons.people} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Nature" icons={icons.nature} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Food" icons={icons.food} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Home" icons={icons.home} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Transport" icons={icons.transport} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Buildings" icons={icons.buildings} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Activities" icons={icons.activities} onSelect={setIcon}/>
                    <Button className="form__button" pill accent label={willCreate ? "Add activity" : "Edit activity"} onClick={() => willCreate ? create() : edit()}/>
                </div>
            </form>
        </div>
    )
}