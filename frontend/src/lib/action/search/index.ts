'use server'

import { getToken } from "@/lib/utils/server"

export const searchAction = async (path: string, searchValue: string) => {
	const token = await getToken()

	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/get/search/${path}?q=${searchValue}`, {
		headers: {
			"Authorization": `Bearer ${token}`
		}
	})

	const data: Item[] = await response.json()

	return data
}