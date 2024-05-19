import {ChangeEvent, useState} from "react";

export function useFormState<D, E>(initial: D, initialErrors: E) {
    const [state, setState] = useState(initial)
    const [errors, _setErrors] = useState(initialErrors)

    function setField(event: ChangeEvent<HTMLInputElement>) {
        event.preventDefault()
        setState({
            ...state,
            [event.target.name] : event.target.value
        })
    }

    function setErrors(errors: any) {
        _setErrors(errors)
    }

    return {state, setField, setState, errors, setErrors}
}