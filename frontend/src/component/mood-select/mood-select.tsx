import "./mood-select.sass"
import {ChangeEvent} from "react";
import Icon from "../icon/icon.tsx";
import classNames from "classnames";

type Props = {
    className?: string
    name: string
    value: number
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function MoodSelect({className, name, value, onChange}: Props) {


    function setMood(mood: number) {
        const input = document.getElementById(name) as HTMLInputElement
        input.defaultValue = mood.toString()
        input.dispatchEvent(new Event('input', {bubbles: true}))

    }

    return (
        <div className={classNames("mood-select", className)}>
            <input type="number" name={name} id={name} className="mood-select__input" onChange={onChange}/>
            <p className="mood-select__label">Как вы себя чувствуете?</p>
            <div className="mood-select__moods">
                <Icon className={classNames("mood-select__mood", {selected: value == 5})} name="sentiment_very_satisfied" onClick={() => setMood(5)} />
                <Icon className={classNames("mood-select__mood", {selected: value == 4})} name="sentiment_satisfied" onClick={() => setMood(4)} />
                <Icon className={classNames("mood-select__mood", {selected: value == 3})} name="sentiment_neutral" onClick={() => setMood(3)} />
                <Icon className={classNames("mood-select__mood", {selected: value == 2})} name="sentiment_dissatisfied" onClick={() => setMood(2)} />
                <Icon className={classNames("mood-select__mood", {selected: value == 1})} name="sentiment_very_dissatisfied" onClick={() => setMood(1)} />
            </div>
        </div>
    )
}