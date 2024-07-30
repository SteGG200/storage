import Link from 'next/link';

export default function Home() {
	return (
		<main className="h-dvh bg-customBlack flex flex-col justify-center items-center space-y-4">
			<p className="text-4xl max-md:text-2xl px-2 text-center">Welcome to SteGG&apos;s Storage</p>
			<Link href="/storage" className="lg:text-lg flex">
				<p className="after:content-space hover:after:content-none hover:underline">
					Click here to enter the storage
				</p>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					className="size-[28px] stroke-customWhite"
					viewBox="0 0 24 24"
					strokeWidth="1.5"
					fill="none"
					strokeLinecap="round"
					strokeLinejoin="round"
				>
					<path stroke="none" d="M0 0h24v24H0z" fill="none" />
					<path d="M5 12l14 0" />
					<path d="M15 16l4 -4" />
					<path d="M15 8l4 4" />
				</svg>
			</Link>
		</main>
	);
}
