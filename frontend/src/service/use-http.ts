import axios from "axios";
import {useNotificationContext} from "../component/notification/notification.tsx";

export function useHttp() {
    const notification = useNotificationContext()

    const http = axios.create({
        baseURL: "http://192.168.1.3:8080",
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
            Authorization: "Bearer " + localStorage.getItem("Access-Token")
        }
    })

    http.interceptors.response.use(
        (response) => response,
        (error) => {
            if (error.response.status == 401 && window.location.pathname != "/login") {
                window.location.href = "/login"
            }

            if (error.response.status == 404 && window.location.pathname != "/not-found") {
                notification.pop(error.response.data)
            }

            if (error.response.status == 500) {
                notification.pop(error.response.data)
            }

            return error.response
        }
    )

    return http
}