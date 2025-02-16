import { queryOptions } from "@tanstack/react-query"

export const listOptions = (path : string) => {
	return queryOptions({
		queryKey: ["get", path],
		queryFn: async () : Promise<Item[]> => {
			const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/get/${path}`)

			if (!response.ok){
        throw response.status
      }

			return response.json();
		},
	})
}