'use client';

import { FolderPlus } from 'lucide-react';
import { Button } from '../ui/button';
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '../ui/dialog';
import { Input } from '../ui/input';
import { FormEvent, useState } from 'react';
import { createFolder } from '@/lib/action/createFolder';
import { useQueryClient } from '@tanstack/react-query';

interface NewFolderDialogProps {
	path: string;
}

export default function NewFolderDialog({ path }: NewFolderDialogProps) {
	const [newFolderName, setNewFolderName] = useState('');
	const [isDialogOpening, setIsDialogOpening] = useState(false);
	const [errorMessage, setErrorMessage] = useState('');
	const queryClient = useQueryClient();

	const submit = async (event: FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		try {
			const formData = new FormData(event.currentTarget);
			await createFolder(path, formData);
			await queryClient.invalidateQueries({
				queryKey: ['get', path],
			});
			setIsDialogOpening(false);
		} catch (e) {
			const message = (e as Error).message;
			setErrorMessage(message);
		}
	};

	return (
		<Dialog open={isDialogOpening} onOpenChange={setIsDialogOpening}>
			<DialogTrigger asChild>
				<Button className="bg-green-600 hover:bg-green-700 text-white">
					<FolderPlus className="mr-2 h-4 w-4" />
					New Folder
				</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Create New Folder</DialogTitle>
				</DialogHeader>
				<form onSubmit={submit} className="space-y-6">
					<Input
						type="text"
						name="name"
						placeholder="Folder Name"
						value={newFolderName}
						onChange={(event) => {
							setNewFolderName(event.target.value);
						}}
						className="bg-gray-700 border-gray-600 text-gray-100"
					/>
					{errorMessage.length > 0 && (
						<p className="text-red-500">{errorMessage}</p>
					)}
					<Button
						type="submit"
						className="bg-green-600 hover:bg-green-700 text-white w-full"
					>
						Create
					</Button>
				</form>
			</DialogContent>
		</Dialog>
	);
}
