import Icon from "../icon/icon";
import "./new-collection.sass"
import {useNavigate} from "react-router-dom";

export default function NewCollection() {
    const navigate = useNavigate()

    return (
        <div className="new-collection" onClick={() => navigate("./new")}>
            <div className="new-collection__icon-container">
                <Icon className="new-collection__icon" name="add" />
            </div>
        </div>
    )
}