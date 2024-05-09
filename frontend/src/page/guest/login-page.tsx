import "./guest-page.sass"
import Field from "../../component/field/field.tsx";
import Button from "../../component/button/button.tsx";
import {Link, useNavigate} from "react-router-dom";
import {FormEvent} from "react";
import {useHttp} from "../../service/use-http.ts";
import {useFormState} from "../../service/use-form-state.ts";

class Payload {
    email = ""
    password = ""
}

class Errors {
    email = ""
    password = ""
}

export default function LoginPage() {
    const http = useHttp()
    const {state, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const navigate = useNavigate()

    async function submit(event: FormEvent) {
        event.preventDefault()
        const {status, data} = await http.post("/users/login", state)
        setErrors({})

        if (status == 422 || status == 401) {
            setErrors(data)
            return
        }

        if (status == 200) {
            localStorage.setItem("Access-Token", data)
            navigate("/")
        }
    }

    return (
        <div className="guest-page">
            <form className="guest-form" onSubmit={submit} method="POST">
                <div className="guest-form__logo-container">
                    <div className="guest-form__logo"></div>
                </div>
                <h1 className="guest-form__title">Login</h1>
                <Field name="email" label="Email" icon="mail" onChange={setField} error={errors.email}/>
                <Field name="password" label="Password" icon="key" password onChange={setField} error={errors.password}/>
                <Button label="Login" className="guest-form__button"/>
                <p className="guest-form__redirect">Don't have an account? <Link className="guest-form__link" to="/signup">Sign up</Link></p>
            </form>
        </div>
    )
}