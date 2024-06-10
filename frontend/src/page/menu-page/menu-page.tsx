import "./menu-page.sass"
import {useNavigate} from "react-router-dom";
import MenuOption from "./menu-option.tsx";
import {Color} from "../../model/color.tsx";
import {useApi} from "../../service/use-api.ts";
import {useNotificationContext} from "../../component/notification/notification.tsx";

export default function MenuPage() {
    const navigate = useNavigate()
    const api = useApi()
    const notification = useNotificationContext()

    async function logout() {
        const {data, status} = await api.post("/api/users/logout", null)

        if (status >= 400) {
            notification.pop(data)
            return
        }

        localStorage.removeItem("Refresh-Token")
        localStorage.removeItem("Access-Token")
        navigate("/login")
    }

    return (
        <div className="menu-page page">
            <div className="menu-page__menu">
                <div className="menu-page__group">
                    <MenuOption icon="bar_chart" label="Статистика" color={Color.Green} onClick={() => navigate("/statistics")}/>
                    <MenuOption icon="post" label="Записи" color={Color.Red} onClick={() => navigate("/records")}/>
                    <MenuOption icon="category" label="Коллекции" color={Color.Orange} onClick={() => navigate("/collections")}/>
                </div>
                <div className="menu-page__group">
                    <MenuOption icon="logout" label="Выйти" color={Color.Accent} onClick={logout}/>
                </div>

            </div>
        </div>
    )
}