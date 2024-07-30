'use client';

import React from 'react';
import {
	Dropdown,
	DropdownTrigger,
	DropdownMenu,
	DropdownItem,
	useDisclosure,
	Modal,
	ModalHeader,
	ModalBody,
	ModalContent,
	ModalFooter,
	Input,
	Button
} from '@nextui-org/react';
import VerticalDots from '../icons/VerticalDotsIcon';
import UploadIcon from '../icons/UploadIcon';
import NewFolderIcon from '../icons/NewFolderIcon';
import Form from '../Form';

interface IMenuProps {
	directory: string;
}

const Menu: React.FC<IMenuProps> = ({ directory }) => {
	const { isOpen, onOpen, onOpenChange } = useDisclosure();
	const [actionType, setActionType] = React.useState('');
	const [isPending, setIsPending] = React.useState(false);
	const [error, setError] = React.useState<string | null>(null);
	const [filenameError, setFilenameError] = React.useState(false);
	const [fileError, setFileError] = React.useState(false);
	const [foldernameError, setFolderNameError] = React.useState(false);

	const uploadFileAction = async (formData: FormData) => {
		const response = await fetch(`/api/upload/${directory}`, {
			method: 'POST',
			body: formData
		});

		const result = await response.json();

		if (result.status === 200) {
			return null;
		} else {
			return result.message;
		}
	};

	const createFolderAction = async (formData: FormData) => {
		const response = await fetch(`/api/${directory}`, {
			method: 'POST',
			body: formData
		});

		const result = await response.json();

		if (result.status === 200) {
			return null;
		} else {
			return result.message;
		}
	};

	const validations = (formData: FormData) => {
		let isValid = true;
		if (actionType === 'upload') {
			const filename = formData.get('name');
			const file = formData.get('file');

			if (typeof filename !== 'string' || filename === '') {
				setFilenameError(true);
				isValid = false;
			}

			if (!(file instanceof File) || (file.name == '' && file.size == 0)) {
				setFileError(true);
				isValid = false;
			}
		} else if (actionType == 'new') {
			const foldername = formData.get('name');

			if (typeof foldername !== 'string' || foldername === '') {
				setFolderNameError(true);
				isValid = false;
			}
		}
		return isValid;
	};

	const actionHandler = async (formData: FormData) => {
		const isValid = validations(formData);
		if (!isValid) {
			setIsPending(false);
			return;
		}
		if (actionType === 'upload') {
			const err = await uploadFileAction(formData);
			setError(err);
		} else if (actionType === 'new') {
			const err = await createFolderAction(formData);
			setError(err);
		} else setError('There was an error');
		setIsPending(false);
		if (!error) {
			window.location.reload();
		}
	};

	return (
		<>
			<Dropdown>
				<DropdownTrigger>
					<button
						className="size-[40px] max-md:size-[33px] flex justify-center items-center hover:bg-slate-800 p-1 rounded-md outline-none"
						title="Menu"
					>
						<VerticalDots className="size-[25px] stroke-customWhite" />
					</button>
				</DropdownTrigger>
				<DropdownMenu
					onAction={(key) => {
						onOpen();
						if (typeof key === 'string') setActionType(key);
						setError(null);
						setIsPending(false);
					}}
				>
					<DropdownItem
						key="upload"
						description="Upload your file to this folder"
						startContent={<UploadIcon className="size-[20px] stroke-customWhite" />}
					>
						Upload File
					</DropdownItem>
					<DropdownItem
						key="new"
						description="Create a new folder in this folder"
						startContent={<NewFolderIcon className="size-[20px] stroke-customWhite" />}
					>
						New Folder
					</DropdownItem>
				</DropdownMenu>
			</Dropdown>
			<Modal isOpen={isOpen} onOpenChange={onOpenChange}>
				<ModalContent>
					<ModalHeader>Upload File</ModalHeader>
					<ModalBody>
						<Form
							onAction={(formData) => {
								setIsPending(true);
								actionHandler(formData);
							}}
						>
							{actionType === 'upload' ? (
								<>
									<div className="space-y-10">
										<Input
											type="text"
											isRequired
											label="File name"
											name="name"
											placeholder="Enter your file's name"
											isDisabled={isPending}
											isInvalid={filenameError}
											onChange={() => setFilenameError(false)}
											errorMessage="Please enter a valid file's name"
										/>
										<Input
											type="file"
											isRequired
											label="File"
											name="file"
											isDisabled={isPending}
											isInvalid={fileError}
											onChange={() => setFileError(false)}
											errorMessage="Please upload a valid file"
										/>
									</div>
									{error && <p className="text-customRed">{error}</p>}
									<Button type="submit" color="primary" isLoading={isPending}>
										Submit
									</Button>
								</>
							) : (
								<>
									<Input
										type="text"
										isRequired
										label="Folder name"
										name="name"
										placeholder="Enter your folder's name"
										isDisabled={isPending}
										isInvalid={foldernameError}
										onChange={() => setFolderNameError(false)}
										errorMessage="Please enter a valid folder's name"
									/>
									{error && <p className="text-customRed">{error}</p>}
									<Button type="submit" color="primary" isLoading={isPending}>
										Submit
									</Button>
								</>
							)}
						</Form>
					</ModalBody>
					<ModalFooter></ModalFooter>
				</ModalContent>
			</Modal>
		</>
	);
};

export default Menu;
