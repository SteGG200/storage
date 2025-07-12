import { sleep } from "@/lib/utils/client"
import { downloadChunk, generateToken, getSizeOfFile, removeOldTmpFile } from "./server"

const CHUNK_SIZE = 1024 * 1024 * 5
const TIMEOUT = 1000 * 0.5

export const downloadFile = async (
	path: string,
	filename: string,
	progressHandler: (currentProgress: number) => void
) => {
	await removeOldTmpFile()

	const totalSize = await getSizeOfFile(path, filename)
			
	const totalNumberChunks = Math.ceil(totalSize / CHUNK_SIZE)

	// Generate a specific token for each temporary file
	const token = await generateToken()

	// Start fetch and read each chunk of file
	let start = 0
	let currentIndexChunk = 1
	while(start < totalSize){
		let end = Math.min(start + CHUNK_SIZE, totalSize)
		await downloadChunk(path, filename, token, start, end)
		progressHandler(currentIndexChunk / totalNumberChunks * 100)
		start = end
		currentIndexChunk++
		await sleep(TIMEOUT)
	}

	const a = document.createElement('a')
	a.href = `/${token}_${filename}`
	a.download = filename
	document.body.appendChild(a)
	a.click()
	document.body.removeChild(a)
}