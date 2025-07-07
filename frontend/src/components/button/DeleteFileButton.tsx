import { Trash2 } from "lucide-react";
import { DropdownMenuItem } from "../ui/dropdown-menu";
import { useAppStore } from "../providers/AppStoreProvider";

interface DeleteFileButtonProps {
	item: Item
}

export default function DeleteFileButton({
	item
}: DeleteFileButtonProps) {
	const { setModifyingItem, setIsDeleteConfirmDialogOpen } = useAppStore((state) => state)

	return (
		<DropdownMenuItem className="hover:bg-gray-700 space-x-2" onClick={() => {
			setModifyingItem(item)
			setIsDeleteConfirmDialogOpen(true)
		}}>
			<Trash2 className="w-4 h-4"/>
			<span>Delete</span>
		</DropdownMenuItem>
	)
}