import type { Metadata } from 'next';
import { Geist } from 'next/font/google';
import './globals.css';
import QueryProvider from '@/components/providers/QueryProvider';
import AppStoreProvider from '@/components/providers/AppStoreProvider';

const geistSans = Geist({
	variable: '--font-geist-sans',
	subsets: ['latin'],
});

export const metadata: Metadata = {
	title: 'Storage',
	description: 'Simple storage management app',
};

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html lang="en">
			<body className={`${geistSans.className} dark antialiased`}>
				<QueryProvider>
					<AppStoreProvider>{children}</AppStoreProvider>
				</QueryProvider>
			</body>
		</html>
	);
}
