import Form from "../../component/form/form.tsx";
import Field from "../../component/field/field.tsx";
import {useFormState} from "../../service/use-form-state.ts";
import {Color} from "../../model/color.tsx";
import {useHttp} from "../../service/use-http.ts";
import {useNavigate} from "react-router-dom";
import base64UrlToString from "../../service/base64.ts";

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

    async function submit() {
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
            <Form title="Sign up" subtitle={"to start working"} submitMessage="Sign up" noBack onSubmit={submit}>
                <Field color={Color.green} name="name" label="What is your name?" icon="face" value={state.name} error={errors.name} onChange={setField}/>
                <Field color={Color.green} name="email" label="Email" icon="mail" value={state.email} error={errors.email} onChange={setField}/>
                <Field color={Color.red} name="password" label="Password" icon="lock" value={state.password} error={errors.password} onChange={setField} password/>
                <Field color={Color.red} name="passwordConfirmation" label="Confirm password" icon="lock" value={state.passwordConfirmation} error={errors.passwordConfirmation} onChange={setField} password/>
            </Form>
        </div>
    )
}
