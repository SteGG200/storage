import LoginForm from '@/components/LoginForm';
import { Button } from '@/components/ui/button';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { ArrowLeft } from 'lucide-react';
import Link from 'next/link';

export default async function LoginInterface({
	params,
}: {
	params: Promise<{ path?: string[] }>;
}) {
	const arrayPath = (await params).path;
	const path = (arrayPath ?? []).join('/');

	return (
		<div className="min-h-dvh flex items-center justify-center p-4">
			<Card className="w-full max-w-md bg-gray-800 border-gray-700">
				<CardHeader className="text-center">
					<CardTitle className="text-2xl font-bold text-gray-100">
						Storage
					</CardTitle>
					<CardDescription className="text-gray-400">
						Enter your password to access this folder
					</CardDescription>
				</CardHeader>
				<CardContent className="space-y-4">
					<LoginForm path={path} />
					{arrayPath && (
						<div className="border-t border-gray-600 pt-4">
							<Button variant="outline" className="w-full" asChild>
								<Link
									className="text-white"
									href={`/storage/${arrayPath.slice(0, arrayPath.length - 1).join('/')}`}
								>
									<ArrowLeft />
									Back
								</Link>
							</Button>
						</div>
					)}
				</CardContent>
			</Card>
		</div>
	);
}
