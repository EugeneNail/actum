import "./button.sass"
import Icon from "../icon/icon.tsx";
import classNames from "classnames";

type ButtonProps = {
    label?: string
    className?: string
    icon?: string
    onClick?: () => void
    negative?: boolean
    pill?: boolean
}

export default function Button({label, className, icon = "", negative, pill, onClick}: ButtonProps) {
    return (
        <button className={classNames("button", className, {round: icon && !label}, {negative: negative}, {pill: pill})} onClick={onClick}>
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