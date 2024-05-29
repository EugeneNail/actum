import "./button.sass"
import {ReactNode} from "react";
import classNames from "classnames";

export enum ButtonStyle {
    Primary,
    Secondary
}

type Props = {
    className?:string
    even?: boolean
    round?: boolean
    shadowed?: boolean
    style: ButtonStyle
    onClick: () => void
    children: ReactNode
}

export default function Button({className, even, round, shadowed, style, onClick, children}: Props) {
    className = classNames(
        "button",
        className,
        {even: even},
        {round: round},
        {shadowed: shadowed},
        {primary: style == ButtonStyle.Primary},
        {secondary: style == ButtonStyle.Secondary}
    )

    return (
        <button className={className} onClick={onClick}>
            {children}
        </button>
    )
}