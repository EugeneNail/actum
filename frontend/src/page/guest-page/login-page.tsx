import "./guest-page.sass"
import Form from "../../component/form/form.tsx";
import Field from "../../component/field/field.tsx";
import {useFormState} from "../../service/use-form-state.ts";
import {useApi} from "../../service/use-api.ts";
import {Link, useNavigate} from "react-router-dom";
import FormButtons from "../../component/form/form-button-container.tsx";
import FormSubmitButton from "../../component/form/form-submit-button.tsx";

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
    const api = useApi()
    const navigate = useNavigate()
    document.title = "Вход - Actum"

    async function login() {
        const {data, status} = await api.post("/api/users/login", state)
        if (status == 422 || status == 401) {
            setErrors(data)
            return
        }

        if (status == 200) {
            localStorage.setItem("Access-Token", data.access)
            localStorage.setItem("Refresh-Token", data.refresh)
            navigate("/")
        }
    }

    return (
        <div className="page">
            <Form title="Привет!" subtitle={"Войдите, чтобы продолжить"}>
                <Field name="email" label="Электронная почта" icon="mail" value={state.email} email max={100} error={errors.email} onChange={setField}/>
                <Field name="password" label="Пароль" icon="lock" value={state.password} max={100} error={errors.password} onChange={setField} password/>
                <FormButtons>
                    <FormSubmitButton label="Войти" onClick={login}/>
                </FormButtons>
            </Form>
            <Link to="/signup" className="guest-page-link">У меня нет аккаунта</Link>
        </div>
    )
}
