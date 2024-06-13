import axios, {AxiosInstance} from "axios";
import {useNotificationContext} from "../component/notification/notification.tsx";
import base64UrlToString from "./base64.ts";
import {useNavigate} from "react-router-dom";

class Api {
    axios: AxiosInstance = axios.create()
    navigate: (_: any) => void = () => {}


    async post(url: string, data: any) {
        await this.refreshToken()
        return this.axios.post(url, data, {
            headers: {
                Authorization: "Bearer " + localStorage.getItem("Access-Token")
            }
        })
    }


    async put(url: string, data: any) {
        await this.refreshToken()
        return this.axios.put(url, data, {
            headers: {
                Authorization: "Bearer " + localStorage.getItem("Access-Token")
            }
        })
    }


    async delete(url: string) {
        await this.refreshToken()
        return this.axios.delete(url, {
            headers: {
                Authorization: "Bearer " + localStorage.getItem("Access-Token")
            }
        })
    }


    async get(url: string) {
        await this.refreshToken()
        return this.axios.get(url, {
            headers: {
                Authorization: "Bearer " + localStorage.getItem("Access-Token")
            }
        })
    }


    async refreshToken() {
        if (window.location.pathname == "/login" || window.location.pathname == "/signup") {
            return
        }

        const accessToken = localStorage.getItem("Access-Token") ?? ""
        const payload = JSON.parse(base64UrlToString(accessToken.split(".")[1]))
        const now = Math.floor(Date.now() / 1000)

        if (now <= payload.exp) {
            return
        }

        const refreshToken = localStorage.getItem("Refresh-Token") ?? ""
        const {data, status} = await this.axios.post("/api/users/refresh-token", {
            "userId": payload.id,
            "uuid": refreshToken
        })

        if (status == 401 || status == 400 || status == 422) {
            this.navigate("/login")
            return
        }

        if (status == 200) {
            localStorage.setItem("Access-Token", data)
        }
    }
}

export function useApi() {
    const notification = useNotificationContext()
    const navigate = useNavigate()
    const api = new Api()
    api.navigate = navigate

    api.axios = axios.create({
        baseURL: import.meta.env.VITE_API_DOMAIN,
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
        }
    })

    api.axios.interceptors.response.use(
        (response) => response,
        (error) => {
            if (error.response.status == 404 && window.location.pathname != "/not-found") {
                notification.pop(error.response.data)
            }

            if (error.response.status == 500 || error.response.status == 400) {
                notification.pop(error.response.data)
            }

            return error.response
        }
    )

    return api
}
