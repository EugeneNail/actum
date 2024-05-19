import "./save-collection-page.sass"
import {useFormState} from "../../../service/use-form-state.ts";
import {useHttp} from "../../../service/use-http.ts";
import {useNavigate, useParams} from "react-router-dom";
import {FormEvent, useEffect} from "react";
import Field from "../../../component/field/field.tsx";
import Button from "../../../component/button/button.tsx";

class Payload {
    name = ""
}

class Errors {
    name = ""
}

export default function SaveCollectionPage() {
    const {state, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const {id} = useParams<string>()
    const http = useHttp()
    const navigate = useNavigate()
    const willCreate = window.location.pathname.includes("/new")

    async function submit(event: FormEvent) {
        event.preventDefault()

        if (willCreate) {
            await create()
        } else {
            await edit()
        }

    }

    async function create() {
        const {data, status} = await http.post("/collections", state)

        if (status == 422 || status == 409) {
            setErrors(data)
            return
        }

        navigate("/settings/collections")
    }

    async function edit() {
        const {data, status} = await http.put("/collections/" + id, state)

        if (status == 422) {
            setErrors(data)
            return
        }

        navigate("/settings/collections")
    }

    return (
        <div className="save-collection-page">
            <form className="collection-form" method="POST" onSubmit={submit}>
                <Field value={state.name} className="collection-form__field" name="name" label="Name" onChange={setField}
                       error={errors.name} icon="category"/>
                <Button className="collection-form__button" icon={willCreate ? "add" : "edit"} label={willCreate ? "Create collection" : "Rename collection"}/>
            </form>
        </div>
    )
}