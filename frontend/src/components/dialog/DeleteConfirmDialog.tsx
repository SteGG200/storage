'use client'
import { deleteItem } from "@/lib/action/delete";
import { useAppStore } from "../providers/AppStoreProvider";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { useQueryClient } from "@tanstack/react-query";
import { AlertTriangle, Trash2 } from "lucide-react";
import { useState } from "react";

interface DeleteConfirmDialogProps {
	path: string
}

export default function DeleteConfirmDialog({
	path
}: DeleteConfirmDialogProps) {
	const { modifyingItem, isDeleteConfirmDialogOpen, setIsDeleteConfirmDialogOpen } = useAppStore((state) => state)
	const [errorMessage, setErrorMessage] = useState("")
	const [isLoading, setIsLoading] = useState(false)
	const queryClient = useQueryClient()

	const remove = async () => {
		setIsLoading(true)
		try {
			if(!modifyingItem?.name) throw new Error("There is an unexpected error")
			await deleteItem(path, modifyingItem?.name)
			await queryClient.invalidateQueries({
				queryKey: ["get", path]
			})
			setIsLoading(false)
			setIsDeleteConfirmDialogOpen(false)
		}catch(err){
			setErrorMessage((err as Error).message)
			setIsLoading(false)
		}
	}

	return (
		<Dialog open={isDeleteConfirmDialogOpen} onOpenChange={setIsDeleteConfirmDialogOpen}>
			<DialogContent>
				<DialogTitle className="flex items-center">
					<AlertTriangle className="mr-2 h-5 w-5 text-red-400" />
					Confirm Delete
				</DialogTitle>
				<div className="space-y-4">
					<div className="p-4 bg-red-900/30 border border-red-700 rounded-md">
						<p className="text-sm text-red-400">Are you sure you want to remove "{modifyingItem?.name}" ?</p>
						<p className="text-xs text-gray-400 mt-2">
							{modifyingItem?.isDirectory ? "This folder and all its contents will be permanently deleted." : "This file will be permanently deleted."}
						</p>
						<p className="text-xs text-gray-400 mt-1">This action cannot be undone.</p>
					</div>
					{errorMessage && <p className="text-sm text-red-500">{errorMessage}</p>}
					<div className="flex space-x-2">
						<Button variant="outline" className="flex-1 text-gray-300 border-gray-700 hover:bg-gray-700 bg-transparent" onClick={() => setIsDeleteConfirmDialogOpen(false)} disabled={isLoading}>Cancel</Button>
						<Button className="flex-1 bg-red-700 hover:bg-red-800 text-white" onClick={remove} disabled={isLoading}>
							<Trash2 className="mr-2 h-4 w-4"/>
							Remove
						</Button>
					</div>
				</div>
			</DialogContent>
		</Dialog>
	)
}