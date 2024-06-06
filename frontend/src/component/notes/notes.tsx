import "./notes.sass"
import {ChangeEvent, useRef} from "react";
import classNames from "classnames";

type Props = {
    className?: string
    label: string
    name: string
    max: number
    value: string
    onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

export default function Notes({className, label, name, max, value, onChange}: Props) {
    const ref = useRef<HTMLTextAreaElement>(document.createElement('textarea'))

    function resizeToContent() {
        ref.current.style.height = ref.current.scrollHeight + "px"
    }

    return (
        <div className={classNames("notes", className)}>
            <label  className="notes__label" htmlFor={name}>{label}</label>
            <textarea className="notes__textarea"
                      ref={ref}
                      placeholder="What interesting things happened?"
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