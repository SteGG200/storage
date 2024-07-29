type ItemType = "file" | "folder"

interface IItem {
	type: ItemType
	name: string
	size: number
	key: string
}

interface IIconProps {
	className: string
}