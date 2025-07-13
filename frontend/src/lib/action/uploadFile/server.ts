'use server'

import { getToken } from "@/lib/utils/server"

export const getUploadSessionToken = async (
	path: string,
	formData: FormData
) => {
	const authToken = await getToken()
	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/upload/token/${path}`, {
		method: 'POST',
		headers: {
			"Authorization": `Bearer ${authToken}`,
		},
		body: formData
	})

	if (!response.ok) {
		const msg = await response.text()
    throw new Error(msg)
  }

	const { token } : { token: string } = await response.json()
	return token
}

export const uploadByChunk = async (
	formData: FormData,
	uploadSessionToken: string
) => {
	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/upload/file`, {
		method: 'POST',
		headers: {
			'Authorization': `Bearer ${uploadSessionToken}`
		},
		body: formData
	})

	if (!response.ok) {
		throw new Error('Failed to upload file')
	}
}