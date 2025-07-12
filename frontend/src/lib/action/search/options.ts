import { queryOptions } from "@tanstack/react-query"
import { searchAction } from "."

export const searchOptions = (
	path: string,
  searchValue: string
) => {
	return queryOptions({
		queryKey: ["search", path, searchValue],
		queryFn: async () =>  await searchAction(path, searchValue),
		enabled: searchValue !== ""
	})
}