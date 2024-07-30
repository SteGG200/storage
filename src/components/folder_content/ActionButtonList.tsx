import { Button, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Input, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Tooltip, useDisclosure } from "@nextui-org/react"
import DownloadIcon from "../icons/DownloadIcon"
import RenameIcon from "../icons/RenameIcon"
import RemoveIcon from "../icons/RemoveIcon"
import VerticalDotsIcon from "../icons/VerticalDotsIcon"
import FolderIcon from "../icons/FolderIcon"
import { useRouter } from "next/navigation"
import React from "react"
import Form from "../Form"

interface IActionButtonListProps {
	type: ItemType
	path: string
	itemname: string
	onUpdate: () => void
}

const test = new Promise((resolve, reject) => {
	setTimeout(() => {
		resolve("here")
	}, 10000)
})

const ActionButtonList: React.FC<IActionButtonListProps> = ({type, path, itemname, onUpdate}) => {
	const router = useRouter()
	const {isOpen, onOpen, onOpenChange} = useDisclosure()
	const [actionType, setActionType] = React.useState("")
	const [filenameError, setFilenameError] = React.useState(false)
	const [renameError, setRenameError] = React.useState(null)
	const [isPending, setIsPending] = React.useState(false)

	const accessAction = () => {
		router.push(`/storage/${path}/${itemname}`)
	}

	const downloadAction = async () => {
		const response = await fetch(`/api/download/${path}/${encodeURI(itemname)}`)
		if(response.headers.get('content-type') === 'application/json'){
			const result = await response.json()
			console.log(result)
		}else{
			const data = await response.blob()
			const tempURL = window.URL.createObjectURL(data)
			const tempLinkTag = document.createElement("a")
			tempLinkTag.href = tempURL
			tempLinkTag.download = itemname
			tempLinkTag.click()
		}
	}

	const renameAction = async (formData : FormData) => {
		const newName = formData.get('name')
		if(typeof newName !== 'string' || newName === ""){
			setFilenameError(true)
			setIsPending(false)
			return false
		}

		const response = await fetch(`/api/${path}/${itemname}`, {
			method: 'PUT',
			body: formData
		})

		const result = await response.json()

		if(result.status === 200){
			setIsPending(false)
			return true
		}

		setRenameError(result.message)
		return false
	}

	const deleteAction = async (onclose : () => void) => {
		const response = await fetch(`/api/${path}/${itemname}`, {
      method: 'DELETE'
    })

    const result = await response.json()

		setIsPending(false)

    if(result.status === 200){
      onUpdate()
			onclose()
		}
		console.log(result)
	}

	return (
		<div className="w-full flex justify-end">
			<div className="max-md:hidden">
				<div className="hidden group-hover:flex justify-end space-x-4">
					<Tooltip showArrow delay={0} closeDelay={100} content="Download">
						<button onClick={downloadAction} className="hover:bg-customDarkGray p-2 rounded-full">
							<DownloadIcon className="size-[20px] stroke-customWhite"/>
						</button>
					</Tooltip>
					<Tooltip showArrow delay={0} closeDelay={100} content="Rename">
						<button className="hover:bg-customDarkGray p-2 rounded-full" onClick={() => {
								setActionType("rename")
								onOpen()
							}}
						>
							<RenameIcon className="size-[20px] stroke-customWhite" />
						</button>
					</Tooltip>
					<Tooltip showArrow delay={0} closeDelay={100} content="Delete">
						<button className="hover:bg-customDarkGray p-2 rounded-full" onClick={() => {
								setActionType("delete")
								onOpen()
							}}
						>
							<RemoveIcon className="size-[20px] stroke-customWhite" />
						</button>
					</Tooltip>
				</div>
			</div>
			<div className="flex justify-end">
				<Dropdown>
					<DropdownTrigger>
						<button className="hover:bg-customDarkGray p-2 rounded-full">
							<VerticalDotsIcon className="size-[20px] stroke-customWhite"/>
						</button>
					</DropdownTrigger>
					<DropdownMenu aria-label="Setting" disabledKeys={type === "folder" ? [] : ["access"]}
						onAction={(key) => {
							if(key === "access"){
								accessAction()
							}else if(key === "download"){
								downloadAction()
							}else if(key === "rename"){
								setActionType("rename")
								onOpen()
							}else if(key === "delete"){
								setActionType("delete")
								onOpen()
							}
						}}
					>
						<DropdownItem
							key="access"
							startContent={<FolderIcon className="size-[20px] stroke-customWhite"/>}
						>
							Access
						</DropdownItem>
						<DropdownItem
							key="download"
							startContent={<DownloadIcon className="size-[20px] stroke-customWhite"/>}
						>
							Download
						</DropdownItem>
						<DropdownItem
							key="rename"
							startContent={<RenameIcon className="size-[20px] stroke-customWhite"/>}
						>
							Rename
						</DropdownItem>
						<DropdownItem
							key="delete"
							startContent={<RemoveIcon className="size-[20px] stroke-customWhite"/>}
						>
							Delete
						</DropdownItem>
					</DropdownMenu>
				</Dropdown>
			</div>
			<Modal isOpen={isOpen} onOpenChange={onOpenChange}>
				<ModalContent>
					{(onclose) => (
						<>
							<ModalHeader>
								{actionType === "rename" && `Rename this ${type}`}
								{actionType === "delete" && `Delete this ${type}`}
							</ModalHeader>
							<ModalBody>
								{actionType === "rename" &&
									<Form onAction={async (formData) => {
										const isValid = await renameAction(formData)
										if(isValid){
											onclose()
											onUpdate()
										}
									}}>
										<Input type="text" isRequired label="New name" name="name" defaultValue={itemname} placeholder="Enter the new name" isDisabled={isPending} isInvalid={filenameError} onChange={() => setFilenameError(false)} errorMessage="Please enter a valid name"/>
										{ renameError && <p className="text-customRed">{renameError}</p>}
										<Button type="submit" color="primary" isLoading={isPending}>Submit</Button>
									</Form>
								}
								{ actionType === "delete" &&
									<p>Are you want to delete this {type} ?</p>
								}
							</ModalBody>
							<ModalFooter>
								{	actionType === "delete" &&
									<Button color="danger" isLoading={isPending} onPress={() => {
										setIsPending(true)
										deleteAction(onclose)
									}}>
										Delete
									</Button>
								}
							</ModalFooter>
						</>
					)}
				</ModalContent>
			</Modal>
		</div>
	)
}

export default ActionButtonList