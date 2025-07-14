'use client'

import { Download } from "lucide-react"
import { DropdownMenuItem } from "../ui/dropdown-menu"
import { downloadFile } from "@/lib/action/download"
import { useAppStore } from "../providers/AppStoreProvider"

interface DownloadFileButtonProps {
	path: string
	item: Item
}

export default function DownloadFileButton({
	path,
	item
}: DownloadFileButtonProps) {
	const { setIsDownloading, setCurrentDownloadProgress } = useAppStore((state) => state)

	const download = async () => {
		try {
			setIsDownloading(true)
			await downloadFile(path, item.name, (progress) => {
				setCurrentDownloadProgress(progress)
				if(progress === 100){
					setIsDownloading(false)
					setCurrentDownloadProgress(0)
				}
			})
		}catch (error) {
			console.log(`Download failed: ${error}`)
		}
	}

	return(
		<>
			<DropdownMenuItem className="hover:bg-gray-700 space-x-2" onClick={download}>
				<Download className="w-4 h-4"/>
				<span>Download</span>
			</DropdownMenuItem>
		</>
	)
}