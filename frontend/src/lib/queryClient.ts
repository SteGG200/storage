import { isServer, QueryClient } from "@tanstack/react-query"

const makeQueryClient = () => {
	return new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: Infinity,
				retry: false,
			}
		}
	})
}

let browserQueryClient : QueryClient | null = null

export const getQueryClient = () => {
	if (isServer) {
		return makeQueryClient()
	}else {
		if (!browserQueryClient) {
      browserQueryClient = makeQueryClient()
    }
    return browserQueryClient
	}
}