import type { Config } from 'tailwindcss';
import { nextui } from '@nextui-org/react';

const config: Config = {
	content: [
		'./src/pages/**/*.{js,ts,jsx,tsx,mdx}',
		'./src/components/**/*.{js,ts,jsx,tsx,mdx}',
		'./src/app/**/*.{js,ts,jsx,tsx,mdx}',
		'./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}'
	],
	theme: {
		extend: {
			colors: {
				'customWhite': '#E3E3E3',
				'customBlack': '#131314',
				'customGray': '#848484',
				'customDarkGray': '#3a3a3a',
				'customRed': '#F31260'
			},
			content: {
				'space': '\"\u2800\"'
			},
			backgroundImage: {
				'uploadIcon': 'url("/uploadIcon.svg")'
			}
		}
	},
	darkMode: 'class',
	plugins: [nextui()]
};
export default config;
