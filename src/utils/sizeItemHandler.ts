export const sizeItemHandler = (size : number) => {
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