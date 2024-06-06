import "./icon-select.sass"
import {ChangeEvent} from "react";
import IconSelectList from "./icon-select-list.tsx";
import classNames from "classnames";

type Props = {
    className?: string
    name: string
    value: number
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function IconSelect({className, name, value, onChange}: Props) {
    const inputId = "icon-select"


    function setIcon(icon: number) {
        const input = document.getElementById(inputId) as HTMLInputElement
        input.defaultValue = icon.toString()
        input.dispatchEvent(new Event('input', {bubbles: true}))
    }


    return (
        <div className={classNames("icon-select", className)}>
            <input id={inputId} className="icon-select__input" name={name} onChange={onChange}/>
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={100} label="Люди" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={200} label="Животные и Насекомые" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={300} label="Еда и Напитки" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={400} label="Природа" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={500} label="Спорт" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={600} label="Места и Путешествия" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={700} label="Дом и Двор" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={800} label="Тело" />
            <IconSelectList setIcon={setIcon} selectedIconId={value} group={900} label="Красота и Мода" />
        </div>
    )
}