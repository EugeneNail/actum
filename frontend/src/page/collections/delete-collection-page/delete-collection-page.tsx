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
                <Form title={`Удалить коллекцию "${name}"?`}>
                    <p className="justified">Удаление коллекции удалит все ее активности.</p>
                    <br/>
                    <p className="justified">Активности также будут удалены из всех ваших записей. Это действие необратимо. Вы действительно хотите удалить коллекцию?</p>
                    <FormButtons>
                        <FormBackButton/>
                        <FormSubmitButton label="Удалить" color={Color.Red} onClick={destroy}/>
                    </FormButtons>
                </Form>
            }
        </div>
    )
}