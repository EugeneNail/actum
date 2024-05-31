import "./form.sass"
import {ReactNode} from "react";

type Props = {
    title: string
    subtitle?: string
    children: ReactNode
}

export default function Form({title, subtitle, children}: Props) {
    return (
        <form className="form" onSubmit={e => e.preventDefault()}>
            <h1 className="form__title">{title}</h1>
            {subtitle && <p className="form__subtitle">{subtitle}</p>}
            <div className="form__content-container">
                {children}
            </div>
        </form>
    )
}