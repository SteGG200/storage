import { Edit } from "lucide-react";
import { DropdownMenuItem } from "../ui/dropdown-menu";
import { useAppStore } from "../providers/AppStoreProvider";

interface RenameButtonProps {
	item: Item
}

export default function RenameButton({
	item,
}: RenameButtonProps) {
	const { setModifyingItem, setIsRenameDialogOpen } = useAppStore((state) => state)

	return (
		<DropdownMenuItem className="hover:bg-gray-700 space-x-2" onClick={() => {
			setModifyingItem(item)
			setIsRenameDialogOpen(true)
		}}>
			<Edit className="w-4 h-4"/>
			<span>Rename</span>
		</DropdownMenuItem>
	)
}