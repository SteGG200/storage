'use client'

import { checkIsAuthorizedOptions } from "@/lib/action/auth/queryOptions"
import { useQuery } from "@tanstack/react-query"
import { useRouter } from "next/navigation"

interface AuthorizationComponentProps {
	path: string
	login: React.ReactNode
	content: React.ReactNode
}

export default function AuthorizationComponent({
	path,
	login,
	content
}: AuthorizationComponentProps) {
	const { data, error, status } = useQuery(checkIsAuthorizedOptions(path))
	const router = useRouter()

	if(status === 'error'){
		throw error
	}

	if(status === 'success'){
		if(data.isAuthorized) {
			return content
		}

		if(data.unauthorizedPath != `/${path}`){
			router.push(`/storage${data.unauthorizedPath}`)
		}

		return login
	}
}