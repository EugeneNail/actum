import "./form.sass"
import Button, {ButtonStyle} from "../button/button.tsx";
import Icon from "../icon/icon.tsx";
import {Color} from "../../model/color.tsx";

type Props = {
    onClick: () => void
}

export default function FormDeleteButton({onClick}: Props) {
    return (
        <Button className="form__delete-button" even style={ButtonStyle.Secondary} color={Color.Red} onClick={onClick}>
            <Icon className="form__delete-button-icon" name="delete" bold/>
        </Button>
    )
}