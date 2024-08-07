interface IVerticalDotsProps {
	className: string;
}

const VerticalDotsIcon: React.FC<IVerticalDotsProps> = ({ className }) => {
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
			<path d="M12 12m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
			<path d="M12 19m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
			<path d="M12 5m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0" />
		</svg>
	);
};

export default VerticalDotsIcon;
