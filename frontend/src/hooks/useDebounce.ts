import { useEffect } from "react"

export const useDebounce = (
	value: string, 
	delay: number,
	onDebounce: () => void
) => {
	useEffect(() => {
		const handler = setTimeout(() => {
			onDebounce()
		}, delay)

		return () => {
			clearTimeout(handler)
		}
	}, [value, delay]);
}