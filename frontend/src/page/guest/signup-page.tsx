import "./guest-page.sass"
import Field from "../../component/field/field.tsx";
import Button from "../../component/button/button.tsx";
import {Link} from "react-router-dom";

export default function SignupPage() {
    return (
        <div className="guest-page">
            <form className="guest-form">
                <div className="guest-form__logo-container">
                    <div className="guest-form__logo"></div>
                </div>
                <h1 className="guest-form__title">Signup</h1>
                <Field name="name" label="Name" icon="face"/>
                <Field name="email" label="Email" icon="mail"/>
                <Field name="password" label="Password" icon="key"/>
                <Field name="passwordConfirmation" icon="key" label="Confirm password"/>
                <Button label="Sign up" className="guest-form__button"/>
                <p className="guest-form__redirect">Already have an account? <Link className="guest-form__link" to="/login">Login</Link></p>
            </form>
        </div>
    )
}