'use client'
import { deleteItem } from "@/lib/action/delete";
import { useAppStore } from "../providers/AppStoreProvider";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { useQueryClient } from "@tanstack/react-query";

interface DeleteConfirmDialogProps {
	path: string
}

export default function DeleteConfirmDialog({
	path
}: DeleteConfirmDialogProps) {
	const { modifyingItem, isDeleteConfirmDialogOpen, setIsDeleteConfirmDialogOpen } = useAppStore((state) => state)
	const queryClient = useQueryClient()

	const remove = async () => {
		try {
			if(!modifyingItem?.name) throw new Error("There is an unexpected error")
			await deleteItem(path, modifyingItem?.name)
			await queryClient.invalidateQueries({
				queryKey: ["get", path]
			})
			setIsDeleteConfirmDialogOpen(false)
		}catch(err){
			console.error(err)
		}
	}

	return (
		<Dialog open={isDeleteConfirmDialogOpen} onOpenChange={setIsDeleteConfirmDialogOpen}>
			<DialogContent>
				<div className="space-y-6">
					<p>Are you sure to delete this {modifyingItem?.isDirectory ? 'folder' : 'file'}</p>
					<div className="flex justify-around">
						<Button variant="destructive" onClick={remove}>Yes</Button>
						<Button variant="secondary" onClick={() => setIsDeleteConfirmDialogOpen(false)}>No</Button>
					</div>
				</div>
			</DialogContent>
		</Dialog>
	)
}