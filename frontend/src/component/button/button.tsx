import "./button.sass"
import Icon from "../icon/icon.tsx";
import classNames from "classnames";
import {MouseEvent} from "react";

type ButtonProps = {
    label?: string
    className?: string
    icon?: string
    onClick?: () => {}
}

export default function Button({label, className, icon = "", onClick}: ButtonProps) {
    function handleClick(event: MouseEvent<HTMLButtonElement>) {
        event.preventDefault()
        onClick?.()
    }

    return (
        <button className={classNames("button", className, {round: icon && !label})} onClick={handleClick}>
            {icon &&
                <div className="button__icon-container">
                    <Icon className="button__icon" name={icon} filled/>
                </div>
            }
            {label &&
                <div className="button__label">
                    {label}
                </div>
            }
        </button>
    )
}