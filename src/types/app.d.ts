type ItemType = "file" | "folder"

interface IItem {
	type: ItemType
	name: string
	size: string
	key: string
}

interface IIconProps {
	className: string
}

interface IParams {
	params: {
		path: string[] | undefined
	}
}