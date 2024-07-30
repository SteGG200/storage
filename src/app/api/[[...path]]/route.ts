import fs from 'fs'
import fsPromise from 'fs/promises'
import { promisify } from 'util'
import fastFolderSize from 'fast-folder-size'

const sizeItemHandler = (size : number) => {
	const maxDecimals = 2
	const units = ['B', 'KB', 'MB', 'GB']
	let currentUnitIndex = 0
	let currentSize = size
	while(currentSize >= 1000 && currentUnitIndex < units.length - 1){
		currentSize /= 1000
		currentUnitIndex++
	}

	return `${Math.round(currentSize * (10 ** maxDecimals)) / 10 ** maxDecimals} ${units[currentUnitIndex]}`
}

export async function GET(request: Request, { params }: IParams){
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path?.join('/') ?? "")}`
	if(!fs.existsSync(path)){
		return Response.json({ status: 404, message: "Directory does not exist!"})
	}

	const tempItems = await fsPromise.readdir(path)

	let items: IItem[] = []

	const stats = await Promise.all(tempItems.map((item) => {
		return fsPromise.stat(`${path}/${item}`)
	}))

	const fastFolderSizeAsync = promisify(fastFolderSize)

	const sizeItems = await Promise.all(tempItems.map((item, index) => {
		if(stats[index].isDirectory()) return fastFolderSizeAsync(`${path}/${item}`)
		return stats[index].size
	}))

	for(let index = 0; index < stats.length; index++){
		const nameItem = tempItems[index]
		const stat = stats[index]
		const size = sizeItems[index]
		if(typeof size === 'undefined'){
			console.log(size)
			return Response.json({ status: 400, message: "There is thing wrong"})
		}

		items.push({
			type: stat.isDirectory() ? 'folder' : 'file',
			name: nameItem,
			size: sizeItemHandler(size),
			key: index.toString()
		})
	}

	return Response.json({ status: 200, items })
}

export async function POST(request: Request, { params }: IParams){
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path?.join('/') ?? "")}`
	
	if(!fs.existsSync(path)){
		return Response.json({ status: 404, message: "Directory does not exist!"})
	}

	const data = await request.formData()
	const foldername = data.get('name')

	if(fs.existsSync(`${path}/${foldername}`)){
		return Response.json({ status: 409, message: "Folder already exists!"})
	}

	await fsPromise.mkdir(`${path}/${foldername}`)

	return Response.json({ status: 200 })
}