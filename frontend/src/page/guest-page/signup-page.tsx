import "./guest-page.sass"
import Form from "../../component/form/form.tsx";
import Field from "../../component/field/field.tsx";
import {useFormState} from "../../service/use-form-state.ts";
import {useApi} from "../../service/use-api.ts";
import {Link, useNavigate} from "react-router-dom";
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
    const api = useApi()
    const navigate = useNavigate()
    document.title = "Регистрация - Actum"

    async function signup() {
        const {data, status} = await api.post("/api/users", state)
        if (status == 422) {
            setErrors(data)
            return
        }

        if (status == 201) {
            localStorage.setItem("Access-Token", data.access)
            localStorage.setItem("Refresh-Token", data.refresh)
            navigate("/")
        }
    }

    return (
        <div className="page">
            <Form title="Регистрация" subtitle={"это только начало"}>
                <Field name="name" label="Как вас зовут?" icon="face" value={state.name} max={20} error={errors.name} onChange={setField}/>
                <Field name="email" label="Электронная почта" icon="mail" value={state.email} email max={100} error={errors.email} onChange={setField}/>
                <Field name="password" label="Пароль" icon="lock" value={state.password} max={100} error={errors.password} onChange={setField} password/>
                <Field name="passwordConfirmation" label="Повторите пароль" icon="lock" value={state.passwordConfirmation} max={100} error={errors.passwordConfirmation} onChange={setField} password/>
                <FormButtons>
                    <FormSubmitButton label="Зарегистрироваться" onClick={signup}/>
                </FormButtons>
            </Form>
            <Link to="/login" className="guest-page-link">У меня уже есть аккаунт</Link>
        </div>
    )
}
