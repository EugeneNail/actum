import "./guest-page.sass"
import Form from "../../component/form/form.tsx";
import Field from "../../component/field/field.tsx";
import {useFormState} from "../../service/use-form-state.ts";
import {Color} from "../../model/color.tsx";
import {useHttp} from "../../service/use-http.ts";
import {Link, useNavigate} from "react-router-dom";
import base64UrlToString from "../../service/base64.ts";

class Payload {
    email: string = ""
    password: string = ""
}

class Errors {
    email: string = ""
    password: string = ""
}

export default function LoginPage() {
    const {state, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const http = useHttp()
    const navigate = useNavigate()

    async function submit() {
        const {data, status} = await http.post("/api/users/login", state)
        if (status == 422 || status == 401) {
            setErrors(data)
            return
        }

        if (status == 200) {
            const decoded = base64UrlToString(data.split(".")[1])
            const payload = JSON.parse(decoded)
            localStorage.setItem("username", payload.name)
            localStorage.setItem("Access-Token", data)
            navigate("/")
        }
    }

    const username = localStorage.getItem("username")
    const greetings = username != null ? `Hello, ${username}!` : "Hello!"

    return (
        <div className="page">
            <Form title={greetings} subtitle={"Log in to continue"} submitMessage="Log in" noBack onSubmit={submit}>
                <Field color={Color.green} name="email" label="Email" icon="mail" value={state.email} error={errors.email} onChange={setField}/>
                <Field color={Color.red} name="password" label="Password" icon="lock" value={state.password} error={errors.password} onChange={setField} password/>
            </Form>
            <Link to="/signup" className="guest-page-link">I don't have an account</Link>
        </div>
    )
}
