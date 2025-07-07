'use server'
import filepath from 'path'

export const deleteItem = async (
	path: string,
	filename: string
) => {
	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/remove/${filepath.join(path,filename)}`, {
		method: 'DELETE'
	})

	if(!response.ok){
		const msg = await response.text()
		throw new Error(msg)
	}
}