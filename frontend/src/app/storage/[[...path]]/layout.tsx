import AuthorizationComponent from "@/components/AuthorizationComponent"
import { checkIsAuthorizedOptions } from "@/lib/action/auth/queryOptions"
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

	await queryClient.prefetchQuery(checkIsAuthorizedOptions(path))

	return(
		<HydrationBoundary state={dehydrate(queryClient)}>
      <AuthorizationComponent path={path} login={login} content={content}/>
    </HydrationBoundary>
	)
}