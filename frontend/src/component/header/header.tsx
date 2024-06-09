import "./header.sass"
import {NavLink, useNavigate} from "react-router-dom";
import {useEffect, useState} from "react";
import Icon from "../icon/icon.tsx";
import Button from "../button/button.tsx";
import {Color} from "../../model/color.tsx";

export default function Header() {
    const navigate = useNavigate()
    const [isVisible, setVisible] = useState(true)

    useEffect(() => {
        let timer: number

        const handleScroll = () => {
            setVisible(true)

            if (timer) {
                clearTimeout(timer)
            }
            timer = setTimeout(() => {
                setVisible(false)
            }, 3000)
        }

        window.addEventListener('wheel', handleScroll)
        window.addEventListener('touchmove', handleScroll)

        return () => {
            window.removeEventListener('wheel', handleScroll)
            window.removeEventListener('touchmove', handleScroll)
            if (timer) {
                clearTimeout(timer)
            }
        }
    }, [])

    return (
        <>
            {isVisible && <header className="header">
                <NavLink className="header-link" to="/statistics">
                    <Icon className="header-link__icon" name="bar_chart"/>
                </NavLink>
                <NavLink className="header-link" to="/records">
                    <Icon className="header-link__icon" name="post"/>
                </NavLink>
                <div className="header__placeholder">
                    <div className="header__button-container">
                        <Button className="header__button" color={Color.Accent} even round
                                onClick={() => navigate("/records/new")}>
                            <Icon className="header__button-icon" name="add" bold/>
                        </Button>
                    </div>
                </div>
                <NavLink className="header-link" to="/collections">
                    <Icon className="header-link__icon" name="category"/>
                </NavLink>
                <NavLink className="header-link" to="/setting">
                    <Icon className="header-link__icon" name="settings"/>
                </NavLink>
            </header>
            }

        </>
    )
}