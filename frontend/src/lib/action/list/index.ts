'use server'

import type { QueryClient } from "@tanstack/react-query"
import { listOptions } from "./options"

// Return false whenever the query is unauthorized
export const listDirectory = async (queryClient: QueryClient, path: string) : Promise<boolean> =>  {
	try{
		await queryClient.fetchQuery(listOptions(path))
		return true
	}catch(statusCode) {
		if (statusCode === 401){
			return false
		}
		return true
	}
}