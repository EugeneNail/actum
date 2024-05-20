import {Outlet, useLocation} from "react-router-dom";
import "./collections-page.sass"
import {useEffect, useState} from "react";
import Collection from "../../../model/collection.ts";
import {useHttp} from "../../../service/use-http.ts";
import CollectionCard from "../../../component/collection-card/collection-card.tsx";
import NewCollection from "../../../component/new-collection/new-collection.tsx";

export default function CollectionsPage() {
    const [collections, setCollections] = useState<Collection[]>([])
    const http = useHttp()
    const location = useLocation()

    useEffect(() => {
        http.get("/collections").then(({data}) => {
            const sorted = data.sort((a: Collection, b: Collection) => a.id - b.id)
            setCollections(sorted)
        })
    }, [location.pathname])

    return (
        <div className="collections-page">
            <Outlet/>
            {collections?.length > 0 && collections.map(collection =>
                <CollectionCard collection={collection}/>
            )}
            {collections.length < 15 && <NewCollection/>}
        </div>
    )
}