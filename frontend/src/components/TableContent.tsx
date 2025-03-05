'use client'

import { useQuery } from "@tanstack/react-query"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table"
import { listOptions } from "@/lib/data/list/options"
import { useMemo, useState } from "react"
import { useAppStore } from "./providers/AppStoreProvider"

interface TableContentProps {
	path: string
}

export default function TableContent({
	path
}: TableContentProps) {
	const { data, error } = useQuery(listOptions(path))
	const { searchResult } = useAppStore((state) => state)
	const dataToShow = useMemo(() => {
		// If searchResult is undefined, that means user is not searching -> return default data
		if (!searchResult){
			return data
		}
		return searchResult
	}, [searchResult, data])

	return (
		<Table>
			<TableHeader>
				<TableRow className="border-b border-gray-700">
					<TableHead className="w-[300px] text-gray-300">Name</TableHead>
					<TableHead className="w-[150px] text-gray-300">Size</TableHead>
					<TableHead className="text-gray-300">Date Modified</TableHead>
					<TableHead className="w-[100px] text-gray-300">Action</TableHead>
				</TableRow>
			</TableHeader>
			<TableBody>
				{(dataToShow && dataToShow.length > 0) ? (
					dataToShow.map((item, index) => (
						<TableRow className="border-b border-gray-700" key={index}>
							<TableCell>
								{item.name}
							</TableCell>
							<TableCell>
								{item.size}
							</TableCell>
							<TableCell>
								{item.date}
							</TableCell>
						</TableRow>
					))
				) : (
					<TableRow>
						<TableCell colSpan={4} className="text-center py-4 text-gray-400">
							Not found and files or folders
						</TableCell>
					</TableRow>
				)}
			</TableBody>
		</Table>
	)
}