import Icon from "../icon/icon.tsx";
import classNames from "classnames";
import "./field.sass"
import {ChangeEvent, useEffect, useState} from "react";

type FieldProps = {
    value: string
    icon?: string
    name: string
    label: string
    max?: number
    className?: string
    error?: string
    email?: boolean
    password?: boolean
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function Field({value, icon = "", name, label, max = 100, className, error = "", email, password, onChange}: FieldProps) {
    const [isVisible, setVisible] = useState(true)

    className = classNames(
        "field",
        className,
        {invalid: error?.length > 0}
    )

    useEffect(() => {
        if (password) {
            setVisible(false)
        }
    }, [])

    function getInputType(): string {
        if (email && !password) {
            return "email"
        }

        if (password) {
            return isVisible ? "text" : "password"
        }

        return "text"
    }

    return (
        <div className={className}>
            <div className={classNames("field__content", className)}>
                <div className="field__icon-container">
                    <Icon name={icon}/>
                </div>
                <input autoComplete="off"
                       value={value}
                       placeholder={label}
                       maxLength={max}
                       type={getInputType()}
                       id={name}
                       name={name}
                       inputMode={email ? "email" : "text"}
                       className="field__input"
                       onChange={onChange}/>
                {password && <Icon className="field__visibility" name={isVisible ? "visibility_off" : "visibility"} onClick={() => setVisible(!isVisible)}/>}
            </div>
            <p className="field__error">{error}</p>
        </div>
    )
}