import classNames from "classnames";
import "./icon.sass"

type IconProps = {
    name: string
    className?: string
    filled?: boolean
}

export default function Icon({name, className, filled = false}: IconProps) {
    return (
        <span className={classNames("icon", "material-symbols-rounded", className, { filled: filled })}>
            {name}
        </span>
    )
}