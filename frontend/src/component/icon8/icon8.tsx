import "./icon8.sass"
import classNames from "classnames";
import {icons8Names} from "./icons8.ts";

type Props = {
    className?: string
    id: number
}

export default function Icon8({className, id}: Props) {
    const name = icons8Names[id]

    return (
        <div className={classNames("icon8", className)}>
            <img className="icon8__icon" src={`/img/icons/${name}.png`} alt={name}/>
        </div>
    )
}