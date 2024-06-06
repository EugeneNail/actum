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

export default function Button({className, color = Color.Green, submit, even, round, shadowed, style = ButtonStyle.Primary, onClick, children}: Props) {
    className = classNames(
        "button",
        className,
        {even: even},
        {round: round},
        {shadowed: shadowed},
        {green: color == Color.Green},
        {yellow: color == Color.Yellow},
        {red: color == Color.Red},
        {blue: color == Color.Blue},
        {orange: color == Color.Orange},
        {purple: color == Color.Purple},
        {accent: color == Color.Accent},
        {primary: style == ButtonStyle.Primary},
        {secondary: style == ButtonStyle.Secondary}
    )

    return (
        <button type={submit ? "submit" : "button"} className={className} onClick={onClick}>
            {children}
        </button>
    )
}