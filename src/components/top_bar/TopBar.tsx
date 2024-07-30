import React from 'react';
import NavBotton from './NavButton';
import Menu from './Menu';
import PathShowAndSearch from './PathShowAndSearch';

const TopBar: React.FC<{ directory: string[] }> = ({ directory }) => {
	return (
		<div className="flex justify-between max-md:justify-center items-center">
			<NavBotton />
			<PathShowAndSearch directory={directory} />
			<Menu directory={directory.join('/')} />
		</div>
	);
};

export default TopBar;
