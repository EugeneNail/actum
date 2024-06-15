import "./header.sass"
import {useNavigate} from "react-router-dom";
import Icon from "../icon/icon.tsx";
import Button from "../button/button.tsx";
import {Color} from "../../model/color.tsx";
import HeaderLink from "./header-link.tsx";

export default function Header() {
    const navigate = useNavigate()

    return (
        <header className="header">
            <HeaderLink label="Отчеты" color={Color.Green} icon="bar_chart" to="/statistics"/>
            <HeaderLink label="Записи" color={Color.Red} icon="post" to="/records"/>
            <div className="header__placeholder">
                <div className="header__button-container">
                    <Button className="header__button" color={Color.Accent} even round onClick={() => navigate("/records/new")}>
                        <Icon className="header__button-icon" name="add" bold/>
                    </Button>
                </div>
            </div>
            <HeaderLink label="Коллекции" color={Color.Orange} icon="category" to="/collections"/>
            <HeaderLink label="Меню" color={Color.Accent} icon="more_horiz" to="/menu"/>
        </header>
    )
}