import { listDirectory } from "@/lib/data/list"
import { getQueryClient } from "@/lib/queryClient"
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"

export default async function StorageLayout({
	params,
	login,
	content
}: {
	params: Promise<{path?: string[]}>,
	login: React.ReactNode,
	content: React.ReactNode,
}) {
	const path = ((await params).path ?? []).join("/")
	const queryClient = getQueryClient()

	const isAuthorized = await listDirectory(queryClient, path)

	return(
		<HydrationBoundary state={dehydrate(queryClient)}>
      {isAuthorized ? (
				<>
					{content}
				</>
			) : (
				<>
					{login}
				</>
			)}
    </HydrationBoundary>
	)
}