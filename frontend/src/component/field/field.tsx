import Icon from "../icon/icon.tsx";
import classNames from "classnames";
import "./field.sass"

type FieldProps = {
    icon?: string
    name: string
    label: string
    className?: string
    error?: string
}

export default function Field({icon = "", name, label, className, error}: FieldProps) {
    return (
        <div className="field-wrapper">
            <div className={classNames("field", className)}>
                <div className="field__icon-container">
                    <Icon name={icon}/>
                </div>
                <div className="field__input-container">
                    <input placeholder="" type="text" id={name} name={name} className="field__input"/>
                    <label htmlFor={name} className="field__label">{label}</label>
                </div>
            </div>
            {error && <p className="field__error">{error}</p>}
        </div>
    )
}