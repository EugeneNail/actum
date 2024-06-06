import "./form.sass"
import Button, {ButtonStyle} from "../button/button.tsx";
import {Color} from "../../model/color.tsx";

type Props = {
    label: string
    color?: Color
    onClick: () => void
}

export default function FormSubmitButton({label, color = Color.Accent, onClick}: Props) {
    return (
        <Button className="form__submit-button" submit style={ButtonStyle.Primary} color={color} onClick={onClick}>
            {label}
        </Button>
    )
}