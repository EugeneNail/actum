import "./save-record-page.sass"
import Form from "../../../component/form/form.tsx";
import FormButtons from "../../../component/form/form-button-container.tsx";
import FormBackButton from "../../../component/form/form-back-button.tsx";
import FormSubmitButton from "../../../component/form/form-submit-button.tsx";
import {useFormState} from "../../../service/use-form-state.ts";
import {useEffect, useState} from "react";
import Collection from "../../../model/collection.ts";
import {useHttp} from "../../../service/use-http.ts";
import {useNotificationContext} from "../../../component/notification/notification.tsx";
import ActivityPicker from "../../../component/activity-picker/activity-picker.tsx";
import {DatePicker} from "../../../component/date-picker/date-picker.tsx";
import MoodSelect from "../../../component/mood-select/mood-select.tsx";
import Notes from "../../../component/notes/notes.tsx";
import {useNavigate} from "react-router-dom";

class Payload {
    mood = 0
    date = ""
    notes = ""
    activities: number[] = []
}

class Errors {
    mood = ""
    date = ""
    notes = ""
    activities = ""
}

export default function SaveRecordPage() {
    const {state, setState, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const [collections, setCollections] = useState<Collection[]>([])
    const http = useHttp()
    const notification = useNotificationContext()
    const willStore = window.location.pathname.includes("/new")
    const navigate = useNavigate()


    useEffect(() => {
        fetchCollections()
    }, [])


    async function fetchCollections() {
        const {data, status} = await http.get("/api/collections")
        if (status == 200) {
            setCollections(data)
            return
        }

        notification.pop(data)
    }


    function addActivity(id: number) {
        if (state.activities.includes(id)) {
            setState({
                ...state,
                activities: state.activities.filter(activityId => activityId!= id)
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
        const {data, status} = await http.post("/api/records", {
            ...state,
            mood: Number(state.mood)
        })

        if (status == 409) {
            setErrors(data)
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

    }



    return (
        <div className="save-record-page page">
            <Form title={willStore ? "New record" : "Record"} subtitle={willStore ? "" : state.date}>
                <MoodSelect name="mood" value={state.mood} onChange={setField}/>
                <DatePicker name="date" label="Date" value={state.date} error={errors.date} onChange={setField}/>
                <ActivityPicker collections={collections} value={state.activities} toggleActivity={addActivity}/>
                <Notes label="Notes" name="notes" max={5000} value={state.notes} onChange={setField}/>
                <FormButtons>
                    <FormBackButton/>
                    <FormSubmitButton label="Save" onClick={save}/>
                </FormButtons>
            </Form>
        </div>
    )
}