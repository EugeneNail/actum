import classNames from "classnames";
import "./icon.sass"

type IconProps = {
    name: string
    className?: string
    bold?: boolean
    onClick?: () => void
}

export default function Icon({name, className, bold, onClick}: IconProps) {
    return (
        <span className={classNames("icon", "material-symbols-rounded", className, { bold: bold })} onClick={onClick}>
            {name}
        </span>
    )
}