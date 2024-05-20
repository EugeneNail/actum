import Activity from "./activity.ts";

export default class Collection {
    id: number = 0
    name: string = ""
    userId: number = 0
    activities: Activity[] = []
}