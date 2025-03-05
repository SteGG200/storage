'use client'

import { useEffect, useState } from "react";
import { Input } from "./ui/input";
import { useAppStore } from "./providers/AppStoreProvider";
import { useDebounce } from "@/hooks/useDebounce";
import { useQuery } from "@tanstack/react-query";
import { searchOptions } from "@/lib/action/search/options";

interface SearchBarProps {
	path: string
}

export default function SearchBar({
	path
}: SearchBarProps){
	const [searchText, setSearchText] = useState(""); // state of user input
	const [ searchValue, setSearchValue ] = useState("") // state of searching value to fetch
	const { setSearchResult } = useAppStore((state) => state)
	const { data } = useQuery(searchOptions(path, searchValue))

	useDebounce(searchText, 500, () => {
		// After delay, set user input to searching value to fetch
		setSearchValue(searchText)
	})

	useEffect(() => {
		setSearchResult(data)
	}, [data])

	return (
		<Input
			type="text"
			value={searchText}
			onChange={(event) => {
				setSearchText(event.target.value)
			}}
			placeholder="Search files and folders..."
			className="w-1/2 bg-gray-800 border-gray-700 text-gray-100 placeholder-gray-400"
		/>
	)
}