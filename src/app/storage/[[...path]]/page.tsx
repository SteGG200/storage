import PathShowAndSearch from "@/components/TopBar/PathShowAndSearch"
import TopBar from "@/components/TopBar/TopBar"

interface IParams {
	params: {
		path: string[]
	}
}

export default function Storage({ params }: IParams){
	return (
		<main className="h-dvh bg-customBlack p-4 max-md:p-2">
			<TopBar>
				<PathShowAndSearch directory={params.path ?? []}/>
			</TopBar>

			Here is your storage on directory: Home{params.path?.map((folder) => {
				return `/${folder}`
			})}
		</main>
	)
}