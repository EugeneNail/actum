import "./header.sass"
import HeaderLink from "./header-link.tsx";
import Button from "../button/button.tsx";
import {useNavigate} from "react-router-dom";
import {useEffect, useState} from "react";

export default function Header() {
    const [isPortrait, setPortrait] = useState(false)
    const navigate = useNavigate()

    useEffect(() => {
        const handleResize = () => {
            const {innerWidth, innerHeight} = window;
            setPortrait(innerWidth < innerHeight);
        };

        handleResize();
        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('resize', handleResize);
        };
    }, [])

    return (
        <header className="header">
            <div className="header__logo-container">
                <img src="/img/logo.png" alt="" className="header__logo-img"/>
                <p className="header__logo-text">Actum</p>
            </div>
            <HeaderLink to="/records" icon="post" label="Records"/>
            <HeaderLink to="/statistics" icon="bar_chart_4_bars" label="Statistics"/>
            <div className="header__placeholder"/>
            <Button className="header__new-record-button" icon="add" label={isPortrait ? "" : "New record"} onClick={() => navigate("/records/new")}/>
            <HeaderLink to="/calendar" icon="calendar_month" label="Calendar"/>
            <HeaderLink to="/settings" icon="settings" label="Settings"/>
        </header>
    )
}