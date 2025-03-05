'use client'
import { Upload } from "lucide-react";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTrigger } from "../ui/dialog";
import { DialogTitle } from "@radix-ui/react-dialog";
import { FormEvent, useState } from "react";
import { Input } from "../ui/input";
import { uploadFile } from "@/lib/data/uploadFile";
import { useQueryClient } from "@tanstack/react-query";

interface UploadFileButtonProps {
	path: string
}

export default function UploadFileButton({
	path
}: UploadFileButtonProps) {
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
				<DialogTrigger className="cursor-pointer" asChild>
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
							<Input
								type="text"
								name="name"
								placeholder="New file name"
								value={newFileName}
								onChange={(event) => setNewFileName(event.target.value)}
								className="bg-gray-700 border-gray-600 text-gray-100"
							/>
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
			<Dialog open={isProgressDialogOpening} onOpenChange={setIsProgressDialogOpening}>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Uploading File</DialogTitle>
					</DialogHeader>
					<div className="space-y-4">
						<div className="relative pt-1">
							<div className="flex mb-2 items-center justify-between">
								<div>
									<span className="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-teal-600 bg-teal-200">
										Progress
									</span>
								</div>
								<div className="text-right">
									<span className="text-xs font-semibold inline-block text-teal-600">
										{currentProgress.toFixed(1)}%
									</span>
								</div>
							</div>
							<div className="overflow-hidden h-2 mb-4 text-xs flex rounded bg-teal-200">
								<div
									style={{ width: `${currentProgress}%` }}
									className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-teal-500 transition-all duration-500 ease-in-out"
								>
								</div>
							</div>
						</div>
						<p className="text-center text-sm text-gray-400">Uploading... Please wait</p>
					</div>
				</DialogContent>
			</Dialog>
		</>
	)	
}