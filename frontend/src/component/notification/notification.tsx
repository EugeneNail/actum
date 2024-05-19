import {createContext, ReactNode, useContext, useState} from "react";
import "./notification.sass"
import classNames from "classnames";
import Icon from "../icon/icon.tsx";

type Props = {
    children: ReactNode
}

type Notification = {
    level: string
    message: string
    isError: boolean
}

const StateContext = createContext({
    notification: {} as Notification|null|undefined,
    pop: (_: string) => {},
    popError: (_: string) => {},
})

export default function Notification({children} : Props) {
    const [notification, setNotification] = useState<Notification|null>()

    function pop(message: string) {
        setNotification({
            level: "Warning",
            message,
            isError: false
        })

        setTimeout(() => {
            setNotification(null)
        }, 5000)
    }

    function popError() {
        setNotification({
            level: "Error",
            message: "The server encountered and error and could not complete your request.",
            isError: true
        })

        setTimeout(() => {
            setNotification(null)
        }, 5000)
    }

    return (
        <StateContext.Provider value={{notification, pop, popError}}>
            <>
                {notification &&
                    <div className={classNames("notification", {"error": notification.isError})}>
                        <h6 className="notification__level">{notification.level}</h6>
                        <div className="notification__icon-container">
                            <Icon name="warning" className="notification__icon"></Icon>
                        </div>
                        <p className="notification__message">{notification.message}</p>
                        <div className="notification__timer"/>
                    </div>
                }
                {children}
            </>
        </StateContext.Provider>
    )
}

export const useNotificationContext = () => useContext(StateContext)