export type Year = {
    value: number,
    months: Month[]
}

export type Month = {
    value: number,
    days: number[]
}

export default function buildCalendar(): Year[] {
    const today = new Date().getDate()
    let yearIndex = -1
    let monthIndex = 0
    const calendar: Year[] = []

    for (let daysAgo = 365; daysAgo >= 0; daysAgo--) {
        const date = new Date()
        date.setDate(today-daysAgo)

        const year = date.getFullYear()
        const month = date.getMonth()
        const day = date.getDate()

        if (calendar.length == 0 || calendar[yearIndex].value != year) {
            calendar.push({
                value: year,
                months: []
            })
            yearIndex++
            monthIndex = -1
        }

        const months = calendar[yearIndex].months
        if (months.length == 0 || months[monthIndex].value != month) {
            months.push({
                value: month,
                days: [...Array(date.getDay())]
            })
            monthIndex++
        }

        months[monthIndex].days.push(day)
    }

    calendar[0].months.reverse()
    calendar[1]?.months.reverse()

    return calendar.reverse()
}

export const monthNames: {[name: number]: string} = {
    0: "January",
    1: "February",
    2: "March",
    3: "April",
    4: "May",
    5: "June",
    6: "July",
    7: "August",
    8: "September",
    9: "October",
    10: "November",
    11: "December",
}