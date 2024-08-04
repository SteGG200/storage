import fsPromise from 'fs/promises';
import fs from 'fs';

export async function POST(request: Request, { params }: IParams) {
	const data = await request.formData();
	const path = `${process.cwd()}/${process.env.UPLOAD_PATH}/${decodeURI(params.path?.join('/') ?? '')}`;

	const filename = data.get('name');
	const file = data.get('file');

	if (typeof filename !== 'string' || !(file instanceof File)) {
		return Response.json({ status: 400, message: 'Invalid' });
	}

	if (fs.existsSync(`${path}/${filename}`)) {
		return Response.json({ status: 409, message: 'This name already exists' });
	}

	const buffer = Buffer.from(await file.arrayBuffer());

	await fsPromise.writeFile(`${path}/${filename}`, buffer);

	return Response.json({ status: 200, message: 'ok' });
}
