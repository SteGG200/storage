'use client'

import { ArrowLeft } from "lucide-react"
import { Button } from "./ui/button"
import { Breadcrumb, BreadcrumbEllipsis, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "./ui/breadcrumb"
import Link from "next/link"
import { Fragment } from "react"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "./ui/dropdown-menu"

interface NavProps {
	path: string
}

export default function Nav({
	path
}: NavProps) {
	const items = path === "" ? [] : path.split("/")

	items.unshift("My Storage")

	const getHref = (index: number) : string => {
		let href = "/storage"
		for (let i = 1; i <= index; i++){
			href += "/" + items[i]
		}
	
		return href
	}

	const maxItemsToDisplay = 5
	const leftIndexCollapse = 1
	const rightIndexCollapse = items.length - maxItemsToDisplay + leftIndexCollapse

	return (
		<div className="flex items-center space-x-4">
			{/* Back to parent directory */}
			<Button
				variant="outline"
				className="text-gray-300 border-gray-700 hover:bg-gray-800"
				asChild={items.length > 1}
				disabled={items.length === 1}
			>
				{items.length > 1 ? (
					<Link href={getHref(items.length - 2)} >
						<ArrowLeft className="mr-2 h-4 w-4"/> Back
					</Link>
				) : (
					<>
						<ArrowLeft className="mr-2 h-4 w-4"/> Back
					</>
				)}
			</Button>

			{/* Display current path */}
			<Breadcrumb>
			  <BreadcrumbList>
					{items.map((item, index) => {
						if ((items.length > maxItemsToDisplay && (index < leftIndexCollapse || index >= rightIndexCollapse)) || items.length <= maxItemsToDisplay){
							return (
								<Fragment key={index}>
									<BreadcrumbItem>
										{ items.length > maxItemsToDisplay && index == rightIndexCollapse ? (
											<DropdownMenu>
												<DropdownMenuTrigger className="outline-none cursor-pointer">
													<BreadcrumbEllipsis/>
												</DropdownMenuTrigger>
												<DropdownMenuContent className="bg-popover">
													{items.map((item, index) => {
														if (leftIndexCollapse <= index && index <= rightIndexCollapse) {
															return (
																<DropdownMenuItem className="cursor-pointer" asChild>
																	<Link href={getHref(index)}>
																		{item}
																	</Link>
																</DropdownMenuItem>
															)
														}
													})}
												</DropdownMenuContent>
											</DropdownMenu>
										) : index === items.length - 1 ? (
											<BreadcrumbPage className="text-blue-400">
												{item}
											</BreadcrumbPage>
										) : (
											<BreadcrumbLink className="text-blue-400 hover:text-blue-300" asChild>
												<Link href={getHref(index)}>
													{item}
												</Link>
											</BreadcrumbLink>
										)}
									</BreadcrumbItem>
									{index < items.length - 1 && <BreadcrumbSeparator/>}
								</Fragment>
							)
						}
					})}
				</BreadcrumbList>
			</Breadcrumb>
		</div>
	)
}