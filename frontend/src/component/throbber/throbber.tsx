import "./throbber.sass"
import classNames from "classnames";

type Props = {
    className?: string
    light?: boolean
}

export default function Throbber({className, light}: Props) {
    return (
        <div className={classNames("throbber", className, {light:light})}>
            <div className="throbber__cover"></div>
        </div>
    )
}