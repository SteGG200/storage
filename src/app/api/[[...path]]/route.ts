import fs from 'fs'
import fsPromise from 'fs/promises'
import path from 'path'

interface IParams {
	params: {
		path: string[] | undefined
	}
}

const sizeItemHandler = (size : number) => {
	const maxDecimals = 2
	const units = ['B', 'KB', 'MB', 'GB']
	let currentUnitIndex = 0
	let currentSize = size
	while(currentSize >= 1000 && currentUnitIndex < units.length - 1){
		currentSize /= 1024
		currentUnitIndex++
	}

	return `${Math.round(currentSize * (10 ** maxDecimals)) / 10 ** maxDecimals} ${units[currentUnitIndex]}`
}

export async function GET(request: Request, { params }: IParams){
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${params.path?.join('/') ?? ""}`
	if(!fs.existsSync(path)){
		return Response.json({ status: 404, message: "Directory does not exist!"})
	}

	const tempItems = await fsPromise.readdir(path)

	let items: IItem[] = []

	for (let index = 0; index < tempItems.length; index++){
		const item = tempItems[index]
		const stat = await fsPromise.stat(`${path}/${item}`)
		items.push({
			type: stat.isDirectory()? "folder" : "file",
      name: item,
      size: sizeItemHandler(stat.size),
      key: index.toString()
		})
	}

	return Response.json({ status: 200, items })
}