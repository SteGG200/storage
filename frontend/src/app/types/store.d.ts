type SearchResult = Item[] | undefined

interface SearchState {
	searchResult: SearchResult
}

interface SearchActions {
	setSearchResult: (newResult: SearchResult) => void
}

type AppStoreProps = SearchState & SearchActions