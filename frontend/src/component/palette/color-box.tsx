import "./palette.sass"
import {Color} from "../../model/color.tsx";
import Icon from "../icon/icon.tsx";
import classNames from "classnames";

type Props = {
    color: Color
    value: number
    setColor: (color: Color) => void
}

export default function ColorBox({color, value, setColor}: Props) {
    const className = classNames(
        "palette__color-box",
        {red: color == Color.Red},
        {orange: color == Color.Orange},
        {yellow: color == Color.Yellow},
        {green: color == Color.Green},
        {blue: color == Color.Blue},
        {purple: color == Color.Purple},
        {selected: value == color}
    )

    return (
        <div onClick={() => setColor(color)} className={className}>
            {value == color && <Icon bold className="palette__check-mark" name="check"/>}
        </div>
    )
}