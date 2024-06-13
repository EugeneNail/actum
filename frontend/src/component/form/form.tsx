import "./form.sass"
import {ReactNode} from "react";
import classNames from "classnames";

type Props = {
    title: string
    subtitle?: string
    noBackground?: boolean
    children: ReactNode
}

export default function Form({title, subtitle, noBackground, children}: Props) {
    return (
        <form className={classNames("form", {"no-background": noBackground})} onSubmit={e => e.preventDefault()}>
            <h1 className="form__title">{title}</h1>
            {subtitle && <p className="form__subtitle">{subtitle}</p>}
            <div className="form__content-container">
                {children}
            </div>
        </form>
    )
}