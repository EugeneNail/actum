import ShortCollection from "./short-collection.ts";
import {Mood} from "./mood.ts";

export default class ShortRecord {
    id: number = 0
    date: string = ""
    mood: Mood = Mood.Neutral
    notes: string = ""
    collections: ShortCollection[] = []
}