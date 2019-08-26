import React from 'react'
import $ from 'jquery';

const DEFAULT_ENV = '';

export default class EnvPicker extends React.PureComponent {
	constructor(...args) {
		super(...args);

		const { env: current = '' } = this.props;

		this.state = {
			envList: [],
			current
		};
	}

	reload() {
		return $.ajax({
			url: '/env/',
		}).then(data => {
			console.debug('c=EnvPicker, m=getData, data=%o', data);
			this.setState({envList: data});
		}, function (err) {
			console.error('c=EnvPicker, m=getData, status=error', err);
		});
	}

	deleteCurrent() {
		const { state: { current: env } } = this;
		console.log('c=EnvPicker, m=deleteCurrent, env=%s', env);

		const defer = $.Deferred();

		if (env === DEFAULT_ENV) {
			const message = 'Deleting default environment is not allowed';

			window.$.notify({
				title: 'Ops!',
				message
			}, {
				type: 'danger'
			});

			defer.rejectWith(new Error(message));
			return defer.promise();
		}

		$.ajax('/env/', {
			type: 'DELETE',
			contentType: 'application/json',
			error: ({ status }) => {
				if (status === 200) {
					defer.resolve();
					return;
				}

				defer.rejectWith(new Error(`HTTP ${status}`));
			},
			success: () => defer.resolve(),
			data: JSON.stringify({
				name: env
			})
		});

		return defer.promise();
	}

	componentDidMount() {
		this.reload();
	}

	handleChanges(ev) {
		const { target: { options, selectedIndex } } = ev;
		const current = options[selectedIndex].value;

		this.setState(
			{ current },
			() => this.props.onChange(current)
		);
	}

	render() {
		const { envList, current } = this.state;
		const deleteEnv = () => {
			this.deleteCurrent()
				.then(() => this.props.onDelete())
				.fail(err => console.error('m=render, err=%o', err))
		};

		console.debug('c=EnvPicker, m=render, env=%s', current);

		return (
			<div className="input-group">
				<select className="form-control"
					onChange={ev => this.handleChanges(ev)}
					value={current}
					name="env"
				>
					{envList.map(
						({ name }, index) => (<option key={name} value={name}>{name.length ? name : 'Default'}</option>)
					)}
				</select>
				<div className="input-group-append">
					<button
						title="Create new env"
						onClick={() => this.props.onToggle()}
						className="btn btn-info"
						type="button"
					>
						<span className="fa fa-plus-circle"></span>
					</button>
					<button
						title="Delete selected env"
						onClick={deleteEnv}
						className="btn btn-danger"
						type="button"
					>
						<span className="fa fa-trash-alt"></span>
					</button>
				</div>
			</div>
		)
	}
}

