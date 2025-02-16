import Nav from "@/components/Nav"
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
			<TableContent path={path}/>
		</div>
	)
}