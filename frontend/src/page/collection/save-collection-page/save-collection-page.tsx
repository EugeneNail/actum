import "./save-collection-page.sass"
import {useFormState} from "../../../service/use-form-state.ts";
import {useHttp} from "../../../service/use-http.ts";
import {useNavigate, useParams} from "react-router-dom";
import {useEffect} from "react";
import Field from "../../../component/field/field.tsx";
import Button from "../../../component/button/button.tsx";
import FormHeader from "../../../component/form-header/form-header.tsx";

class Payload {
    name = ""
}

class Errors {
    name = ""
}

export default function SaveCollectionPage() {
    const {state, setState, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const {id} = useParams<string>()
    const http = useHttp()
    const navigate = useNavigate()
    const willCreate = window.location.pathname.includes("/new")

    useEffect(() => {
        if (willCreate) {
            return
        }
        fetchCollection()
    }, [])

    async function fetchCollection() {
        const {data, status} = await http.get("/collections/" + id)

        if (status == 403) {
            navigate("/settings/collections")
            return
        }

        setState({
            ...state,
            name: data.name
        })
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
            <div className="cover" onClick={() => navigate("/settings/collections")}/>
            <form className="collection-form" method="POST" onClick={e => e.preventDefault()}>
                <FormHeader icon="category" title={willCreate ? "Add collection" : "Rename collection"}/>
                <div className="collection-form__content">
                    <Field value={state.name} className="collection-form__field" name="name" label="Name" onChange={setField} error={errors.name} icon="category"/>
                    <div className="collection-form__button-container">
                        <Button className="collection-form__button" label={willCreate ? "Add" : "Rename"} pill accent onClick={() => willCreate ? create() : edit()}/>
                        <Button className="collection-form__button cancel" label="Cancel" pill onClick={() => navigate("/settings/collections")}/>
                    </div>
                </div>

            </form>
        </div>
    )
}