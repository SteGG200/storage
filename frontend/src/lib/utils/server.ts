'use server';

import { cookies } from 'next/headers';

export const saveToken = async (token: string) => {
	const cookieStore = await cookies();
	cookieStore.set('token', token);
};

export const getToken = async () => {
	const cookieStore = await cookies();
	const data = cookieStore.get('token');
	return data?.value ?? '';
};
