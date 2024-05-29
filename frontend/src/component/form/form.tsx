import "./form.sass"
import {ReactNode} from "react";
import Button, {ButtonStyle} from "../button/button.tsx";
import Icon from "../icon/icon.tsx";

type Props = {
    title: string
    subtitle?: string
    submitMessage: string
    noBack?:boolean
    onSubmit: () => void
    children: ReactNode
}

export default function Form({title, subtitle, submitMessage, noBack, onSubmit, children}: Props) {
    return (
        <form className="form" onSubmit={e => e.preventDefault()}>
            <h1 className="form__title">{title}</h1>
            {subtitle && <p className="form__subtitle">{subtitle}</p>}
            <div className="form__content-container">
                {children}
            </div>
            <div className="form__button-container">
                {!noBack && <Button even style={ButtonStyle.Secondary} onClick={() => window.history.back()} >
                    <Icon name="west" bold/>
                </Button>}

                <Button className="form__submit-button" style={ButtonStyle.Primary} onClick={onSubmit} >
                    {submitMessage} <Icon name="east" bold/>
                </Button>
            </div>
        </form>
    )
}