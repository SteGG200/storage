import { checkIsAuthorized } from "@/lib/action/auth"
import { redirect } from "next/navigation"

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

	const response = await checkIsAuthorized(path)

	if(response.isAuthorized){
		return content
	}

	if(response.unauthorizedPath.slice(1) != path){
		redirect(`/storage/${response.unauthorizedPath.slice(1)}`)
	}

	return login
}