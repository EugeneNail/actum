import "./form-header.sass"
import Icon from "../icon/icon.tsx";

type Props = {
    icon: string
    title: string
}

export default function FormHeader({icon, title}: Props) {
    return (
        <div className="form-header">
            <div className="form-header__background">
                <div className="form-header__icon-container">
                    <Icon name={icon} className="form-header__icon"/>
                </div>
            </div>
            <h3 className="form-header__title">{title}</h3>
        </div>
    )
}