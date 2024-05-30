import "./form.sass"
import {ReactNode} from "react";

type Props = {
    children: ReactNode
}

export default function FormButtons({children}: Props) {
    return (
        <div className="form__button-container">
            {children}
        </div>
    )
}