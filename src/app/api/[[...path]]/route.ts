import fs from 'fs';
import fsPromise from 'fs/promises';
import { promisify } from 'util';
import fastFolderSize from 'fast-folder-size';
import { sizeItemHandler } from '@/utils/sizeItemHandler';

export async function GET(request: Request, { params }: IParams) {
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path?.join('/') ?? '')}`;
	if (!fs.existsSync(path)) {
		return Response.json({ status: 404, message: 'Directory does not exist!' });
	}

	const tempItems = await fsPromise.readdir(path);

	let items: IItem[] = [];

	const stats = await Promise.all(
		tempItems.map((item) => {
			return fsPromise.stat(`${path}/${item}`);
		})
	);

	const fastFolderSizeAsync = promisify(fastFolderSize);

	const sizeItems = await Promise.all(
		tempItems.map((item, index) => {
			if (stats[index].isDirectory()) return fastFolderSizeAsync(`${path}/${item}`);
			return stats[index].size;
		})
	);

	for (let index = 0; index < stats.length; index++) {
		const nameItem = tempItems[index];
		const stat = stats[index];
		const size = sizeItems[index];
		if (typeof size === 'undefined') {
			// console.log(size);
			return Response.json({ status: 400, message: 'There is thing wrong' });
		}

		items.push({
			type: stat.isDirectory() ? 'folder' : 'file',
			name: nameItem,
			size: sizeItemHandler(size),
			key: index.toString()
		});
	}

	return Response.json({ status: 200, items });
}

export async function POST(request: Request, { params }: IParams) {
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path?.join('/') ?? '')}`;

	if (!fs.existsSync(path)) {
		return Response.json({ status: 404, message: 'Directory does not exist!' });
	}

	const data = await request.formData();
	const foldername = data.get('name');

	if (typeof foldername !== 'string' || foldername === '') {
		return Response.json({ status: 400, message: 'Invalid folder name' });
	}

	if (fs.existsSync(`${path}/${foldername}`)) {
		return Response.json({ status: 409, message: 'This name already exists!' });
	}

	await fsPromise.mkdir(`${path}/${foldername}`);

	return Response.json({ status: 200 });
}

export async function PUT(request: Request, { params }: IParams) {
	if (typeof params.path === 'undefined') {
		return Response.json({ status: 400, message: 'Invalid path' });
	}
	const oldName = params.path.pop();
	if (!oldName) {
		return Response.json({ status: 400, message: 'Invalid path' });
	}
	const data = await request.formData();
	const newName = data.get('name');
	if (typeof newName !== 'string' || newName === '') {
		return Response.json({ status: 400, message: 'Invalid name' });
	}

	const parentPath = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path.join('/'))}`;
	const oldPath = `${parentPath}/${decodeURI(oldName)}`;
	const newPath = `${parentPath}/${decodeURI(newName)}`;

	if (!fs.existsSync(oldPath)) {
		return Response.json({ status: 404, message: 'Directory does not exist!' });
	}

	if (fs.existsSync(newPath)) {
		return Response.json({ status: 404, message: 'This name already exist!' });
	}

	await fsPromise.rename(oldPath, newPath);

	return Response.json({ status: 200 });
}

export async function DELETE(request: Request, { params }: IParams) {
	if (typeof params.path === 'undefined') {
		return Response.json({ status: 400, message: 'Invalid path' });
	}

	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path.join('/'))}`;

	if (!fs.existsSync(path)) {
		return Response.json({ status: 404, message: 'Directory does not exist!' });
	}

	await fsPromise.rm(path, {
		recursive: true
	});

	return Response.json({ status: 200 });
}
