'use client'

import { useQuery } from "@tanstack/react-query"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table"
import { listOptions } from "@/lib/action/list/options"
import { useMemo, useState } from "react"
import { useAppStore } from "./providers/AppStoreProvider"
import { File, Folder, MoreVertical, Trash2 } from "lucide-react"
import Link from "next/link"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "./ui/dropdown-menu"
import { Button } from "./ui/button"
import DownloadFileButton from "./button/DownloadFileButton"
import RenameButton from "./button/RenameButton"
import RenameDialog from "./dialog/RenameDialog"
import React from "react"

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
	const [selectedItem, setSelectedItem] = useState<Item | undefined>(undefined)
	const [isRenameDialogOpen, setIsRenameDialogOpen] = useState(false)

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
						<React.Fragment key={index}>
							<TableRow className="border-b border-gray-700">
								<TableCell>
									<div className="flex items-center space-x-2">
										{item.isDirectory ? (
											<>
											<Folder className="h-4 w-4 text-blue-400"/>
												<Link className="hover:underline" href={`/storage/${item.path}`}>
													{item.name}
												</Link>
											</>
										) : (
											<>
												<File className="h-4 w-4"/>
												<p>{item.name}</p>
											</>
										)}
									</div>
								</TableCell>
								<TableCell>
									{item.size}
								</TableCell>
								<TableCell>
									{item.date}
								</TableCell>
								<TableCell>
									<DropdownMenu>
										<DropdownMenuTrigger asChild>
											<Button variant="ghost" className="h-8 w-8 ring-0">
												<span className="sr-only">Open menu</span>
												<MoreVertical className="h-4 w-4"/>
											</Button>
										</DropdownMenuTrigger>
										<DropdownMenuContent align="end" className="bg-gray-800 text-gray-100 border-gray-700">
											{!item.isDirectory && (
												<DownloadFileButton path={path} item={item}/>
											)}
											<RenameButton item={item}/>
											<DropdownMenuItem className="hover:bg-gray-700 space-x-2">
												<Trash2 className="w-4 h-4"/>
												<span>Delete</span>
											</DropdownMenuItem>
										</DropdownMenuContent>
									</DropdownMenu>
								</TableCell>
							</TableRow>
						</React.Fragment>
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