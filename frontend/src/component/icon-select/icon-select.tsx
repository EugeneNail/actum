import "./icon-select.sass"
import {ChangeEvent} from "react";
import IconSelectList from "./icon-select-list.tsx";
import classNames from "classnames";
import {Color} from "../../model/color.tsx";

type Props = {
    className?: string
    name: string
    value: number
    color: Color
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function IconSelect({className, name, value, color, onChange}: Props) {
    const inputId = "icon-select"


    function setIcon(icon: number) {
        const input = document.getElementById(inputId) as HTMLInputElement
        input.defaultValue = icon.toString()
        input.dispatchEvent(new Event('input', {bubbles: true}))
    }


    return (
        <div className={classNames("icon-select", className)}>
            <input id={inputId} className="icon-select__input" name={name} onChange={onChange}/>
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={100} color={color} label="Люди" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={200} color={color} label="Животные и Насекомые" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={300} color={color} label="Еда и Напитки" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={400} color={color} label="Природа" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={500} color={color} label="Спорт" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={600} color={color} label="Места и Путешествия" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={700} color={color} label="Дом и Двор" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={800} color={color} label="Тело" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={900} color={color} label="Красота и Мода" />
        </div>
    )
}