import "./notes.sass"
import {useRef} from "react";
import classNames from "classnames";

type Props = {
    className?: string
    name: string
    max: number
    value: string
    onChange: (event: any) => void
}

export default function Notes({className, name, max, value, onChange}: Props) {
    const ref = useRef<HTMLTextAreaElement>(document.createElement('textarea'))

    function resizeToContent() {
        ref.current.style.height = ref.current.scrollHeight + "px"
    }

    return (
        <div className={classNames("notes", className)}>
            <label className="notes__label" htmlFor={name}>Как прошел день?</label>
            <textarea className="notes__textarea"
                      ref={ref}
                      placeholder="Расскажите, что интересного случилось"
                      value={value}
                      name={name}
                      id={name}
                      onChange={onChange}
                      autoComplete="off"
                      autoCorrect="on"
                      onInput={resizeToContent}
                      maxLength={max}/>
            <p className="notes__limit">{value.length} / {max}</p>
        </div>
    )
}