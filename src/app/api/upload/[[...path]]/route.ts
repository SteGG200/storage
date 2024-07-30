import fsPromise from 'fs/promises'
import fs from 'fs'

const getFileExtension = (filename: string) => {
	let extension = ""
	let canAdd = false
	for (const character of filename){
		if(character == "." && !canAdd) canAdd = true
		if(canAdd) extension += character
	}

	return extension
}

export async function POST(request: Request, { params }: IParams){
	const data = await request.formData()
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${params.path?.join('/') ?? ""}`

	const name = data.get('name')
	const file = data.get('file')

	if(typeof name !== 'string' || !(file instanceof File)){
		return Response.json({ status: 400, message: "Invalid"})
	}

	if(fs.existsSync(`${path}/${name}${getFileExtension(file.name)}`)){
		return Response.json({ status: 409, message: "File already exists"})
	}

	const buffer = Buffer.from(await file.arrayBuffer())

	await fsPromise.writeFile(`${path}/${name}${getFileExtension(file.name)}`, buffer)

	return Response.json({ status: 200, message: "ok"})
}