import NewFolderButton from "@/components/button/NewFolderButton"
import UploadFileButton from "@/components/button/UploadFileButton"
import Nav from "@/components/Nav"
import DownloadInfoDialog from "@/components/dialog/DownloadInfoDialog"
import SearchBar from "@/components/SearchBar"
import TableContent from "@/components/TableContent"
import RenameDialog from "@/components/dialog/RenameDialog"
import DeleteConfirmDialog from "@/components/dialog/DeleteConfirmDialog"
import { getQueryClient } from "@/lib/queryClient"
import { listOptions } from "@/lib/action/list/options"
import { dehydrate, HydrationBoundary } from "@tanstack/react-query"
import { checkNeedAuth } from "@/lib/action/auth"
import ManagePasswordDialog from "@/components/dialog/ManagePasswordDialog"
import SetNewPasswordDialog from "@/components/dialog/SetNewPasswordDialog"

export default async function ContentPage({
	params
}: {
	params: Promise<{path?: string[]}>
}) {
	const path = ((await params).path ?? []).join('/')
	const needAuth = await checkNeedAuth(path)

	const queryClient = getQueryClient()
	await queryClient.prefetchQuery(listOptions(path))

	return (
		<HydrationBoundary state={dehydrate(queryClient)}>
			<div className="min-h-dvh container mx-auto p-4 space-y-8">
				<h1 className="text-2xl font-bold">My Storage</h1>
				<Nav path={path}/>
				<div className="flex justify-between items-center mb-4">
					<SearchBar path={path}/>
					<div className="space-x-2">
						<UploadFileButton path={path}/>
						<NewFolderButton path={path} />
						{needAuth ? (
							<ManagePasswordDialog path={path}/>
						): (
							<SetNewPasswordDialog path={path}/>
						)}
					</div>
				</div>
				<TableContent path={path}/>
			</div>
			<DownloadInfoDialog/>
			<RenameDialog path={path}/>
			<DeleteConfirmDialog path={path}/>
		</HydrationBoundary>
	)
}