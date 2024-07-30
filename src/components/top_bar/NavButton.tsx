'use client';

import { useRouter } from 'next/navigation';
import React from 'react';

interface ICurrentPath {
	path: string[];
}

const NavBotton = () => {
	const router = useRouter();

	const goBack = () => {
		router.back();
	};

	const goForward = () => {
		router.forward();
	};

	return (
		<div className="max-md:hidden space-x-2">
			<button
				className="group enabled:hover:bg-slate-800 p-1 rounded-md outline-none"
				onClick={goBack}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					className="size-[37px] stroke-customWhite group-disabled:stroke-customDarkGray"
					viewBox="0 0 24 24"
					strokeWidth="1.5"
					fill="none"
					strokeLinecap="round"
					strokeLinejoin="round"
				>
					<path stroke="none" d="M0 0h24v24H0z" fill="none" />
					<path d="M20 15h-8v3.586a1 1 0 0 1 -1.707 .707l-6.586 -6.586a1 1 0 0 1 0 -1.414l6.586 -6.586a1 1 0 0 1 1.707 .707v3.586h8a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1z" />
				</svg>
			</button>
			<button
				className="group enabled:hover:bg-slate-800 p-1 rounded-md outline-none"
				onClick={goForward}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					className="size-[37px] stroke-customWhite group-disabled:stroke-customGray"
					viewBox="0 0 24 24"
					strokeWidth="1.5"
					fill="none"
					strokeLinecap="round"
					strokeLinejoin="round"
				>
					<path stroke="none" d="M0 0h24v24H0z" fill="none" />
					<path d="M4 9h8v-3.586a1 1 0 0 1 1.707 -.707l6.586 6.586a1 1 0 0 1 0 1.414l-6.586 6.586a1 1 0 0 1 -1.707 -.707v-3.586h-8a1 1 0 0 1 -1 -1v-4a1 1 0 0 1 1 -1z" />
				</svg>
			</button>
		</div>
	);
};

export default NavBotton;
