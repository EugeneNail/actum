import "./guest-page.sass"
import Field from "../../component/field/field.tsx";
import Button from "../../component/button/button.tsx";
import {Link, useNavigate} from "react-router-dom";
import {FormEvent} from "react";
import {useHttp} from "../../service/use-http.ts";
import {useFormState} from "../../service/use-form-state.ts";

class Payload {
    name = ""
    email = ""
    password = ""
    passwordConfirmation = ""
}

class Errors {
    name = ""
    email = ""
    password = ""
    passwordConfirmation = ""
}

export default function SignupPage() {
    const http = useHttp()
    const {state, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const navigate = useNavigate()

    async function submit(event: FormEvent) {
        event.preventDefault()
        const {status, data} = await http.post("/api/users", state)
        setErrors({})

        if (status == 422) {
            setErrors(data)
            return
        }

        if (status == 201) {
            navigate("/")
        }
    }

    return (
        <div className="guest-page">
            <form className="guest-form" onSubmit={submit} method="POST">
                <div className="guest-form__logo-container">
                    <div className="guest-form__logo"></div>
                </div>
                <h1 className="guest-form__title">Signup</h1>
                <Field name="name" label="Name" icon="face" onChange={setField} error={errors.name}/>
                <Field name="email" label="Email" icon="mail" onChange={setField} error={errors.email}/>
                <Field name="password" label="Password" icon="key" password onChange={setField} error={errors.password}/>
                <Field name="passwordConfirmation" label="Confirm password" icon="key" password onChange={setField} error={errors.passwordConfirmation}/>
                <Button label="Sign up" className="guest-form__button"/>
                <p className="guest-form__redirect">Already have an account? <Link className="guest-form__link"                                                            to="/login">Login</Link></p>
            </form>
        </div>
    )
}