import "./weather-select.sass"
import {ChangeEvent} from "react";
import {Weather, WeatherIcons} from "../../model/weather.ts";
import Icon from "../icon/icon.tsx";
import classNames from "classnames";

type Props = {
    name: string
    value: Weather
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function WeatherSelect({name, value, onChange}: Props) {
    function setWeather(weather: Weather) {
        const input = document.getElementById(name) as HTMLInputElement
        input.defaultValue = weather.toString()
        input.dispatchEvent(new Event('input', {bubbles: true}))
    }

    function getWeatherClassName(weather: Weather) {
        return classNames(
            "weather-select__weather",
            {selected: value == weather},
            {heat: weather == Weather.Heat},
            {sunny: weather == Weather.Sunny},
            {cloudy: weather == Weather.Cloudy},
            {windy: weather == Weather.Windy},
            {rainy: weather == Weather.Rainy},
            {thunder: weather == Weather.Thunder},
            {foggy: weather == Weather.Foggy},
            {snowy: weather == Weather.Snowy},
            {cold: weather == Weather.Cold},
        )
    }

    return (
        <div className="weather-select">
            <input name={name} id={name} type="text" className="weather-select__input" onChange={onChange}/>
            <p className="weather-select__label">Как погода?</p>
            <div className="weather-select__weathers">
                <Icon className={getWeatherClassName(Weather.Heat)} name={WeatherIcons[Weather.Heat]} onClick={() => setWeather(Weather.Heat)} />
                <Icon className={getWeatherClassName(Weather.Sunny)} name={WeatherIcons[Weather.Sunny]} onClick={() => setWeather(Weather.Sunny)} />
                <Icon className={getWeatherClassName(Weather.Cloudy)} name={WeatherIcons[Weather.Cloudy]} onClick={() => setWeather(Weather.Cloudy)} />
                <Icon className={getWeatherClassName(Weather.Windy)} name={WeatherIcons[Weather.Windy]} onClick={() => setWeather(Weather.Windy)} />
                <Icon className={getWeatherClassName(Weather.Rainy)} name={WeatherIcons[Weather.Rainy]} onClick={() => setWeather(Weather.Rainy)} />
                <Icon className={getWeatherClassName(Weather.Thunder)} name={WeatherIcons[Weather.Thunder]} onClick={() => setWeather(Weather.Thunder)} />
                <Icon className={getWeatherClassName(Weather.Foggy)} name={WeatherIcons[Weather.Foggy]} onClick={() => setWeather(Weather.Foggy)} />
                <Icon className={getWeatherClassName(Weather.Snowy)} name={WeatherIcons[Weather.Snowy]} onClick={() => setWeather(Weather.Snowy)} />
                <Icon className={getWeatherClassName(Weather.Cold)} name={WeatherIcons[Weather.Cold]} onClick={() => setWeather(Weather.Cold)} />
            </div>
        </div>
    )
}