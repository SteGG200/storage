interface INewFolderIconProps {
	className: string;
}

const NewFolderIcon: React.FC<INewFolderIconProps> = ({ className }) => {
	return (
		<svg
			xmlns="http://www.w3.org/2000/svg"
			className={className}
			viewBox="0 0 24 24"
			strokeWidth="1.5"
			fill="none"
			strokeLinecap="round"
			strokeLinejoin="round"
		>
			<path stroke="none" d="M0 0h24v24H0z" fill="none" />
			<path d="M12 19h-7a2 2 0 0 1 -2 -2v-11a2 2 0 0 1 2 -2h4l3 3h7a2 2 0 0 1 2 2v3.5" />
			<path d="M16 19h6" />
			<path d="M19 16v6" />
		</svg>
	);
};

export default NewFolderIcon;
