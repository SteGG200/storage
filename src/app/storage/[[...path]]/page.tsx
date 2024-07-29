import MainContent from "@/components/folder_content/MainContent"
import PathShowAndSearch from "@/components/top_bar/PathShowAndSearch"
import TopBar from "@/components/top_bar/TopBar"

interface IParams {
	params: {
		path: string[]
	}
}

export default function Storage({ params }: IParams){
	const items : IItem[] = [
		{
			type: "folder",
      name: "flkajsdlfjasdlfjalsdjflakfkalsjdfladjfaldfjalsdkfjalsdjfladkjfalsdfajsdlfjasdlkfjaskdflaskdjffasdfka",
      size: 0,
			key: "0"
		},
		{
			type: "file",
      name: "File 1.txt",
      size: 1024,
			key: "1"
		}
	]

	return (
		<main className="h-dvh bg-customBlack p-4 max-md:p-2 space-y-4">
			<TopBar>
				<PathShowAndSearch directory={params.path ?? []}/>
			</TopBar>

			<MainContent items={items}/>
		</main>
	)
}