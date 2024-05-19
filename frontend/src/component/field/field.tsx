import Icon from "../icon/icon.tsx";
import classNames from "classnames";
import "./field.sass"
import {ChangeEvent, useEffect, useState} from "react";

type FieldProps = {
    value: string
    icon?: string
    name: string
    label: string
    className?: string
    error?: string
    password?: boolean
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function Field({value, icon = "", name, label, className, error, password, onChange}: FieldProps) {
    const [isVisible, setVisible] = useState(true)

    useEffect(() => {
        if (password) {
            setVisible(false)
        }
    }, [])

    return (
        <div className="field-wrapper">
            <div className={classNames("field", className)}>
                <Icon name={icon}/>
                <div className="field__input-container">
                    <input value={value} placeholder="" type={isVisible ? "text" : "password"} id={name} name={name} className="field__input" onChange={onChange}/>
                    <label htmlFor={name} className="field__label">{label}</label>
                </div>
                {password && <Icon className="field__visibility" name={isVisible ? "visibility_off" : "visibility"} onClick={() => setVisible(!isVisible)}/>}
            </div>
            {error && <p className="field__error">{error}</p>}
        </div>
    )
}