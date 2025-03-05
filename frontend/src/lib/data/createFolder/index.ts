export const createFolder = async (
	path: string,
	formData: FormData
) => {
	const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/create/${path}`, {
		method: 'POST',
    body: formData
	})

	if (!response.ok) {
    const msg = await response.text()
    throw new Error(msg)
  }
}