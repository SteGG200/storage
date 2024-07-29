import React from "react";
import NavBotton from "./NavButton";
import Menu from "./Menu";

const TopBar : React.FC<{ children : React.ReactNode }> = ({ children }) => {
	return (
		<div className="flex justify-between max-md:justify-center items-center">
			<NavBotton/>
			{children}
			<Menu/>
		</div>
	)
}


export default TopBar