type SearchResult = Item[] | undefined;

interface SearchStates {
	searchResult: SearchResult;
}

interface SearchActions {
	setSearchResult: (newResult: SearchResult) => void;
}

interface DownloadingStates {
	isDownloading: boolean;
	currentDonwloadProgress: number;
}

interface DownloadingActions {
	setIsDownloading: (state: boolean) => void;
	setCurrentDownloadProgress: (progress: number) => void;
}

interface TableContentStates {
	modifyingItem: Item | undefined;
}

interface TableContentActions {
	setModifyingItem: (item: Item | undefined) => void;
}

interface RenameDialogStates {
	isRenameDialogOpen: boolean;
}

interface RenameDialogActions {
	setIsRenameDialogOpen: (open: boolean) => void;
}

interface DeleteConfirmDialogStates {
	isDeleteConfirmDialogOpen: boolean;
}

interface DeleteConfirmDialogActions {
	setIsDeleteConfirmDialogOpen: (open: boolean) => void;
}

type States = SearchStates &
	DownloadingStates &
	TableContentStates &
	RenameDialogStates &
	DeleteConfirmDialogStates;
type Actions = SearchActions &
	DownloadingActions &
	TableContentActions &
	RenameDialogActions &
	DeleteConfirmDialogActions;

type AppStoreProps = States & Actions;
