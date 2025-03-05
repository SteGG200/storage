import { createStore } from "zustand"

const defaultState : SearchState = {
	searchResult: []
}

export const createAppStore = (
	initState: SearchState = defaultState
) => {
	return createStore<AppStoreProps>()((set) => {
		return {
			...initState,
			setSearchResult: (newResult: SearchResult) => {
				set(() => {
					if (!newResult) return { searchResult: undefined}
					return { searchResult: newResult.map((item) => ({...item}))}
				})
			}
		}
	})
}