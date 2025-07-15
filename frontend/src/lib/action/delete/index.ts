'use server';
import { getToken, saveToken } from '@/lib/utils/server';
import filepath from 'path';

export const deleteItem = async (path: string, filename: string) => {
	const token = await getToken();
	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/remove/${filepath.join(path, filename)}`,
		{
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${token}`,
			},
		}
	);

	if (!response.ok) {
		if (response.status === 401) {
			throw new Error('You have to log in this folder first to remove it');
		}
		const msg = await response.text();
		throw new Error(msg);
	}

	const data: { token: string } = await response.json()
	const newToken = data.token
	await saveToken(newToken)
};
