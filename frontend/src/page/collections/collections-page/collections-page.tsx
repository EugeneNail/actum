import "./collections-page.sass"
import {useEffect, useState} from "react";
import Collection from "../../../model/collection";
import {useHttp} from "../../../service/use-http.ts";
import CollectionCard from "../../../component/collection-card/collection-card.tsx";
import Icon from "../../../component/icon/icon.tsx";
import {useNavigate} from "react-router-dom";

export default function CollectionsPage() {
    const [collections, setCollections] = useState<Collection[]>([])
    const http = useHttp()
    const navigate = useNavigate()

    useEffect(() => {
        http.get("/api/collections").then(({data}) => {
            setCollections(data)
        })
    }, [])

    return (
        <div className="collections-page page">
            {collections && collections.map((collection) =>
                <CollectionCard key={collection.id} collection={collection}/>
            )}
            <div className="collections-page-button" onClick={() => navigate("./new")}>
                <Icon name="add" className="collections-page-button__icon" bold/>
                <p className="collections-page-button__label">Add collection</p>
            </div>
        </div>
    )
}