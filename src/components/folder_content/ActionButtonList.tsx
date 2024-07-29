import { Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Tooltip } from "@nextui-org/react"
import DownloadIcon from "../icons/DownloadIcon"
import RenameIcon from "../icons/RenameIcon"
import RemoveIcon from "../icons/RemoveIcon"
import VerticalDotsIcon from "../icons/VerticalDotsIcon"
import FolderIcon from "../icons/FolderIcon"

interface IActionButtonListProps {
	accessible: boolean
}

const ActionButtonList: React.FC<IActionButtonListProps> = ({accessible}) => {
	return (
		<div className="w-full flex justify-end">
			<div className="max-md:hidden">
				<div className="hidden group-hover:flex justify-end space-x-4">
					<Tooltip showArrow delay={0} closeDelay={100} content="Download">
						<button className="hover:bg-customDarkGray p-2 rounded-full">
							<DownloadIcon className="size-[20px] stroke-customWhite"/>
						</button>
					</Tooltip>
					<Tooltip showArrow delay={0} closeDelay={100} content="Rename">
						<button className="hover:bg-customDarkGray p-2 rounded-full">
							<RenameIcon className="size-[20px] stroke-customWhite" />
						</button>
					</Tooltip>
					<Tooltip showArrow delay={0} closeDelay={100} content="Delete">
						<button className="hover:bg-customDarkGray p-2 rounded-full">
							<RemoveIcon className="size-[20px] stroke-customWhite" />
						</button>
					</Tooltip>
				</div>
			</div>
			<div className="flex justify-end">
				<Dropdown>
					<DropdownTrigger>
						<button className="hover:bg-customDarkGray p-2 rounded-full">
							<VerticalDotsIcon className="size-[20px] stroke-customWhite"/>
						</button>
					</DropdownTrigger>
					<DropdownMenu aria-label="Setting" disabledKeys={accessible ? [] : ["access"]}>
						<DropdownItem
							key="access"
							startContent={<FolderIcon className="size-[20px] stroke-customWhite"/>}
						>
							Access
						</DropdownItem>
						<DropdownItem
							key="download"
							startContent={<DownloadIcon className="size-[20px] stroke-customWhite"/>}
						>
							Download
						</DropdownItem>
						<DropdownItem
							key="rename"
							startContent={<RenameIcon className="size-[20px] stroke-customWhite"/>}
						>
							Rename
						</DropdownItem>
						<DropdownItem
							key="delete"
							startContent={<RemoveIcon className="size-[20px] stroke-customWhite"/>}
						>
							Delete
						</DropdownItem>
					</DropdownMenu>
				</Dropdown>
			</div>
		</div>
	)
}

export default ActionButtonList