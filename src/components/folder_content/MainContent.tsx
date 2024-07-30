'use client';

import {
	getKeyValue,
	Spinner,
	Table,
	TableBody,
	TableCell,
	TableColumn,
	TableHeader,
	TableRow
} from '@nextui-org/react';
import FolderIcon from '../icons/FolderIcon';
import FileDefault from '../icons/FileDefault';
import ActionButtonList from './ActionButtonList';
import type { Selection } from '@nextui-org/react';
import React from 'react';
import { usePathname, useRouter } from 'next/navigation';

interface IItemIconProps {
	type: ItemType;
}

interface IMainContentProps {
	path: string;
}

const ItemIcon: React.FC<IItemIconProps> = ({ type }) => {
	if (type === 'folder') {
		return <FolderIcon className="flex-none size-[20px] mr-2 max-md:mr-1 stroke-customWhite" />;
	}
	return <FileDefault className="flex-none size-[20px] mr-2 max-md:mr-1 stroke-customWhite" />;
};

const MainContent: React.FC<IMainContentProps> = ({ path }) => {
	const [items, setItems] = React.useState<IItem[]>([]);
	const [loading, setLoading] = React.useState(true);
	const [selectedKeys, setSelectedKeys] = React.useState<Selection>(new Set([]));
	const [errorTable, setErrorTable] = React.useState('Not found and files or folders');
	const [canUpdated, setCanUpdated] = React.useState(true);

	const router = useRouter();
	const pathname = usePathname();

	React.useEffect(() => {
		const getItems = async () => {
			const response = await fetch(`/api/${path}`);

			const result = await response.json();

			if (result.status === 200) {
				setItems(result.items);
			} else {
				setErrorTable(result.message);
			}

			setLoading(false);
			setCanUpdated(false);
		};

		if (canUpdated) {
			getItems();
		}
	}, [canUpdated]);

	return (
		<Table
			aria-label="Sth"
			selectionMode="single"
			selectionBehavior="replace"
			layout="fixed"
			selectedKeys={selectedKeys}
			onSelectionChange={(key) => {
				console.log(key);
				setSelectedKeys(key);
			}}
		>
			<TableHeader>
				<TableColumn key="name">Name</TableColumn>
				<TableColumn key="size">Size</TableColumn>
				<TableColumn key="action">{''}</TableColumn>
			</TableHeader>
			<TableBody
				items={items}
				emptyContent={errorTable}
				isLoading={loading}
				loadingContent={<Spinner label="Loading..." />}
			>
				{(item) => (
					<TableRow
						className="group md:h-[52px]"
						key={item.key}
						onDoubleClick={() => {
							if (item.type === 'folder') {
								router.push(`${pathname}/${encodeURI(item.name)}`);
							}
						}}
					>
						{(columnKey) => {
							if (columnKey === 'action') {
								return (
									<TableCell>
										<ActionButtonList
											path={path}
											itemname={item.name}
											type={item.type}
											onUpdate={() => {
												setCanUpdated(true);
											}}
										/>
									</TableCell>
								);
							}

							return (
								<TableCell>
									<div className="flex items-center">
										{columnKey === 'name' && <ItemIcon type={item.type} />}
										<p className="truncate">{getKeyValue(item, columnKey)}</p>
									</div>
								</TableCell>
							);
						}}
					</TableRow>
				)}
			</TableBody>
		</Table>
	);
};

export default MainContent;
