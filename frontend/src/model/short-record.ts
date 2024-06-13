import ShortCollection from "./short-collection.ts";
import {Mood} from "./mood.ts";
import {Weather} from "./weather.ts";

export default class ShortRecord {
    id: number = 0
    date: string = ""
    mood: Mood = Mood.Neutral
    weather: Weather = Weather.Sunny
    notes: string = ""
    collections: ShortCollection[] = []
    photos: string[] = []
}