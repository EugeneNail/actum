import "./save-activity-page.sass"
import {useFormState} from "../../../service/use-form-state.ts";
import {useLocation, useNavigate} from "react-router-dom";
import {useHttp} from "../../../service/use-http.ts";
import {FormEvent} from "react";
import Field from "../../../component/field/field.tsx";
import Button from "../../../component/button/button.tsx";
import IconsList from "./icons-list.tsx";
import {icons} from "../../../assets/icons.ts"
import Icon from "../../../component/icon/icon.tsx";

class Payload {
    name: string = ""
    icon: string = ""
    collectionId: number = 0
}

class Errors {
    name: string = ""
    icon: string = "man"
    collectionId: string = ""
}

export default function SaveActivityPage() {
    const {state, setField, setState, errors, setErrors} = useFormState(new Payload(), new Errors())
    const http = useHttp()
    const navigate = useNavigate()
    const location = useLocation()

    function setIcon(icon: string) {
        setState({
            ...state,
            icon: icon
        })
    }

    async function submit(event: FormEvent) {
        event.preventDefault()
        setState({
            ...state,
            collectionId: location.state.collectionId
        })

        console.log(state)

        const {data, status} = await http.post("/activities", state)

        if (status == 422) {
            setErrors(data)
            return
        }

        if (status == 201) {
            navigate("/settings/collections")
        }
    }

    return (
        <div className="save-activity-page">
            <div className="cover" onClick={() => navigate("/settings/collections")}/>
            <form className="form" onSubmit={submit}>
                <div className="form__cover">
                    <div className="form__icon-container">
                        <Icon name={state?.icon} className="form__icon" filled/>
                    </div>
                </div>
                <h3 className="form__title">Add activity to collection "{location.state.collectionName}"</h3>
                <div className="form__content">
                    <Field value={state.name} name="name" label="Name" icon="label" error={errors.name} onChange={setField}/>
                        <IconsList currentIcon={state.icon} label="People" icons={icons.people} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Nature" icons={icons.nature} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Food" icons={icons.food} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Home" icons={icons.home} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Transport" icons={icons.transport} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Buildings" icons={icons.buildings} onSelect={setIcon}/>
                        <IconsList currentIcon={state.icon} label="Activities" icons={icons.activities} onSelect={setIcon}/>
                    <Button className="form__button" pill label="Add activity"/>
                </div>
            </form>
        </div>
    )
}