import {Outlet, useLocation} from "react-router-dom";
import "./collections-page.sass"
import {useEffect, useState} from "react";
import Collection from "../../../model/collection.ts";
import {useHttp} from "../../../service/use-http.ts";
import CollectionCard from "../../../component/collection/collection-card.tsx";
import NewCollection from "../../../component/new-collection/new-collection.tsx";

export default function CollectionsPage() {
    const [collections, setCollections] = useState<Collection[]>([])
    const http = useHttp()
    const location = useLocation()

    useEffect(() => {
        http.get("/collections").then(({data}) => {
            setCollections(data)
        })
    }, [location.pathname])

    return (
        <div className="collections-page">
            <Outlet/>
            {collections?.length > 0 && collections.map(collection =>
                <CollectionCard collection={collection}/>
            )}
            <NewCollection/>
        </div>
    )
}