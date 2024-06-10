import "./collections-page.sass"
import {useEffect, useState} from "react";
import Collection from "../../../model/collection";
import {useApi} from "../../../service/use-api.ts";
import CollectionCard from "../../../component/collection-card/collection-card.tsx";
import Icon from "../../../component/icon/icon.tsx";
import {useNavigate} from "react-router-dom";
import Throbber from "../../../component/throbber/throbber.tsx";

export default function CollectionsPage() {
    const [isLoading, setLoading] = useState(true)
    const [collections, setCollections] = useState<Collection[]>([])
    const api = useApi()
    const navigate = useNavigate()

    useEffect(() => {
        document.title = "Коллекции"
        setLoading(true)
        api.get("/api/collections").then(({data}) => {
            setCollections(data)
            setLoading(false)
        })
    }, [])

    return (
        <>
            {isLoading &&
                <div className="page">
                    <Throbber/>
                </div>
            }
            {!isLoading &&
                <div className="collections-page page">
                    {collections && collections.map((collection) =>
                        <CollectionCard key={collection.id} collection={collection}/>
                    )}
                    <div className="collections-page-button" onClick={() => navigate("./new")}>
                        <div className="collections-page-button__title-container">
                            <Icon name="add" className="collections-page-button__icon" bold/>
                            <p className="collections-page-button__label">Добавить коллекцию</p>
                        </div>
                    </div>
                </div>
            }
        </>
    )
}