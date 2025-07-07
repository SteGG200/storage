type SearchResult = Item[] | undefined

interface SearchState {
	searchResult: SearchResult
}

interface SearchActions {
	setSearchResult: (newResult: SearchResult) => void
}

interface DownloadingState {
	isDownloading: boolean
	currentDonwloadProgress: number
}

interface DownloadingActions {
	setIsDownloading: (state: boolean) => void
	setCurrentDownloadProgress: (progress: number) => void
}

type States = SearchState & DownloadingState
type Actions = SearchActions & DownloadingActions

type AppStoreProps = States & Actions