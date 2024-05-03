import classNames from "classnames";
import "./icon.sass"

type IconProps = {
    name: string
    className?: string
    filled?: boolean
    onClick?: () => void
}

export default function Icon({name, className, filled, onClick}: IconProps) {
    return (
        <span className={classNames("icon", "material-symbols-rounded", className, { filled: filled })} onClick={onClick}>
            {name}
        </span>
    )
}