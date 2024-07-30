interface FormProps {
	onAction: (formData: FormData) => void
	children: React.ReactNode
}

const Form : React.FC<FormProps> = ({ onAction, children }) => {
	return (
		<form className="space-y-6" action={onAction}>
			{children}
		</form>
	)
}

export default Form