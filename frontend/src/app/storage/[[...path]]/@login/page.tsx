import LoginForm from "@/components/LoginForm";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default async function LoginInterface({
	params
}: {
	params: Promise<{path?: string[]}>
}) {
	const path = ((await params).path ?? []).join('/')

	return (
		<div className="min-h-dvh flex items-center justify-center p-4">
			<Card className="w-full max-w-md bg-gray-800 border-gray-700">
				<CardHeader className="text-center">
					<CardTitle className="text-2xl font-bold text-gray-100">Storage</CardTitle>
					<CardDescription className="text-gray-400">Enter your password to access this folder</CardDescription>
				</CardHeader>
				<CardContent>
					<LoginForm path={path}/>
				</CardContent>
			</Card>
		</div>
	)
}