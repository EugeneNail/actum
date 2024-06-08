import "./save-collection-page.sass"
import {useHttp} from "../../../service/use-http.ts";
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

class Payload {
    name = ""
    color: Color = Color.Red
}

class Errors {
    name = ""
    color = ""
}

export default function SaveCollectionPage() {
    const http = useHttp()
    const {state, setState, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const willStore = window.location.pathname.includes("/new")
    const navigate = useNavigate()
    const notification = useNotificationContext()
    const {id} = useParams()
    const [initialName, setInitialName] = useState("")


    useEffect(() => {
        document.title = "Новая коллекция"
        if (!willStore) {
            fetchCollection()
        }
    }, [])


    async function fetchCollection() {
        const {data, status} = await http.get(`/api/collections/${id}`)

        if (status == 403) {
            notification.pop(data)
            navigate("/collections")
            return
        }

        if (status == 200) {
            setInitialName(data.name)
            setState({
                name: data.name,
                color: data.color
            })
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
        const {data, status} = await http.post("/api/collections", {
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
        const {data, status} = await http.put(`/api/collections/${id}`, {
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
            <Form title={willStore ? "Новая коллекция" : "Коллекция"} subtitle={initialName ? initialName : ""}>
                <Field name="name" label="Название" icon="category" value={state.name} error={errors.name} onChange={setField}/>
                <Palette name="color" value={state.color} onChange={setField}/>
                <FormButtons>
                    <FormBackButton/>
                    <FormSubmitButton label="Сохранить" onClick={save}/>
                    {!willStore && <FormDeleteButton onClick={() => navigate("./delete")}/>}
                </FormButtons>
            </Form>
            <Outlet/>
        </div>
    )
}