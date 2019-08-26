import React from 'react'
import $ from 'jquery';

const createEnvironment = (name, callback) => {
	const sanitizedInput = name.replace(/[^a-z0-9\s_-]+/gi, '').trim();

	if (!sanitizedInput.length) {
		window.$.notify({
			title: 'Ops',
			message: 'Only letters, numbers, _ and - are allowed in environment name'
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
	.then(
		() => {
			window.$.notify({
				message: `Environment '${sanitizedInput}' created successfully`
			}, {
				type: 'success'
			});

			callback(sanitizedInput);
		}, response => {
			window.$.notify({
				title: 'Ops!',
				message: 'An error ocurred',
			}, {
				type: 'danger'
			})

			console.error('Could not create environment: %o', response);
		}
	);
}

const EnvForm = ({ onCreate, onCancel }) => {
	const textInput = React.createRef();

	return (
		<div className="input-group">
			<input placeholder="my new env..." type="text" className="form-control" ref={textInput} />
			<div className="input-group-append">
				<button
					title="Save"
					className="btn btn-info fa fa-save"
					type="button"
					onClick={ev => createEnvironment(textInput.current.value, onCreate)}
				></button>

				<button
					title="Cancel"
					className="btn btn-dark fa fa-times"
					onClick={() => onCancel()}
					type="button"
				></button>
			</div>
		</div>
	)
};

export default EnvForm;
