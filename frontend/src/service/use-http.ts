import axios from "axios";

export function useHttp() {
    const http = axios.create({
        baseURL: "http://192.168.1.3:8080/api",
        withCredentials: true,
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
                window.location.href = "/not-found"
            }

            return error.response
        }
    )

    return http
}