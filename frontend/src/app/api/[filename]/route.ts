import fs from 'fs/promises'

export async function GET(request: Request, { params }: { params: Promise<{ filename: string }> }){
	const { filename } = await params

	const buffer = await fs.readFile(`../storage/${filename}`)

	return new Response(buffer, {
		status: 200
	})
}