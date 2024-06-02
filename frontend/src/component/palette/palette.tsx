import "./palette.sass"
import {ChangeEvent} from "react";
import {Color} from "../../model/color";
import ColorBox from "./color-box.tsx";

type Props = {
    name: string
    value: number
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function Palette({name, value, onChange}: Props) {
    function setColor(color: Color) {
        const input = document.getElementById(name) as HTMLInputElement
        input.defaultValue = color.toString()
        input.dispatchEvent(new Event('input', {bubbles: true}))
    }


    return (
        <div className="palette">
            <input type="text" name={name} id={name} className="palette__input" onChange={onChange}/>
            <ColorBox color={Color.Red} value={value} setColor={setColor}/>
            <ColorBox color={Color.Orange} value={value} setColor={setColor}/>
            <ColorBox color={Color.Yellow} value={value} setColor={setColor}/>
            <ColorBox color={Color.Green} value={value} setColor={setColor}/>
            <ColorBox color={Color.Blue} value={value} setColor={setColor}/>
            <ColorBox color={Color.Purple} value={value} setColor={setColor}/>
        </div>
    )
}