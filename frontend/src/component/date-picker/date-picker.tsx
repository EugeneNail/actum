import "./date-picker.sass"
import {ChangeEvent, useEffect, useState} from "react";
import classNames from "classnames";
import Icon from "../icon/icon.tsx";
import Button, {ButtonStyle} from "../button/button.tsx";

type Props = {
    className?: string
    active: boolean
    name: string
    value: string
    error: string
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}


export function DatePicker({className, active, name, value, error, onChange}: Props) {
    const [isVisible, setVisible] = useState(false)
    const [year, setYear] = useState(new Date(value).getFullYear())
    const [month, setMonth] = useState(new Date(value).getMonth())
    const [days, setDays] = useState<number[]>([])
    className = classNames(
        "field",
        "date-picker",
        className,
        {invalid: error?.length > 0}
    )
    const monthNames = ["Января", "Февраля", "Марта", "Апреля", "Мая", "Июня", "Июля", "Августа", "Сентября", "Октября", "Ноября", "Декабря"]


    useEffect(() => {
        if (days.length == 0) {
            buildCalendar(year, month)
        }
    }, [])


    function goTo(back: boolean = false) {
        let destinationMonth = month + (back ? -1 : 1)
        let destinationYear = year

        if (destinationMonth == -1) {
            destinationMonth = 11
            destinationYear--
        }
        if (destinationMonth == 12) {
            destinationYear++
            destinationMonth = 0
        }

        setYear(destinationYear)
        setMonth(destinationMonth)
        buildCalendar(destinationYear, destinationMonth)
    }


    function buildCalendar(year: number, month: number) {
        const days: number[] = [...Array(new Date(year, month, 1).getDay() - 1)]

        let day = 1
        let date = new Date(year, month, day)

        while (date.getMonth() == month) {
            days.push(date.getDate())
            day++
            date = new Date(year, month, day)
        }

        const minCalendarSize = 41
        setDays([...days, ...Array(minCalendarSize - days.length)])
    }


    function setDate(year: number, month: number, day: number) {
        const input = document.getElementById(name) as HTMLInputElement
        input.defaultValue = new Date(year, month, day+ 1).toISOString().split("T")[0]
        input.dispatchEvent(new Event('input', {bubbles: true}))
        toggleCalendar()
    }


    function toggleCalendar() {
        if (isVisible) {
            setVisible(false)
            document.body.style.overflow = "auto"
        } else {
            setVisible(true)
            document.body.style.overflow = "hidden"
        }
    }


    function getDate(): string {
        const date = new Date(value)

        if (value == new Date().toISOString().split("T")[0]) {
            return `Сегодня, ${dateToUserFriendly(date)}`
        }

        const yesterday = new Date();
        yesterday.setDate(yesterday.getDate() - 1)
        if (value == yesterday.toISOString().split("T")[0]) {
            return `Вчера, ${dateToUserFriendly(date)}`
        }

        return dateToUserFriendly(date)
    }


    function dateToUserFriendly(date: Date): string {
        return `${date.getDate()} ${monthNames[date.getMonth()]} ${date.getFullYear()}`
    }


    function checkUnavailability(year: number, month: number, day: number): boolean {
        const date = new Date(year, month, day)

        return date > new Date()
    }


    function checkSelection(year: number, month: number, day: number): boolean {
        return value == `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`
    }


    return (
        <div className={className}>
            <div className={classNames("date-picker__label", {active: active})} onClick={toggleCalendar}>
                <Icon className="date-picker__icon" name="event"/>
                <p className="date-picker__date">{getDate()}</p>
            </div>
            <input type="text" id={name} name={name} className="date-picker__input" onChange={onChange}/>
            {isVisible && <div className="date-picker__cover" onClick={toggleCalendar}></div>}
            {isVisible && <div className="date-picker__calendar">
                <div className="date-picker__calendar-header">
                    <p className="date-picker__calendar-label">
                        {monthNames[month]} {year}
                    </p>
                    <Button className="date-picker__button" style={ButtonStyle.Secondary} even onClick={() => goTo(true)}>
                        <Icon bold name="west"/>
                    </Button>
                    <Button className="date-picker__button" style={ButtonStyle.Secondary} even onClick={() => goTo()}>
                        <Icon bold name="east"/>
                    </Button>
                </div>

                <div className="date-picker__days">
                    <div className="date-picker__day weekday">Пн</div>
                    <div className="date-picker__day weekday">Вт</div>
                    <div className="date-picker__day weekday">Ср</div>
                    <div className="date-picker__day weekday">Чт</div>
                    <div className="date-picker__day weekday">Пт</div>
                    <div className="date-picker__day weekday">Сб</div>
                    <div className="date-picker__day weekday">Вс</div>
                    {days && days.map(day => (
                        <div key={Math.random()}
                            className={classNames(
                            "date-picker__day",
                            {filler: day == undefined},
                            {selected: checkSelection(year, month, day)},
                            {unavailable: checkUnavailability(year, month, day)}
                        )}
                             onClick={() => setDate(year, month, day)}>
                            {day ?? 0}
                        </div>
                    ))}
                </div>
            </div>}
        </div>
    )
}