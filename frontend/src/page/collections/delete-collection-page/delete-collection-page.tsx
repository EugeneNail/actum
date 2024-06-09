import "./delete-collection-page.sass"
import {useEffect, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import {useHttp} from "../../../service/use-http";
import Form from "../../../component/form/form";
import FormButtons from "../../../component/form/form-button-container";
import FormBackButton from "../../../component/form/form-back-button";
import FormSubmitButton from "../../../component/form/form-submit-button.tsx";
import {Color} from "../../../model/color.tsx";
import Throbber from "../../../component/throbber/throbber.tsx";

export default function DeleteCollectionPage() {
    const [isLoading, setLoading] = useState(true)
    const {id} = useParams<string>()
    const http = useHttp()
    const navigate = useNavigate()
    const [name, setName] = useState("")

    useEffect(() => {
        fetchCollection().then()
    }, [])

    async function fetchCollection() {
        const {data, status} = await http.get(`/api/collections/${id}`)

        if (status == 403) {
            navigate("/collections")
            return
        }

        setName(data.name)
        document.title = data.name + " - Коллекции"
        setLoading(false)
    }

    async function destroy() {
        const {status} = await http.delete(`/api/collections/${id}`)
        if (status == 204) {
            navigate("/collections")
        }
    }

    return (
        <div className="delete-collection-page page">
            {isLoading && <Throbber/>}
            {!isLoading &&
                <Form title={`Delete collection "${name}"?`}>
                    <p className="justified">Deleting collection will remove all activities within.</p>
                    <br/>
                    <p className="justified">Activities will also be removed from your records. You can also edit collection. Do you want to delete the collection?</p>
                    <FormButtons>
                        <FormBackButton/>
                        <FormSubmitButton label="Delete" color={Color.Red} onClick={destroy}/>
                    </FormButtons>
                </Form>
            }
        </div>
    )
}