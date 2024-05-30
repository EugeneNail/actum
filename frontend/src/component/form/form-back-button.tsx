import "./form.sass"
import Button, {ButtonStyle} from "../button/button.tsx";
import Icon from "../icon/icon.tsx";

export default function FormBackButton() {
    return (
        <Button even style={ButtonStyle.Secondary} onClick={() => window.history.back()} >
            <Icon name="west" bold/>
        </Button>
    )
}