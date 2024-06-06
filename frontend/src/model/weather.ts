export enum Weather {
    Heat = 1,
    Sunny = 2,
    Cloudy = 3,
    Windy = 4,
    Rainy = 5,
    Thunder = 6,
    Foggy = 7,
    Snowy = 8,
    Cold = 9
}

export const WeatherIcons: {[id: number]: string} = {
    [Weather.Heat]: "heat",
    [Weather.Sunny]: "sunny",
    [Weather.Cloudy]: "cloud",
    [Weather.Windy]: "air",
    [Weather.Rainy]: "rainy",
    [Weather.Thunder]: "thunderstorm",
    [Weather.Foggy]: "foggy",
    [Weather.Snowy]: "cloudy_snowing",
    [Weather.Cold]: "ac_unit",
}

export const WeatherNames: {[id: number]: string} = {
    [Weather.Heat]: "Зной",
    [Weather.Sunny]: "Солнечно",
    [Weather.Cloudy]: "Пасмурно",
    [Weather.Windy]: "Ветрено",
    [Weather.Rainy]: "Дождь",
    [Weather.Thunder]: "Гроза",
    [Weather.Foggy]: "Туман",
    [Weather.Snowy]: "Снег",
    [Weather.Cold]: "Мороз",
}