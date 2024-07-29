"use client"

import { getKeyValue, Table, TableBody, TableCell, TableColumn, TableHeader, TableRow } from "@nextui-org/react"
import FolderIcon from "../icons/FolderIcon"
import FileDefault from "../icons/FileDefault"
import ActionButtonList from "./ActionButtonList"
import type { Selection } from "@nextui-org/react"
import React from "react"
import { usePathname, useRouter } from "next/navigation"

interface IItemIconProp {
	type: ItemType
}

interface IMainContentProps {
	items: IItem[]
}

const ItemIcon : React.FC<IItemIconProp> = ({ type }) => {
	if(type === "folder"){
		return <FolderIcon className="flex-none size-[20px] mr-2 max-md:mr-1 stroke-customWhite"/>
	}
	return <FileDefault className="flex-none size-[20px] mr-2 max-md:mr-1 stroke-customWhite"/>
}

const MainContent : React.FC<IMainContentProps> = ({ items }) => {
	const [selectedKeys, setSelectedKeys] = React.useState<Selection>(new Set([]));

	const router = useRouter()
	const pathname = usePathname()

	return (
		<Table
			aria-label="Sth"
			selectionMode="single"
			selectionBehavior="replace"
			layout="fixed"
			selectedKeys={selectedKeys}
			onSelectionChange={(key) => {
				console.log(key)
				setSelectedKeys(key)
			}}
		>
			<TableHeader>
				<TableColumn key="name">
					Name
				</TableColumn>
				<TableColumn key="size">
					Size
				</TableColumn>
				<TableColumn key="action">
					{""}
				</TableColumn>
			</TableHeader>
			<TableBody items={items}>
				{(item) => (
					<TableRow className="group md:h-[52px]" key={item.key} onDoubleClick={() => {
						if(item.type === "folder"){
							router.push(`${pathname}/${item.name}`)
						}
					}}>
						{(columnKey) => {
							if(columnKey === "action"){
								return (
									<TableCell>
										<ActionButtonList accessible={item.type === "folder"}/>
									</TableCell>
								)
							}

							return (
								<TableCell>
									<div className="flex items-center">
										{ columnKey === "name" && <ItemIcon type={item.type} /> }
										<p className="truncate">{getKeyValue(item, columnKey)}</p>
									</div>
								</TableCell>
							)
						}}
					</TableRow>
				)}
			</TableBody>
		</Table>
	)
}

export default MainContent