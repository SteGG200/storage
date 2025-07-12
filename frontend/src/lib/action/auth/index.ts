'use server'

import { cookies } from "next/headers"

export const checkIsAuthorized = async (
	path: string
) : Promise<{
	isAuthorized: true
} | {
	isAuthorized: false,
	unauthorizedPath: string
}> => {
	const cookieStore = await cookies()
	const token = cookieStore.get('token')

	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/auth/check/${path}`, {
		headers: {
			"Authorization": `Bearer ${token?.value ?? ""}`
		}
	})

	if(!response.ok) {
		if(response.status === 401) {
			const { unauthorizedPath } = await response.json()
			return {
				isAuthorized: false,
				unauthorizedPath
			}
		}else {
			const message = await response.text()
			throw new Error(message)
		}
	}

	return {
		isAuthorized: true
	}
}

export const login = async (
	path: string,
	formData: FormData
) => {
	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/auth/login/${path}`,{
		method: 'POST',
		body: formData
	})

	if(!response.ok){
		const message = await response.text()
		throw new Error(message)
	}

	const data : { token : string } = await response.json()
	const cookieStore = await cookies()
	cookieStore.set('token', data.token)
}