import "./delete-activity-page.sass"
import {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {useHttp} from "../../../service/use-http";
import Form from "../../../component/form/form";
import FormButtons from "../../../component/form/form-button-container";
import FormBackButton from "../../../component/form/form-back-button";
import FormSubmitButton from "../../../component/form/form-submit-button.tsx";
import {Color} from "../../../model/color.tsx";

export default function DeleteActivityPage() {
    const {activityId} = useParams<string>()
    const http = useHttp()
    const navigate = useNavigate()
    const [name, setName] = useState("")

    useEffect(() => {
        fetchActivity()
    }, [])

    async function fetchActivity() {
        const {data, status} = await http.get(`/api/activities/${activityId}`)

        if (status == 403) {
            navigate("/collections")
            return
        }

        console.log(data)

        setName(data.name)
    }

    async function destroy() {
        const {status} = await http.delete(`/api/activities/${activityId}`)
        if (status == 204) {
            navigate("/collections")
        }
    }

    return (
        <div className="delete-activity-page page">
            <Form title={`Delete activity "${name}"?`}>
                <p className="justified">Deleting activity will remove it from all records.</p>
                <br/>
                <p className="justified">You can also edit activity. Do you want to delete the activity?</p>
                <FormButtons>
                    <FormBackButton/>
                    <FormSubmitButton label="Delete" color={Color.red} onClick={destroy}/>
                </FormButtons>
            </Form>
        </div>
    )
}