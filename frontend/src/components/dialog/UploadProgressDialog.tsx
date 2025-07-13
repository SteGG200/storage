import { Dialog, DialogContent, DialogHeader, DialogTitle } from "../ui/dialog";

interface UploadProgressDialogProps {
	open: boolean
	progress: number
}

export default function UploadProgressDialog({
	open,
	progress
}: UploadProgressDialogProps) {
	return(
		<Dialog open={open}>
			<DialogContent className="[&>button]:hidden">
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
									{progress.toFixed(1)}%
								</span>
							</div>
						</div>
						<div className="overflow-hidden h-2 mb-4 text-xs flex rounded bg-teal-200">
							<div
								style={{ width: `${progress}%` }}
								className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-teal-500 transition-all duration-500 ease-in-out"
							>
							</div>
						</div>
					</div>
					<p className="text-center text-sm text-gray-400">Uploading... Please wait</p>
				</div>
			</DialogContent>
		</Dialog>
	)
}