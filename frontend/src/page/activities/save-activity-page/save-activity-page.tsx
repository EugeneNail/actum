import "./save-activity-page.sass"
import Form from "../../../component/form/form.tsx";
import {useEffect, useState} from "react";
import {useHttp} from "../../../service/use-http.ts";
import {useFormState} from "../../../service/use-form-state.ts";
import {Outlet, useNavigate, useParams} from "react-router-dom";
import {useNotificationContext} from "../../../component/notification/notification.tsx";
import Field from "../../../component/field/field.tsx";
import FormButtons from "../../../component/form/form-button-container.tsx";
import FormBackButton from "../../../component/form/form-back-button.tsx";
import FormSubmitButton from "../../../component/form/form-submit-button.tsx";
import FormDeleteButton from "../../../component/form/form-delete-button.tsx";
import IconSelect from "../../../component/icon-select/icon-select.tsx";
import Throbber from "../../../component/throbber/throbber.tsx";

class Payload {
    name: string = ""
    icon: number = 0
    collectionId: number = 0
}

class Errors {
    name: string = ""
    icon: string = ""
    collectionId: string = ""
}

export default function SaveActivityPage() {
    const [isCollectionLoading, setCollectionLoading] = useState(true)
    const [isActivityLoading, setActivityLoading] = useState(false)
    const http = useHttp()
    const {state, setState, setField, errors, setErrors} = useFormState(new Payload(), new Errors())
    const willStore = window.location.pathname.includes("/new")
    const navigate = useNavigate()
    const notification = useNotificationContext()
    const {collectionId, activityId} = useParams<string>()
    const [collectionName, setCollectionName] = useState<string>()

    useEffect(() => {
        document.title = "Новая активность"
        setState({
            ...state,
            collectionId: parseInt(collectionId ?? "0"),
            icon: 100
        })
        fetchCollection().then()

        if (!willStore) {
            setActivityLoading(true)
            fetchActivity().then()
        }
    }, [])


    async function fetchCollection() {
        const {data, status} = await http.get(`/api/collections/${collectionId}`)
        if (status == 403 || status == 404) {
            notification.pop(data)
            navigate("/collections")
        }
        setCollectionName(data.name)
        setCollectionLoading(false)
    }


    async function fetchActivity() {
        const {data, status} = await http.get(`/api/activities/${activityId}`)
        if (status == 403) {
            notification.pop(data)
            navigate("/collections")
            return
        }

        setState({
            ...state,
            name: data.name,
            icon: data.icon
        })
        document.title = data.name + " - Активности"
        setActivityLoading(false)
    }


    async function save() {
        if (willStore) {
            await store()
        } else {
            await update()
        }
    }


    async function store() {
        const {data, status} = await http.post("/api/activities", {
            name: state.name,
            icon: Number(state.icon),
            collectionId: parseInt(collectionId ?? "0")
        })

        if (status == 422 || status == 409) {
            window.scrollTo({top: 0, left: 0, behavior: "smooth"})
            setErrors(data)
            return
        }

        if (status == 400) {
            return
        }

        navigate("/collections")
    }


    async function update() {
        const {data, status} = await http.put(`/api/activities/${activityId}`, {
            name: state.name,
            icon: Number(state.icon)
        })

        if (status == 403) {
            notification.pop(data)
            return
        }

        if (status == 422) {
            window.scrollTo({top: 0, left: 0, behavior: "smooth"})
            setErrors(data)
            return
        }

        if (status == 400) {
            return
        }

        navigate("/collections")
    }


    return (
        <>
            {(isActivityLoading || isCollectionLoading) &&
                <div className="page">
                    <Throbber/>
                </div>
            }
            {!isActivityLoading && !isCollectionLoading &&
                <div className="save-activity-page page">
                    <Form title={willStore ? "Новая активность" : state.name} subtitle={(willStore ? "" : "Активность") + ` коллекции "${collectionName}"`}>
                        <Field name="name" label="Название" icon="webhook" value={state.name} error={errors.name} onChange={setField}/>
                        <IconSelect className="save-activity-page__icon-select" name="icon" value={state.icon} onChange={setField}/>
                        <FormButtons>
                            <FormBackButton/>
                            <FormSubmitButton label="Сохранить" onClick={save}/>
                            {!willStore && <FormDeleteButton onClick={() => navigate("./delete")}/>}
                        </FormButtons>
                    </Form>
                    <Outlet/>
                </div>
            }
        </>
    )
}