'use server';
import { getToken } from '@/lib/utils/server';
import filepath from 'path';

export const renameItem = async (
	path: string,
	filename: string,
	formData: FormData
) => {
	const token = await getToken();
	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/rename/${filepath.join(path, filename)}`,
		{
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`,
			},
			body: formData,
		}
	);

	if (!response.ok) {
		const msg = await response.text();
		throw new Error(msg);
	}
};
