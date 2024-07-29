'use client'

import React from "react"
import { Dropdown, DropdownTrigger, DropdownMenu, DropdownItem } from "@nextui-org/react"
import VerticalDots from "../icons/VerticalDotsIcon"
import UploadIcon from "../icons/UploadIcon"
import NewFolderIcon from "../icons/NewFolderIcon"

const Menu = () => {


	return(
		<Dropdown>
			<DropdownTrigger>
				<button className="size-[40px] max-md:size-[33px] flex justify-center items-center hover:bg-slate-800 p-1 rounded-md outline-none" title="Menu">
					<VerticalDots className="size-[25px] stroke-customWhite"/>
				</button>
			</DropdownTrigger>
			<DropdownMenu onAction={(key) => {
				console.log(key)
			}}>
				<DropdownItem key="upload"
				  description="Upload your file to this folder"
					startContent={<UploadIcon className="size-[20px] stroke-customWhite"/>}
				>
					Upload File
				</DropdownItem>
				<DropdownItem key="new"
					description="Create a new folder in this folder"
					startContent={<NewFolderIcon className="size-[20px] stroke-customWhite"/>}
				>
					New Folder
				</DropdownItem>
			</DropdownMenu>
		</Dropdown>
	)
}

export default Menu