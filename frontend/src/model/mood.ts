export enum Mood {
    Awful = 1,
    Bad = 2,
    Neutral = 3,
    Happy = 4,
    Radiating = 5
}

export const MoodIcons: {[id: number]: string} = {
    [Mood.Awful]: "sentiment_very_dissatisfied",
    [Mood.Bad]: "sentiment_dissatisfied",
    [Mood.Neutral]: "sentiment_neutral",
    [Mood.Happy]: "sentiment_satisfied",
    [Mood.Radiating]: "sentiment_very_satisfied",
}