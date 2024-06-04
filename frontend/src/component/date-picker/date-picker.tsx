import "./date-picker.sass"
import {ChangeEvent, useEffect, useState} from "react";
import classNames from "classnames";
import Icon from "../icon/icon.tsx";
import buildCalendar, {monthNames, Year} from "./build-calendar.ts";

type Props = {
    className?: string
    label: string
    name: string
    value: string
    error: string
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}


export function DatePicker({className, label, name, value, error, onChange}: Props) {
    const [isCalendarVisible, setCalendarVisible] = useState(false)
    const [calendar, setCalendar] = useState<Year[]>([])
    className = classNames(
        "field",
        "date-picker",
        className,
        {invalid: error?.length > 0}
    )


    useEffect(() => {
        if (calendar?.length == 0) {
            setCalendar(buildCalendar)
        }
    }, [])


    function toggleCalendar() {
        const form = document.getElementsByTagName("form")[0] as HTMLFormElement
        form.style.overflow = "hidden"
        if (isCalendarVisible) {
            setCalendarVisible(!isCalendarVisible)
            form.style.height = "fit-content"
        } else {
            form.style.height = "0px"
        }
        setCalendarVisible(!isCalendarVisible)
    }


    function checkSelection(year: number, month: number, day: number): boolean {
        return value == `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`
    }


    function setDate(year: number, month: number, day: number) {
        const input = document.getElementById(name) as HTMLInputElement
        input.defaultValue = `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`
        input.dispatchEvent(new Event('input', {bubbles: true}))
        toggleCalendar()
    }


    return (
        <div className={className}>
            {isCalendarVisible && <div className="date-picker__calendar">
                {calendar && calendar.map(year =>
                    year.months.map(month => (
                        <div key={Math.random()} className="date-picker__month">
                            <h6 className="date-picker__month-name">{monthNames[month.value]} {year.value}</h6>
                            <div className="date-picker__month-days">
                                <div className="date-picker__day weekday">Su</div>
                                <div className="date-picker__day weekday">Mo</div>
                                <div className="date-picker__day weekday">Tu</div>
                                <div className="date-picker__day weekday">We</div>
                                <div className="date-picker__day weekday">Th</div>
                                <div className="date-picker__day weekday">Fr</div>
                                <div className="date-picker__day weekday">Sa</div>
                                {month.days && month.days.map(day => (
                                    <div key={Math.random()} className={classNames("date-picker__day", {inactive: day == undefined}, {selected: checkSelection(year.value, month.value, day)})} onClick={() => setDate(year.value, month.value, day)}>{day ?? 0}</div>
                                ))}
                            </div>
                        </div>
                    )))
                }
            </div>
            }
            <div className="field__content">
                <div className="field__icon-container">
                    <Icon name="event"/>
                </div>
                <input placeholder={label} onFocus={e => e.target.blur()} type="text" id={name} name={name} className="field__input date-picker__input" onChange={onChange} onClick={toggleCalendar}/>
            </div>
            <p className="field__error">{error}</p>
        </div>
    )





}