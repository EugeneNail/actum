import "./save-collection-page.sass"
import {useFormState} from "../../../service/use-form-state.ts";
import {useHttp} from "../../../service/use-http.ts";
import {useNavigate} from "react-router-dom";
import {FormEvent} from "react";
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
    const http = useHttp()
    const navigate = useNavigate()

    async function submit(event: FormEvent) {
        event.preventDefault()

        const {data, status} = await http.post("/collections", state)

        if (status == 422 || status == 409) {
            setErrors(data)
            return
        }

        navigate("/collections")
    }

    return (
        <div className="save-collection-page">
            <form className="collection-form" method="POST" onSubmit={submit}>
                <Field className="collection-form__field" name="name" label="Name" onChange={setField}
                       error={errors.name} icon="category"/>
                <Button className="collection-form__button" icon="add" label="Create collection"/>
            </form>
        </div>
    )
}