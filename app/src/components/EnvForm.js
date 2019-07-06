import React from 'react'
import $ from 'jquery';

const createEnvironment = (name, callback) => {
	const sanitizedInput = name.replace(/[^a-z0-9\s_-]+/gi, '').trim();

	if (!sanitizedInput.length) {
		window.$.notify({
			title: 'Ops',
			message: 'O nome do ambiente deve conter apenas: letras, números, _ e -'
		}, {
			type: 'warning'
		});

		return;
	}

	new Promise((resolve, reject) => {
		$.ajax('/env/', {
			type: 'POST',
			contentType: 'application/json',
			error: response => {
				if (response.status === 200) {
					resolve();
					return;
				}

				reject(new Error(`HTTP ${response.status}`));
			},
			success: () => resolve(),
			data: JSON.stringify({
				name: sanitizedInput
			}),
		});
	})
	.then(() => {
		window.$.notify({
			message: `Ambiente '${sanitizedInput}' criado com sucesso`
		}, {
			type: 'success'
		});

		callback();
	})
	.catch(response => {
		window.$.notify({
			title: 'Ops!',
			message: 'Ocorreu um erro ao criar o ambiente',
		}, {
			type: 'danger'
		})

		console.error('Could not create environment: %o', response);
	});
}

const EnvForm = ({ onCreate, onCancel }) => {
	const textInput = React.createRef();

	return (
		<>
			<div className="mb-3">
				<input type="text" className="form-control" ref={textInput} />
			</div>
			<div className="mb-3">
				<button
					className="btn btn-primary mr-3"
					type="button"
					onClick={ev => createEnvironment(textInput.current.value, onCreate)}
				>Criar novo ambiente</button>

				<button
					className="btn btn-secondary"
					onClick={() => onCancel()}
					type="button"
				>Cancelar</button>
			</div>
		</>
	)
};

export default EnvForm;
