'use client'

import { FormEvent, useState } from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog"
import { Input } from "../ui/input"
import { Button } from "../ui/button"
import { ArrowLeft, ShieldOff } from "lucide-react"
import { removePassword } from "@/lib/action/auth"

interface RemovePasswordDialogProps {
	path: string
	open: boolean
	onOpenChange: (open: boolean) => void
	onBackAction: (open: boolean) => void
}

export default function RemovePasswordDialog({
	path,
	open,
	onOpenChange,
	onBackAction
}: RemovePasswordDialogProps){
	const [isLoading, setIsLoading] = useState(false)
	const [errorMessage, setErrorMessage] = useState("")

	const backAction = () => {
		onOpenChange(false)
		onBackAction(true)
	}

	const submit = async (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault()
		setIsLoading(true)
		try{
			const formData = new FormData(event.currentTarget)

			await removePassword(path, formData)

			setIsLoading(false)
			onOpenChange(false)
		}catch(err){
			setErrorMessage((err as Error).message)
			setIsLoading(false)
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Remove Password Protection</DialogTitle>
				</DialogHeader>
				<form onSubmit={submit} className="space-y-4">
					<div className="p-3 bg-red-900/30 border border-red-700 rounded-md">
						<p className="text-sm text-red-400">
							Warning: This will remove password protection from this folder. Anyone will be able to
							access this directory.
						</p>
					</div>
					<Input
						name="oldPassword"
						type="password"
						placeholder="Enter Current Password to Confirm"
						disabled={isLoading}
						required
					/>
					{errorMessage && <p className="text-sm text-red-400">{errorMessage}</p>}
					<div className="flex space-x-2">
						<Button type="button" variant="outline" className="flex-1 text-gray-300 border-gray-700 hover:bg-gray-700 bg-transparent" onClick={backAction} disabled={isLoading}>
							<ArrowLeft className="mr-2 h-4 w-4"/>
							Back
						</Button>
						<Button type="submit" className="flex-1 bg-red-700 hover:bg-red-800 text-white" disabled={isLoading}>
							<ShieldOff className="mr-2 h-4 w-4" />
							Remove Password
						</Button>
					</div>
				</form>
			</DialogContent>
		</Dialog>
	)
}