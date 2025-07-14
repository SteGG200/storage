'use client';

import { ChevronRight, Edit, Shield, ShieldOff } from 'lucide-react';
import { Button } from '../ui/button';
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '../ui/dialog';
import { useState } from 'react';
import ChangePasswordDialog from './ChangePasswordDialog';
import RemovePasswordDialog from './RemovePasswordDialog';

interface ManagePasswordDialogProps {
	path: string;
}

export default function ManagePasswordDialog({
	path,
}: ManagePasswordDialogProps) {
	const [isManagePasswordDialogOpen, setIsManagePasswordDialogOpen] =
		useState(false);
	const [isChangePasswordDialogOpen, setIsChangePasswordDialogOpen] =
		useState(false);
	const [isRemovePasswordDialogOpen, setIsRemovePasswordDialogOpen] =
		useState(false);

	const openChangePasswordDialog = () => {
		setIsManagePasswordDialogOpen(false);
		setIsChangePasswordDialogOpen(true);
	};

	const openRemovePasswordDialog = () => {
		setIsManagePasswordDialogOpen(false);
		setIsRemovePasswordDialogOpen(true);
	};

	return (
		<>
			<Dialog
				open={isManagePasswordDialogOpen}
				onOpenChange={setIsManagePasswordDialogOpen}
			>
				<DialogTrigger asChild>
					<Button className="bg-orange-600 hover:bg-orange-700 text-white">
						<Shield className="mr-2 h-4 w-4" />
						Manage Password
					</Button>
				</DialogTrigger>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Manage Password for this folder</DialogTitle>
					</DialogHeader>
					<div className="space-y-4">
						<div className="flex items-center p-3 bg-green-900/30 border border-green-700 rounded-md">
							<Shield className="h-4 w-4 text-green-400 mr-2" />
							<span className="text-sm text-green-400">
								This directory is currently password protected
							</span>
						</div>
						<div className="space-y-3">
							<Button
								className="bg-orange-400 hover:bg-orange-500 text-white w-full flex items-center justify-between"
								onClick={openChangePasswordDialog}
							>
								<div className="flex items-center">
									<Edit className="mr-2 h-4 w-4" />
									Change Password
								</div>
								<ChevronRight className="h-4 w-4" />
							</Button>
							<Button
								className="bg-red-700 hover:bg-red-800 text-white w-full flex items-center justify-between"
								onClick={openRemovePasswordDialog}
							>
								<div className="flex items-center">
									<ShieldOff className="mr-2 h-4 w-4" />
									Remove Password Protection
								</div>
								<ChevronRight className="h-4 w-4" />
							</Button>
						</div>
					</div>
				</DialogContent>
			</Dialog>
			<ChangePasswordDialog
				path={path}
				open={isChangePasswordDialogOpen}
				onOpenChange={setIsChangePasswordDialogOpen}
				onBackAction={setIsManagePasswordDialogOpen}
			/>
			<RemovePasswordDialog
				path={path}
				open={isRemovePasswordDialogOpen}
				onOpenChange={setIsRemovePasswordDialogOpen}
				onBackAction={setIsManagePasswordDialogOpen}
			/>
		</>
	);
}
