import "./button.sass"
import {ReactNode} from "react";
import classNames from "classnames";
import {Color} from "../../model/color.tsx";

export enum ButtonStyle {
    Primary,
    Secondary
}

type Props = {
    className?:string
    color?: Color
    submit?: boolean
    even?: boolean
    round?: boolean
    shadowed?: boolean
    style?: ButtonStyle
    onClick: () => void
    children: ReactNode
}

export default function Button({className, color = Color.green, submit, even, round, shadowed, style = ButtonStyle.Primary, onClick, children}: Props) {
    className = classNames(
        "button",
        className,
        {even: even},
        {round: round},
        {shadowed: shadowed},
        {green: color == Color.green},
        {yellow: color == Color.yellow},
        {red: color == Color.red},
        {primary: style == ButtonStyle.Primary},
        {secondary: style == ButtonStyle.Secondary}
    )

    return (
        <button type={submit ? "submit" : "button"} className={className} onClick={onClick}>
            {children}
        </button>
    )
}