'use client';

import Link from 'next/link';
import React from 'react';

interface IFolderLinkProps {
	isCurrentFolder: boolean;
	href: string;
	children: React.ReactNode;
}

interface IDirectoryProps {
	directory: string[];
}

const FolderLink: React.FC<IFolderLinkProps> = ({ isCurrentFolder, children, href }) => {
	if (isCurrentFolder) {
		return <div className="flex items-center px-2 whitespace-nowrap">{children}</div>;
	}

	return (
		<Link
			className="flex items-center hover:bg-[#777777] text-customGray hover:text-customWhite whitespace-nowrap group rounded mr-0.5 px-2"
			href={href}
		>
			{children}
		</Link>
	);
};

const DirectoryList: React.FC<IDirectoryProps> = ({ directory }) => {
	let tempPath = '';

	return (
		<div className="w-fit min-w-full h-full flex items-center space-x-0.5 absolute right-0 p-0.5">
			<FolderLink isCurrentFolder={directory.length == 0} href="/storage">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					className="size-[20px] stroke-customGray group-hover:stroke-customWhite mr-1"
					viewBox="0 0 24 24"
					strokeWidth="1.5"
					fill="none"
					strokeLinecap="round"
					strokeLinejoin="round"
				>
					<path stroke="none" d="M0 0h24v24H0z" fill="none" />
					<path d="M5 12l-2 0l9 -9l9 9l-2 0" />
					<path d="M5 12v7a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-7" />
					<path d="M9 21v-6a2 2 0 0 1 2 -2h2a2 2 0 0 1 2 2v6" />
				</svg>
				Home
			</FolderLink>
			{directory.map((folder, index) => {
				tempPath += `${folder}/`;

				return (
					<React.Fragment key={index}>
						/
						<FolderLink
							isCurrentFolder={index === directory.length - 1}
							href={`/storage/${tempPath}`}
						>
							{decodeURI(folder)}
						</FolderLink>
					</React.Fragment>
				);
			})}
		</div>
	);
};

const PathShowAndSearch: React.FC<IDirectoryProps> = ({ directory }) => {
	const [isSearching, setIsSearching] = React.useState(false);

	const handleKeyEvent = (event: KeyboardEvent) => {
		if (event.key === 'Escape') {
			setIsSearching(false);
		} else if (
			!event.ctrlKey &&
			!event.altKey &&
			!event.metaKey &&
			event.key.length === 1 &&
			(('A' <= event.key && event.key <= 'Z') ||
				('a' <= event.key && event.key <= 'z') ||
				('0' <= event.key && event.key <= '9') ||
				('!' <= event.key && event.key <= '/') ||
				(':' <= event.key && event.key <= '@') ||
				('[' <= event.key && event.key <= '`') ||
				('{' <= event.key && event.key <= '~'))
		) {
			setIsSearching(true);
		}
	};

	React.useEffect(() => {
		document.addEventListener('keydown', handleKeyEvent, false);

		return () => {
			document.removeEventListener('keydown', handleKeyEvent, false);
		};
	}, []);

	return (
		<div className="w-1/2 max-md:w-11/12 h-[47.5px] p-2 flex space-x-2">
			{isSearching ? (
				<input
					className="w-full h-full bg-customDarkGray bg-uploadIcon bg-no-repeat bg-[5px_5px] rounded-md outline-none pl-9"
					type="text"
					placeholder="Search..."
					autoFocus
				/>
			) : (
				<div className="w-full h-full bg-customDarkGray rounded-md overflow-hidden relative">
					<DirectoryList directory={directory} />
				</div>
			)}
			<button
				className={`${isSearching ? 'bg-slate-700' : 'bg-inherit md:hover:bg-slate-800'} p-1 rounded-md outline-none`}
				onClick={() => setIsSearching(!isSearching)}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					className="size-[25px] stroke-customWhite"
					viewBox="0 0 24 24"
					strokeWidth="1.5"
					fill="none"
					strokeLinecap="round"
					strokeLinejoin="round"
				>
					<path stroke="none" d="M0 0h24v24H0z" fill="none" />
					<path d="M10 10m-7 0a7 7 0 1 0 14 0a7 7 0 1 0 -14 0" />
					<path d="M21 21l-6 -6" />
				</svg>
			</button>
		</div>
	);
};

export default PathShowAndSearch;
