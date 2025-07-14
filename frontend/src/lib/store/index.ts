import { createStore } from 'zustand';

const defaultState: States = {
	searchResult: [],
	isDownloading: false,
	currentDonwloadProgress: 0,
	modifyingItem: undefined,
	isRenameDialogOpen: false,
	isDeleteConfirmDialogOpen: false,
};

export const createAppStore = (initState: States = defaultState) => {
	return createStore<AppStoreProps>()((set) => {
		return {
			...initState,
			setSearchResult: (newResult: SearchResult) => {
				set(() => {
					if (!newResult) return { searchResult: undefined };
					return { searchResult: newResult.map((item) => ({ ...item })) };
				});
			},
			setIsDownloading: (state: boolean) => {
				set(() => {
					return {
						isDownloading: state,
					};
				});
			},
			setCurrentDownloadProgress: (progress: number) => {
				set(() => {
					return {
						currentDonwloadProgress: progress,
					};
				});
			},
			setModifyingItem: (item: Item | undefined) => {
				set(() => {
					return {
						modifyingItem: item,
					};
				});
			},
			setIsRenameDialogOpen: (open: boolean) => {
				set(() => {
					return {
						isRenameDialogOpen: open,
					};
				});
			},
			setIsDeleteConfirmDialogOpen: (open: boolean) => {
				set(() => {
					return {
						isDeleteConfirmDialogOpen: open,
					};
				});
			},
		};
	});
};
