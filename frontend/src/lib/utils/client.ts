import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const sleep = async (ms: number) => {
  return new Promise(resolve => {
    setTimeout(resolve, ms)
  })
}

export const formatSize = (size: number) => {
  const units = ["B", "KB", "MB", "GB", "TB"]
  let unitIndex = 0
  let sizeInUnit = size
  while(sizeInUnit >= 1024 && unitIndex < units.length - 1){
    sizeInUnit /= 1024
    unitIndex++
  }

  return `${sizeInUnit.toFixed(1)} ${units[unitIndex]}`
}

export const formatDate = (dateStringISO: string) => {
  const date = new Date(dateStringISO)
  return date.toLocaleString("en-GB", {
    timeZone: "UTC",
  })
}