import "./delete-activity-page.sass"
import {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {useApi} from "../../../service/use-api.ts";
import Form from "../../../component/form/form";
import FormButtons from "../../../component/form/form-button-container";
import FormBackButton from "../../../component/form/form-back-button";
import FormSubmitButton from "../../../component/form/form-submit-button.tsx";
import {Color} from "../../../model/color.tsx";
import Throbber from "../../../component/throbber/throbber.tsx";

export default function DeleteActivityPage() {
    const [isLoading, setLoading] = useState(true)
    const {activityId} = useParams<string>()
    const api = useApi()
    const navigate = useNavigate()
    const [name, setName] = useState("")

    useEffect(() => {
        fetchActivity().then()
    }, [])

    async function fetchActivity() {
        const {data, status} = await api.get(`/api/activities/${activityId}`)

        if (status == 403) {
            navigate("/collections")
            return
        }

        setName(data.name)
        document.title = data.name + " - Активности"
        setLoading(false)
    }

    async function destroy() {
        const {status} = await api.delete(`/api/activities/${activityId}`)
        if (status == 204) {
            navigate("/collections")
        }
    }

    return (
        <div className="delete-activity-page page">
            {isLoading && <Throbber/>}
            {!isLoading &&
                <Form title={`Удалить активность "${name}"?`}>
                    <p className="justified">Удаление активности также удалит ее из всех ваших записей.</p>
                    <br/>
                    <p className="justified">Это действие необратимо. Вы действительно хотите удалить активность?</p>
                    <FormButtons>
                        <FormBackButton/>
                        <FormSubmitButton label="Удалить" color={Color.Red} onClick={destroy}/>
                    </FormButtons>
                </Form>
            }
        </div>
    )
}