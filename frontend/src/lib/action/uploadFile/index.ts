import { sleep } from "@/lib/utils/client"

const CHUNK_SIZE = 1024 * 1024 * 5

export const uploadFile = async (path: string, 
	formData: FormData, 
	progressHandler: (currentProgress: number) => void
) => {
	const filename = formData.get('name')
	const file = formData.get('file')

	if (!filename ||!file) {
		throw new Error('Filename or file not provided')
  }

	if (typeof filename !== 'string' || !(file instanceof File)){
		throw new Error('Invalid filename or file')
	}

	// Get JWT token of upload session
	const requestTokenFormData = new FormData()
	requestTokenFormData.append('name', filename)

	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/upload/token/${path}`, {
		method: 'POST',
		body: requestTokenFormData
	})

	if (!response.ok) {
		const msg = await response.text()
    throw new Error(msg)
  }

	const { token } : { token: string } = await response.json()


	// Upload file progress
	const uploadFileFormData = new FormData()
	uploadFileFormData.append('name', filename)
	
	const size = file.size
	const totalNumberChunks = Math.ceil(size / CHUNK_SIZE)

	uploadFileFormData.append('file', "")
	uploadFileFormData.append('isLast', "0")
	
	let startIndex = 0
	let currentIndexChunk = 1
	while(startIndex < size) {
		const endIndex = Math.min(startIndex + CHUNK_SIZE, size)
    const chunk = file.slice(startIndex, endIndex)

    uploadFileFormData.set('file', chunk)
		uploadFileFormData.set('isLast', endIndex === size? "1" : "0")

		const uploadResponse = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/upload/file`, {
			method: 'POST',
			headers: {
        'Authorization': `Bearer ${token}`
      },
      body: uploadFileFormData
		})

		if (!uploadResponse.ok) {
      throw new Error('Failed to upload file')
    }

		progressHandler(currentIndexChunk / totalNumberChunks * 100)
		startIndex = endIndex
		currentIndexChunk++
		await sleep(500)
	}
}