'use client'

import { getQueryClient } from "@/lib/queryClient"
import { QueryClientProvider } from "@tanstack/react-query"
import { ReactQueryDevtools } from "@tanstack/react-query-devtools"

export default function QueryProvider({
	children
}: {
	children: React.ReactNode
}){
	const queryClient = getQueryClient()

	return (
		<QueryClientProvider client={queryClient}>
			{children}
			{process.env.NODE_ENV === "development" && (
				<ReactQueryDevtools/>
			)}
		</QueryClientProvider>
	)
}