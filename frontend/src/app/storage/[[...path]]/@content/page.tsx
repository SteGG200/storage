import NewFolderButton from "@/components/button/NewFolderButton"
import UploadFileButton from "@/components/button/UploadFileButton"
import Nav from "@/components/Nav"
import SearchBar from "@/components/SearchBar"
import TableContent from "@/components/TableContent"

export default async function ContentPage({
	params
}: {
	params: Promise<{path?: string[]}>
}) {
	const path = ((await params).path ?? []).join('/')

	return (
		<div className="min-h-dvh container mx-auto p-4 space-y-8">
			<h1 className="text-2xl font-bold">My Storage</h1>
			<Nav path={path}/>
			<div className="flex justify-between items-center mb-4">
				<SearchBar path={path}/>
				<div className="space-x-2">
					<UploadFileButton path={path}/>
					<NewFolderButton path={path} />
				</div>
			</div>
			<TableContent path={path}/>
		</div>
	)
}