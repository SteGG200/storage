import { queryOptions } from "@tanstack/react-query"

export const searchOptions = (
	path: string,
  searchText: string
) => {
	return queryOptions({
		queryKey: ["search", path, searchText],
		queryFn: async (): Promise<Item[]> => {
			const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/get/search/${path}?q=${searchText}`)

			return response.json()
		},
		enabled: searchText !== ""
	})
}