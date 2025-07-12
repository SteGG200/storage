import { queryOptions } from "@tanstack/react-query"
import { checkIsAuthorized } from "."

export const checkIsAuthorizedOptions = (path : string) => {
	return queryOptions({
		queryKey: ['checkAuth', path],
		queryFn: async () => {
			return await checkIsAuthorized(path)
		}
	})
}