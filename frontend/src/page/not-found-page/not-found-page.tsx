import "./not-found-page.sass"
import Button from "../../component/button/button.tsx";

export default function NotFoundPage() {
    function goHome() {
        window.location.href = "/"
    }

    return (
        <div className="not-found-page">
            <div className="not-found-page__text-container">
                <h1 className="not-found-page__title">404</h1>
                <h2 className="not-found-page__subtitle">page not found</h2>
                <p className="not-found-page__message">The page you are looking for was moved, removed, renamed, or might never exist in this world!</p>
                <Button className="not-found-page__button" icon="west" label="Back to Homepage" onClick={goHome}/>
            </div>
            <img src="/img/404.png" alt="" className="not-found-page__image"/>
        </div>
    )
}