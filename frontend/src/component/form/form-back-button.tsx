import "./form.sass"
import Button, {ButtonStyle} from "../button/button.tsx";
import Icon from "../icon/icon.tsx";
import {Color} from "../../model/color.tsx";

type Props = {
    color?: Color
}

export default function FormBackButton({color = Color.Accent}: Props) {
    return (
        <Button even color={color} style={ButtonStyle.Secondary} onClick={() => window.history.back()} >
            <Icon name="west" bold/>
        </Button>
    )
}