import "./save-record-page.sass"
import Form from "../../../component/form/form.tsx";
import FormButtons from "../../../component/form/form-button-container.tsx";
import FormBackButton from "../../../component/form/form-back-button.tsx";
import FormSubmitButton from "../../../component/form/form-submit-button.tsx";
import {useFormState} from "../../../service/use-form-state.ts";
import {useEffect, useState} from "react";
import Collection from "../../../model/collection.ts";
import {useApi} from "../../../service/use-api.ts";
import {useNotificationContext} from "../../../component/notification/notification.tsx";
import ActivityPicker from "../../../component/activity-picker/activity-picker.tsx";
import {DatePicker} from "../../../component/date-picker/date-picker.tsx";
import MoodSelect from "../../../component/mood-select/mood-select.tsx";
import Notes from "../../../component/notes/notes.tsx";
import {useNavigate, useParams} from "react-router-dom";
import WeatherSelect from "../../../component/weather-select/weather-select.tsx";
import {Mood} from "../../../model/mood.ts";
import {Weather} from "../../../model/weather.ts";
import Throbber from "../../../component/throbber/throbber.tsx";
import PhotoUploader from "../../../component/photo-uploader/photo-uploader.tsx";

class Payload {
    mood = Mood.Neutral
    date = new Date().toISOString().split("T")[0]
    weather = Weather.Sunny
    notes = ""
    activities: number[] = []
    photos: string[] = []
}

class Errors {
    mood = ""
    date = ""
    notes = ""
    activities = ""
}

const months = ["Января", "Февраля", "Марта", "Апреля", "Мая", "Июня", "Июля", "Августа", "Сентября", "Октября", "Ноября", "Декабря"]

export default function SaveRecordPage() {
    const willStore = window.location.pathname.includes("/new")
    const [isRecordLoading, setRecordLoading] = useState(!willStore)
    const [areCollectionsLoading, setCollectionsLoading] = useState(true)
    const {state, setState, setField, errors} = useFormState(new Payload(), new Errors())
    const [collections, setCollections] = useState<Collection[]>([])
    const api = useApi()
    const notification = useNotificationContext()
    const navigate = useNavigate()
    const {id} = useParams<string>()


    useEffect(() => {
        document.title = "Новая запись"
        setCollectionsLoading(true)
        if (!willStore) {
            setRecordLoading(true)
            fetchRecord().then()
        }
        fetchCollections().then()
    }, [])


    async function fetchRecord() {
        const {data, status} = await api.get(`/api/records/${id}`)
        if (status != 200) {
            notification.pop(data)
            return
        }

        setState({
            mood: data.mood,
            weather: data.weather,
            date: data.date,
            notes: data.notes,
            activities: data.activities,
            photos: data.photos ?? []
        })
        document.title = dateToTitle(data.date)

        setRecordLoading(false)
    }


    function dateToTitle(date: string): string {
        const day = date.substring(8, 10)
        const month = Number(date.substring(5, 7))
        return `${day} ${months[month]} - Записи`
    }


    async function fetchCollections() {
        const {data, status} = await api.get("/api/collections")
        if (status == 200) {
            setCollections(data)
            setCollectionsLoading(false)
            return
        }

        notification.pop(data)
    }


    function addActivity(id: number) {
        if (state.activities.includes(id)) {
            setState({
                ...state,
                activities: state.activities.filter(activityId => activityId != id)
            })
        } else {
            setState({
                ...state,
                activities: [...state.activities, id]
            })
        }
    }


    async function save() {
        if (willStore) {
            await store()
        } else {
            await update()
        }
    }


    async function store() {
        const {data, status} = await api.post("/api/records", {
            ...state,
            mood: Number(state.mood),
            weather: Number(state.weather),
        })

        if (status == 422) {
            window.scrollTo({top: 0, left: 0, behavior: "smooth"})
            if (data.activities != null) {
                notification.pop("Выберите хотя бы одну активность")
            }
            return
        }

        if (status == 409) {
            notification.pop(data.date)
            window.scrollTo({top: 0, left: 0, behavior: "smooth"})
            return
        }

        if (status == 404) {
            notification.pop(data.activities)
            return
        }

        if (status == 400) {
            notification.pop(data)
            return
        }

        if (status == 403) {
            notification.pop(data)
            return
        }

        navigate("/records")
    }


    async function update() {
        const {data, status} = await api.put(`/api/records/${id}`, {
            mood: Number(state.mood),
            weather: Number(state.weather),
            notes: state.notes,
            activities: state.activities,
            photos: state.photos ?? []
        })

        if (status == 422) {
            if (data.activities != null) {
                notification.pop("Выберите хотя бы одну активность")
            }
            window.scrollTo({top: 0, left: 0, behavior: "smooth"})
            return
        }

        if (status == 404) {
            notification.pop(data)
            return
        }

        navigate("/records")
    }

    function addPhotos(urls: string[]) {
        setState({
            ...state,
            photos: [...state.photos, ...urls]
        })
    }


    async function deletePhoto(name: string) {
        const {status} = await api.delete(`/api/photos/${name}`)
        if (status == 204) {
            setState({
                ...state,
                photos: state.photos.filter(photoName => photoName != name)
            })
        } else {
            notification.pop("Не удалось удалить фотографию")
        }
    }


    return (
        <div className="save-record-page page">
            {isRecordLoading && <Throbber/>}
            {!isRecordLoading &&
                <Form title={willStore ? "Новая запись" : "Запись"} noBackground>
                    <DatePicker active={willStore} name="date" value={state.date} error={errors.date} onChange={setField}/>
                    <MoodSelect name="mood" value={state.mood} onChange={setField}/>
                    <WeatherSelect name="weather" value={state.weather} onChange={setField}/>
                    {areCollectionsLoading && <Throbber/>}
                    {!areCollectionsLoading &&
                        <ActivityPicker collections={collections} value={state.activities} toggleActivity={addActivity}/>
                    }
                    <Notes name="notes" max={5000} value={state.notes} onChange={setField}/>
                    <PhotoUploader name="photos" values={state.photos} onPhotosUploaded={addPhotos} deletePhoto={deletePhoto} />
                    <FormButtons>
                        <FormBackButton/>
                        <FormSubmitButton label="Сохранить" onClick={save}/>
                    </FormButtons>
                </Form>
            }
        </div>
    )
}