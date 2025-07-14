'use server';

import { getToken } from '@/lib/utils/server';

export const listDirectory = async (path: string) => {
	const token = await getToken();

	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/get/${path}`,
		{
			headers: {
				Authorization: `Bearer ${token}`,
			},
		}
	);

	const data: Item[] = await response.json();
	return data;
};
