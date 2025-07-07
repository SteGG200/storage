"use client";
import { useQueryClient } from "@tanstack/react-query";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";
import { Input } from "../ui/input";
import { FormEvent, useEffect, useState } from "react";
import { Button } from "../ui/button";
import { renameItem } from "@/lib/action/rename";
import { useAppStore } from "../providers/AppStoreProvider";

interface RenameDialogProps {
  path: string;
}

export default function RenameDialog({
  path,
}: RenameDialogProps) {
  const queryClient = useQueryClient();
  const { modifyingItem, isRenameDialogOpen, setIsRenameDialogOpen } = useAppStore((state) => state)
  const [newName, setNewName] = useState("");
	const [errorMessage, setErrorMessage] = useState("")

  useEffect(() => {
    if(modifyingItem?.name){
      setNewName(modifyingItem.name)
    }
  }, [modifyingItem])

	const submit = async (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault()
		try {
      if(!modifyingItem?.name) throw new Error("There is an unexpected error")
			const formData = new FormData(event.currentTarget)
			await renameItem(path, modifyingItem?.name, formData)
			await queryClient.invalidateQueries({
				queryKey: ["get", path]
			})
      setErrorMessage("")
			setIsRenameDialogOpen(false)
		}catch(e){
			const message = (e as Error).message
			setErrorMessage(message)
		}
	}

  return (
    <Dialog open={isRenameDialogOpen} onOpenChange={setIsRenameDialogOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Rename {modifyingItem?.isDirectory ? "Folder" : "File"}</DialogTitle>
        </DialogHeader>
        <form onSubmit={submit} className="space-y-6">
          <Input
            type="text"
            name="newName"
            placeholder="New Name"
            value={newName}
            onChange={(event) => {
              setNewName(event.target.value);
            }}
						className="bg-gray-700 border-gray-600 text-gray-100"
          />
					{errorMessage.length > 0 && (
						<p className="text-red-500">{errorMessage}</p>
					)}
					<Button type="submit" className="bg-green-600 hover:bg-green-700 text-white">
						Submit
					</Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
