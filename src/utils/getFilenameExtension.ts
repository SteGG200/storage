export function getFilenameExtension(filename: string) {
	let extension = '';
	let temp = '';
	for (let i = filename.length - 1; i >= 0; i--) {
		if (filename[i].match(/[a-zA-z]/i)) {
			temp = filename[i] + temp;
		} else if (filename[i] === '.') {
			if (temp === '') break;
			extension = '.' + temp + extension;
			temp = '';
		} else break;
	}
	// console.log(extension);
	return extension;
}
