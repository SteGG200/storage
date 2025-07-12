import { queryOptions } from "@tanstack/react-query"
import { listDirectory } from "."

export const listOptions = (path : string) => {
	return queryOptions({
		queryKey: ["get", path],
		queryFn: async () => await listDirectory(path)
	})
}