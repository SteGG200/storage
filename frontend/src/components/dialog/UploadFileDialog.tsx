'use client'
import { Upload } from "lucide-react";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTrigger } from "../ui/dialog";
import { DialogTitle } from "@radix-ui/react-dialog";
import { FormEvent, useState } from "react";
import { Input } from "../ui/input";
import { uploadFile } from "@/lib/action/uploadFile";
import { useQueryClient } from "@tanstack/react-query";
import UploadProgressDialog from "./UploadProgressDialog";

interface UploadFileDialogProps {
	path: string
}

export default function UploadFileButton({
	path
}: UploadFileDialogProps) {
	const [selectedFile, setSelectedFile] = useState<File | null>(null)
	const [newFileName, setNewFileName] = useState<string>("")
	const [isDialogOpening, setIsDialogOpening] = useState<boolean>(false)
	const [isProgressDialogOpening, setIsProgressDialogOpening] = useState<boolean>(false)
	const [currentProgress, setCurrentProgress] = useState<number>(0)
	const [errorMessage, setErrorMessage] = useState<string>('')
	const queryClient = useQueryClient()

	const submit = async (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault()
		const formData = new FormData(event.currentTarget)
		setIsDialogOpening(false)
		setIsProgressDialogOpening(true)
		try {
			await uploadFile(path, formData, (progress) => {
				setCurrentProgress(progress)
				if (progress === 100){
					setIsProgressDialogOpening(false)
				}
			})
			setSelectedFile(null)
			setNewFileName("")
			await queryClient.invalidateQueries({
				queryKey: ["get", path]
			})
		}catch (err) {
			const msg = (err as Error).message
			setIsProgressDialogOpening(false)
			setIsDialogOpening(true)
      setErrorMessage(msg)
		}
	}

	return (
		<>
			<Dialog open={isDialogOpening} onOpenChange={setIsDialogOpening}>
				<DialogTrigger asChild>
					<Button className="bg-blue-600 hover:bg-blue-700 text-white">
						<Upload className="mr-2 h-4 w-4"/>
						Upload File
					</Button>
				</DialogTrigger>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Upload File</DialogTitle>
					</DialogHeader>
					<form onSubmit={submit} className="space-y-5">
						<Input
							type="file"
							name="file"
							onChange={(event) => {
								if(event.target.files && event.target.files.length > 0) {
									const file = event.target.files[0];
									setSelectedFile(file)
									setNewFileName(file.name)
								}
							}}
							className="bg-gray-700 border-gray-600 text-gray-100"
						/>
						{selectedFile && (
							<>
								<label className="text-sm font-medium text-gray-300">
									Enter your file name:
								</label>
								<Input
									type="text"
									name="name"
									placeholder="New file name"
									value={newFileName}
									onChange={(event) => setNewFileName(event.target.value)}
									className="mt-2 bg-gray-700 border-gray-600 text-gray-100"
								/>
							</>
						)}
						{ errorMessage.length > 0 && (
							<p className="text-red-500">{errorMessage}</p>
						)}
						<Button
							type="submit"
							disabled={!selectedFile || !newFileName.trim()}
							className="bg-blue-600 hover:bg-blue-700 text-white"
						>
							Upload
						</Button>
					</form>
				</DialogContent>
			</Dialog>
			<UploadProgressDialog open={isProgressDialogOpening} progress={currentProgress}/>
		</>
	)	
}