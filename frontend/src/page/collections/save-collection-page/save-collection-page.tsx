import "./save-collection-page.sass"
import {useApi} from "../../../service/use-api.ts";
import {useFormState} from "../../../service/use-form-state.ts";
import {useEffect, useState} from "react";
import {Outlet, useNavigate, useParams} from "react-router-dom";
import {useNotificationContext} from "../../../component/notification/notification.tsx";
import Form from "../../../component/form/form.tsx";
import Field from "../../../component/field/field.tsx";
import FormButtons from "../../../component/form/form-button-container.tsx";
import FormBackButton from "../../../component/form/form-back-button.tsx";
import FormSubmitButton from "../../../component/form/form-submit-button.tsx";
import FormDeleteButton from "../../../component/form/form-delete-button.tsx";
import Palette from "../../../component/palette/palette.tsx";
import {Color} from "../../../model/color.tsx";
import Throbber from "../../../component/throbber/throbber.tsx";

class Payload {
    name = ""
    color: Color = Color.Red
}

class Errors {
    name = ""
    color = ""
}

export default function SaveCollectionPage() {
    const willStore = window.location.pathname.includes("/new")
    const [isLoading, setLoading] = useState(!willStore)
    const api = useApi()
    const {state, setState, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const navigate = useNavigate()
    const notification = useNotificationContext()
    const {id} = useParams()


    useEffect(() => {
        document.title = "Новая коллекция"
        if (!willStore) {
            fetchCollection().then()
        }
    }, [])


    async function fetchCollection() {
        const {data, status} = await api.get(`/api/collections/${id}`)

        if (status == 403) {
            notification.pop(data)
            navigate("/collections")
            return
        }

        if (status == 200) {
            setState({
                name: data.name,
                color: data.color
            })
            setLoading(false)
        }
        document.title = data.name + " - Коллекции"
    }


    async function save() {
        if (willStore) {
            await store()
        } else {
            await update()
        }
    }


    async function store() {
        const {data, status} = await api.post("/api/collections", {
            name: state.name,
            color: Number(state.color)
        })

        if (status == 422 || status == 409) {
            setErrors(data)
            return
        }

        navigate("/collections")
    }


    async function update() {
        const {data, status} = await api.put(`/api/collections/${id}`, {
            name: state.name,
            color: Number(state.color)
        })

        if (status == 403) {
            notification.pop(data)
            return
        }

        if (status == 422) {
            setErrors(data)
            return
        }

        navigate("/collections")
    }

    return (
        <div className="save-collection-page page">
            {isLoading && <Throbber/>}
            {!isLoading &&
                <>
                    <Form title={willStore ? "Новая коллекция" : state.name} subtitle={willStore ? "" : "Коллекция"}>
                        <Field name="name" label="Название" icon="category" color={state.color} value={state.name} max={20} error={errors.name} onChange={setField}/>
                        <Palette name="color" value={state.color} onChange={setField}/>
                        <FormButtons>
                            <FormBackButton color={state.color}/>
                            <FormSubmitButton color={state.color} label="Сохранить" onClick={save}/>
                            {!willStore && <FormDeleteButton onClick={() => navigate("./delete")}/>}
                        </FormButtons>
                    </Form>
                    <Outlet/>
                </>
            }
        </div>
    )
}