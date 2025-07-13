'use server'
import fs from 'fs/promises'
import crypto from 'crypto'
import path from 'path'
import { getToken } from '@/lib/utils/server'

export const removeOldTmpFile = async () => {
	for(const file of await fs.readdir(process.env.TMP_DIR!)){
		await fs.unlink(path.join(process.env.TMP_DIR!, file))
	}
}

export const getSizeOfFile = async (
	path: string,
	filename: string
) => {
	// Get total size of file
	const token = await getToken()
	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/download/${path}/${filename}`, {
		method: 'HEAD',
		headers: {
			"Authorization": `Bearer ${token}`
		}
	})

	return parseInt(response.headers.get('content-length') ??  "0")
}

export const generateRandomHash = async () => {
	return crypto.randomBytes(16).toString('hex')
}

export const downloadChunk = async (
	path: string, 
	filename: string,
	hashString: string,
	start: number,
	end: number
) => {
	const token = await getToken()
	const chunkResponse = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/download/${path}/${filename}`, {
		headers: {
			"Authorization": `Bearer ${token}`,
			'Range': `bytes=${start}-${end - 1}`
		}
	})

	const data = await chunkResponse.arrayBuffer()

	const buffer = Buffer.from(data)
	await fs.appendFile(`${process.env.TMP_DIR}/${hashString}_${filename}`, buffer)

}