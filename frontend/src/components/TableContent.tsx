'use client'

import { useQuery } from "@tanstack/react-query"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "./ui/table"
import { listOptions } from "@/lib/data/list/options"

interface TableContentProps {
	path: string
}

export default function TableContent({
	path
}: TableContentProps) {
	const { data, error } = useQuery(listOptions(path))

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
				{(data && data.length > 0) ? (
					data.map((item, index) => (
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