'use client'

import { ArrowLeft } from "lucide-react"
import { Button } from "../ui/button"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog"
import { Input } from "../ui/input"
import { FormEvent, useState } from "react"
import { setPassword } from "@/lib/action/auth"

interface ChangePasswordDialogProps {
	path: string
	open: boolean
	onOpenChange: (open: boolean) => void
	onBackAction: (open: boolean) => void
}

export default function ChangePasswordDialog({
	path,
	open,
	onOpenChange,
	onBackAction
}: ChangePasswordDialogProps){
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

			const newPassword = formData.get("password")
			const confirmPassword = formData.get("confirm")
			if(newPassword != confirmPassword){
				throw new Error("Passwords do not match")
			}
			formData.delete("confirm")

			await setPassword(path, formData)

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
					<DialogTitle>Change Password for this folder</DialogTitle>
				</DialogHeader>
				<form className="space-y-4" onSubmit={submit}>
					<Input
						name="oldPassword"
						type="password"
						placeholder="Current Password"
						className="bg-gray-700 border-gray-600 text-gray-100"
						disabled={isLoading}
						required
					/>
					<Input
						name="password"
						type="password"
						placeholder="New Password"
						className="bg-gray-700 border-gray-600 text-gray-100"
						disabled={isLoading}
						required
					/>
					<Input
						name="confirm"
						type="password"
						placeholder="Confirm Password"
						className="bg-gray-700 border-gray-600 text-gray-100"
						disabled={isLoading}
						required
					/>
					{errorMessage && <p className="text-sm text-red-400">{errorMessage}</p>}
					<div className="flex space-x-2">
						<Button type="button" variant="outline" className="flex-1 text-gray-300 border-gray-700 hover:bg-gray-700 bg-transparent" onClick={backAction} disabled={isLoading}>
							<ArrowLeft className="mr-2 h-4 w-4"/>
							Back
						</Button>
						<Button type="submit" className="flex-1 bg-orange-600 hover:bg-orange-700 text-white" disabled={isLoading}>
							Update Password
						</Button>
					</div>
				</form>
			</DialogContent>
		</Dialog>
	)
}