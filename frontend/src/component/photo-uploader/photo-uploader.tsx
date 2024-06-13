import "./photo-uploader.sass"
import {ChangeEvent} from "react";
import Icon from "../icon/icon.tsx";
import {useApi} from "../../service/use-api.ts";

type Props = {
    name: string
    values: string[]
    deletePhoto: (name: string) => void
    onPhotosUploaded: (photos: string[]) => void
}

export default function PhotoUploader({name, values, deletePhoto, onPhotosUploaded}: Props) {
    const api = useApi()


    async function loadBase64Images(event: ChangeEvent<HTMLInputElement>) {
        const files = event.target.files ?? [] as File[]
        const length = files.length
        const names: string[] = []

        for (let i = 0; i < length; i++) {
            const {data, status} = await api.post("/api/photos", {
                image: await fileToBase64(files[i])
            })

            if (status == 201) {
                names.push(data)
            }
        }

        onPhotosUploaded(names)
    }


    async function fileToBase64(file: File) {
        return new Promise<string>((resolve, reject) => {
            const reader = new FileReader()
            reader.readAsDataURL(file)
            reader.onloadend = () => {
                const value = (reader.result as string).replace("data:", "").replace(/^.+,/, "")
                resolve(value)
            }
            reader.onerror = (error) => reject(error)
        })
    }

    return (
        <div className="photo-uploader">
            <p className="photo-uploader__label">Фотографии</p>
            <input className="photo-uploader__input" name={name} id={name} type="file" multiple accept="image/*" onChange={loadBase64Images}/>
            <div className="photo-uploader__photos">
                {values && values.map((name) => (
                    <div className="photo-uploader__image-container">
                        <div className="photo-uploader__delete-button" onClick={() => deletePhoto(name)}>
                            <Icon className="photo-uploader__delete-button-icon" bold name="close"/>
                        </div>
                        <img className="photo-uploader__img" src={`${import.meta.env.VITE_API_DOMAIN}/api/photos/${name}`}
                             alt={name}/>
                    </div>
                ))}
                <label htmlFor={name} className="photo-uploader__button">
                    <Icon className="photo-uploader__icon" name="add_photo_alternate"/>
                </label>
            </div>

        </div>
    )
}