import fs from 'fs'
import fsPromises from 'fs/promises'

export async function GET(request: Request, { params }: IParams){
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path?.join('/')?? "")}`
  if(!fs.existsSync(path)){
    return Response.json({ status: 404, message: "Directory does not exist!"})
  }

	const stat = await fsPromises.stat(path)

	if(stat.isDirectory()){
		return Response.json({ status: 400, message: "Cannot download a directory" })
	}else{
		const data = await fsPromises.readFile(path)
		return new Response(data)
	}
}