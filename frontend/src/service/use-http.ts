import axios from "axios";

export function useHttp() {
    const http = axios.create({
        baseURL: "http://localhost:8080",
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json"
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