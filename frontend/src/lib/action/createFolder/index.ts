'use server'

import { getToken } from "@/lib/utils/server"

export const createFolder = async (
	path: string,
	formData: FormData
) => {
	const token = await getToken()
	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/create/${path}`, {
		method: 'POST',
		headers: {
			"Authorization": `Bearer ${token}`
		},
    body: formData
	})

	if (!response.ok) {
    const msg = await response.text()
    throw new Error(msg)
  }
}