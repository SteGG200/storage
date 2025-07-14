import { sleep } from '@/lib/utils/client';
import { getUploadSessionToken, uploadByChunk } from './server';

const CHUNK_SIZE = 1024 * 1024 * 5;

export const uploadFile = async (
	path: string,
	formData: FormData,
	progressHandler: (currentProgress: number) => void
) => {
	const filename = formData.get('name');
	const file = formData.get('file');

	if (!filename || !file) {
		throw new Error('Filename or file not provided');
	}

	if (typeof filename !== 'string' || !(file instanceof File)) {
		throw new Error('Invalid filename or file');
	}

	// Get JWT token of upload session
	const requestTokenFormData = new FormData();
	requestTokenFormData.append('name', filename);

	const uploadSessionToken = await getUploadSessionToken(
		path,
		requestTokenFormData
	);

	// Upload file progress
	const uploadFileFormData = new FormData();
	uploadFileFormData.append('name', filename);

	const size = file.size;
	const totalNumberChunks = Math.ceil(size / CHUNK_SIZE);

	uploadFileFormData.append('file', '');
	uploadFileFormData.append('isLast', '0');

	let startIndex = 0;
	let currentIndexChunk = 1;
	while (startIndex < size) {
		const endIndex = Math.min(startIndex + CHUNK_SIZE, size);
		const chunk = file.slice(startIndex, endIndex);

		uploadFileFormData.set('file', chunk);
		uploadFileFormData.set('isLast', endIndex === size ? '1' : '0');

		await uploadByChunk(uploadFileFormData, uploadSessionToken);

		progressHandler((currentIndexChunk / totalNumberChunks) * 100);
		startIndex = endIndex;
		currentIndexChunk++;
		await sleep(200);
	}
};
