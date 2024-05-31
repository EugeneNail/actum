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

export default function Field({value, icon = "", name, label, className, error = "", password, onChange}: FieldProps) {
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

    return (
        <div className={className}>
            <div className={classNames("field__content", className)}>
                <div className="field__icon-container">
                    <Icon name={icon}/>
                </div>
                <input autoComplete="off" value={value} placeholder={label} type={isVisible ? "text" : "password"} id={name} name={name} className="field__input" onChange={onChange}/>
                {password && <Icon className="field__visibility" name={isVisible ? "visibility_off" : "visibility"} onClick={() => setVisible(!isVisible)}/>}
            </div>
            <p className="field__error">{error}</p>
        </div>
    )
}