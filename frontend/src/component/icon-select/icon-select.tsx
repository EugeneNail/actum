import "./icon-select.sass"
import {ChangeEvent} from "react";
import IconSelectList from "./icon-select-list.tsx";
import {icons} from "../../assets/icons.ts";
import classNames from "classnames";

type Props = {
    className?: string
    name: string
    value: string
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function IconSelect({className, name, value, onChange}: Props) {
    const inputId = "icon-select"


    function setIcon(icon: string) {
        const input = document.getElementById(inputId) as HTMLInputElement
        input.defaultValue = icon
        input.dispatchEvent(new Event('input', {bubbles: true}))
    }


    return (
        <div className={classNames("icon-select", className)}>
            <input id={inputId} className="icon-select__input" name={name} onChange={onChange}/>
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.people} label="People" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.nature} label="Nature" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.home} label="Home" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.food} label="Food & Drinks" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.sport} label="Sport" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.goals} label="Goals" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.transport} label="Transport" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.buildings} label="Buildings" />
            <IconSelectList setIcon={setIcon} selectedIcon={value} icons={icons.miscellaneous} label="Miscellaneous" />
        </div>
    )
}