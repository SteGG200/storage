import MainContent from "@/components/folder_content/MainContent"
import TopBar from "@/components/top_bar/TopBar"

export default function Storage({ params }: IParams){
	return (
		<main className="h-dvh bg-customBlack p-4 max-md:p-2 space-y-4">
			<TopBar directory={params.path ?? []}/>

			<MainContent path={params.path?.join('/') ?? ""}/>
		</main>
	)
}