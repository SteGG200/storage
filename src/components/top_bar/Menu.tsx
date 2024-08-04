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

const test = new Promise((resolve, reject) => {
	setTimeout(() => {
    resolve('here');
  }, 10000);
})

const Menu: React.FC<IMenuProps> = ({ directory }) => {
	const { isOpen, onOpen, onOpenChange } = useDisclosure();
	const [actionType, setActionType] = React.useState('');
	const [isPending, setIsPending] = React.useState(false);
	const [error, setError] = React.useState<string | null>(null);

	// State for upload file action
	const [fileName, setFileName] = React.useState("")
	const [filenameError, setFilenameError] = React.useState(false);
	const [fileError, setFileError] = React.useState(false);

	// State for create folder action
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
			return result.message as string;
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
			return result.message as string;
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
		setIsPending(true);
		const isValid = validations(formData);
		if (!isValid) {
			setIsPending(false);
			return;
		}
		let err: string | null;
		if (actionType === 'upload') {
			err = await uploadFileAction(formData);
		} else if (actionType === 'new') {
			err = await createFolderAction(formData);
		} else {
			setError('Invalid action type');
			setIsPending(false);
			return;
		}
		if (err) {
			const testM = await test
			console.log(testM)
			setError(err);
			setIsPending(false);
			return;
		}
		window.location.reload();
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
											value={fileName}
											isDisabled={isPending}
											isInvalid={filenameError}
											onChange={(event) => {
												setFileName(event.target.value)
												setFilenameError(false)
											}}
											errorMessage="Please enter a valid file's name"
										/>
										<Input
											type="file"
											isRequired
											label="File"
											name="file"
											isDisabled={isPending}
											isInvalid={fileError}
											onChange={(event) => {
												if(!event.target.files || event.target.files.length > 1){
													setFileError(true)
													return
												}
												setFileName(event.target.files[0].name)
												setFileError(false)
												setFilenameError(false)
											}}
											errorMessage="Please upload a valid file"
										/>
									</div>
								</>
							) : (
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
							)}
							{error && <p className="text-customRed">{error}</p>}
							<Button type="submit" color="primary" isLoading={isPending}>
								Submit
							</Button>
						</Form>
					</ModalBody>
					<ModalFooter></ModalFooter>
				</ModalContent>
			</Modal>
		</>
	);
};

export default Menu;
