import { sleep } from "@/lib/utils/client"
import { downloadChunk, generateRandomHash, getSizeOfFile, removeOldTmpFile } from "./server"

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

	// Generate a specific random hash string for each temporary file
	const hashString = await generateRandomHash()

	// Start fetch and read each chunk of file
	let start = 0
	let currentIndexChunk = 1
	while(start < totalSize){
		let end = Math.min(start + CHUNK_SIZE, totalSize)
		await downloadChunk(path, filename, hashString, start, end)
		progressHandler(currentIndexChunk / totalNumberChunks * 100)
		start = end
		currentIndexChunk++
		await sleep(TIMEOUT)
	}

	const a = document.createElement('a')
	a.href = `/${hashString}_${filename}`
	a.download = filename
	document.body.appendChild(a)
	a.click()
	document.body.removeChild(a)
}