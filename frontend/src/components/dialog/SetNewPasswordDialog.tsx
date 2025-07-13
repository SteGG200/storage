'use client'

import { Shield } from "lucide-react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "../ui/dialog"
import { Button } from "../ui/button"
import { FormEvent, useState } from "react"
import { Input } from "../ui/input"
import { setPassword } from "@/lib/action/auth"

interface SetNewPasswordDialogProps {
	path: string
}

export default function SetNewPasswordDialog({
	path
}: SetNewPasswordDialogProps) {
	const [errorMessage, setErrorMessage] = useState("")
	const [isLoading, setIsLoading] = useState(false)
	const [isDialogOpen, setIsDialogOpen] = useState(false)

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
			setIsDialogOpen(false)

		}catch(err){
			setErrorMessage((err as Error).message)
			setIsLoading(false)
		}
	}

	return (
		<Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
			<DialogTrigger asChild>
				<Button className="bg-orange-600 hover:bg-orange-700 text-white">
					<Shield className="mr-2 h-4 w-4"/>
					Set Password
				</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Set Password for this folder</DialogTitle>
				</DialogHeader>
				<form onSubmit={submit} className="space-y-4">
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
					<Button type="submit" className="bg-orange-600 hover:bg-orange-700 text-white w-full" disabled={isLoading}>
						Set Password
					</Button>
				</form>
			</DialogContent>
		</Dialog>
	)
}