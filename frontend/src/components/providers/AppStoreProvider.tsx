'use client';

import { createAppStore } from '@/lib/store';
import { createContext, useContext, useRef } from 'react';
import { useStore } from 'zustand';

type AppStoreApi = ReturnType<typeof createAppStore>;

const AppStoreContext = createContext<AppStoreApi | null>(null);

export default function AppStoreProvider({
	children,
}: {
	children: React.ReactNode;
}) {
	const appStoreRef = useRef<AppStoreApi>(null);

	if (!appStoreRef.current) {
		appStoreRef.current = createAppStore();
	}

	return (
		<AppStoreContext.Provider value={appStoreRef.current}>
			{children}
		</AppStoreContext.Provider>
	);
}

export const useAppStore = <T,>(selector: (state: AppStoreProps) => T): T => {
	const appStore = useContext(AppStoreContext);

	if (!appStore) {
		throw new Error('useAppStore must be used within an AppStoreProvider');
	}

	return useStore(appStore, selector);
};
