'use server';

import { getToken, saveToken } from '@/lib/utils/server';
import { redirect } from 'next/navigation';

export const checkIsAuthorized = async (
	path: string
): Promise<
	| {
			isAuthorized: true;
	  }
	| {
			isAuthorized: false;
			unauthorizedPath: string;
	  }
> => {
	const token = await getToken();

	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/auth/check/${path}`,
		{
			headers: {
				Authorization: `Bearer ${token}`,
			},
		}
	);

	if (!response.ok) {
		if (response.status === 401) {
			const { unauthorizedPath } = await response.json();
			return {
				isAuthorized: false,
				unauthorizedPath,
			};
		} else {
			const message = await response.text();
			throw new Error(message);
		}
	}

	return {
		isAuthorized: true,
	};
};

export const checkNeedAuth = async (path: string) => {
	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/auth/checkNeedAuth/${path}`
	);

	if (!response.ok) {
		const message = await response.text();
		throw new Error(message);
	}

	const data: { needAuth: boolean } = await response.json();

	return data.needAuth;
};

export const login = async (path: string, formData: FormData) => {
	const token = await getToken();

	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/auth/login/${path}`,
		{
			method: 'POST',
			headers: {
				Authorization: `Bearer ${token}`,
			},
			body: formData,
		}
	);

	if (!response.ok) {
		const message = await response.text();
		throw new Error(message);
	}

	const data: { token: string } = await response.json();
	await saveToken(data.token);

	redirect(`/storage/${path}`);
};

export const setPassword = async (path: string, formData: FormData) => {
	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/auth/setPassword/${path}`,
		{
			method: 'POST',
			body: formData,
		}
	);

	if (!response.ok) {
		const message = await response.text();
		throw new Error(message);
	}

	redirect(`/storage/${path}`);
};

export const removePassword = async (path: string, formData: FormData) => {
	const response = await fetch(
		`${process.env.NEXT_PUBLIC_SERVER_URL}/auth/removePassword/${path}`,
		{
			method: 'DELETE',
			body: formData,
		}
	);

	if (!response.ok) {
		const message = await response.text();
		throw new Error(message);
	}

	redirect(`/storage/${path}`);
};
