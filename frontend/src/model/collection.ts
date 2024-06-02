import Activity from "./activity.ts";
import {Color} from "./color.tsx";

export default class Collection {
    id: number = 0
    name: string = ""
    userId: number = 0
    color: Color = Color.Red
    activities: Activity[] = []
}