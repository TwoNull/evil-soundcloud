import { FieldValues, useForm } from "react-hook-form";
import { FetchPlaylist } from "../wailsjs/go/main/App"

function App() {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm();
    
    async function fetchPlaylist(vals: FieldValues) {
        await FetchPlaylist(vals.url)
    }

    return (
        <div className="min-h-screen bg-white grid grid-cols-1 place-items-center justify-items-center mx-auto py-8">
            <div className="text-blue-900 text-2xl font-bold font-mono">
                <h1 className="content-center">Evil Soundcloud</h1>
            </div>
            <div className="w-fit max-w-md">
                <form className="flex flex-col justify-center items-between" onSubmit={handleSubmit(fetchPlaylist)}>
                    <label>Soundcloud Playlist URL:</label>
                    <input className="bg-neutral-100" type="text" {...register("url")} />
                    <button type="submit">Submit</button>
                </form>
            </div>
        </div>
    )
}

export default App
