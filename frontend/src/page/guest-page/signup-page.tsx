import "./guest-page.sass"
import Form from "../../component/form/form.tsx";
import Field from "../../component/field/field.tsx";
import {useFormState} from "../../service/use-form-state.ts";
import {useHttp} from "../../service/use-http.ts";
import {Link, useNavigate} from "react-router-dom";
import base64UrlToString from "../../service/base64.ts";
import FormButtons from "../../component/form/form-button-container.tsx";
import FormSubmitButton from "../../component/form/form-submit-button.tsx";

class Payload {
    name: string = ""
    email: string = ""
    password: string = ""
    passwordConfirmation = ""
}

class Errors {
    name: string = ""
    email: string = ""
    password: string = ""
    passwordConfirmation = ""
}

export default function SignupPage() {
    const {state, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const http = useHttp()
    const navigate = useNavigate()

    async function signup() {
        const {data, status} = await http.post("/api/users", state)
        if (status == 422) {
            setErrors(data)
            return
        }

        if (status == 201) {
            const decoded = base64UrlToString(data.split(".")[1])
            const payload = JSON.parse(decoded)
            localStorage.setItem("username", payload.name)
            localStorage.setItem("Access-Token", data)
            navigate("/")
        }
    }

    return (
        <div className="page">
            <Form title="Sign up" subtitle={"to start working"}>
                <Field name="name" label="What is your name?" icon="face" value={state.name} error={errors.name} onChange={setField}/>
                <Field name="email" label="Email" icon="mail" value={state.email} error={errors.email} onChange={setField}/>
                <Field name="password" label="Password" icon="lock" value={state.password} error={errors.password} onChange={setField} password/>
                <Field name="passwordConfirmation" label="Confirm password" icon="lock" value={state.passwordConfirmation} error={errors.passwordConfirmation} onChange={setField} password/>
                <FormButtons>
                    <FormSubmitButton label="Sign up" onClick={signup}/>
                </FormButtons>
            </Form>
            <Link to="/login" className="guest-page-link">I already have an account</Link>
        </div>
    )
}
