export default async function StorageLayout({
	params,
	login
}: {
	params: Promise<{path?: string[]}>,
	login: React.ReactNode
}) {
	const path = (await params).path ?? []
	return(
		<>
      <h1>Storage Page</h1>
      <p>Current Path: {path.join("/")}</p>
      {login}
    </>
	)
}