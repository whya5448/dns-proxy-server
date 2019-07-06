import React from 'react'
import $ from 'jquery';

export default class EnvPicker extends React.PureComponent {
	constructor() {
		super();
		this.state = {
			envList: [],
			current: ''
		};
	}

	reload() {
		return $.ajax({
			url: '/env/',
		}).then(data => {
			console.debug('m=getData, data=%o', data);
			this.setState({envList: data});
		}, function (err) {
			console.error('m=getData, status=error', err);
		});
	}

	activate() {
		const { state: { current: env } } = this;
		console.log('c=EnvPicker, m=activate, env=%s', env);

		const defer = $.Deferred();

		// API is returning an empty body with content type 'application/json', this is causing a parse error
		$.ajax('/env/active', {
			type: 'PUT',
			contentType: 'application/json',
			error: response => {
				if (response.status === 200) {
					defer.resolve();
					return;
				}

				defer.reject();
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
			() => this.activate()
				.fail(err => console.error('c=EnvPicker, m=handleChanges, error=%o', err))
				.always(() => this.props.onChange(current))
		);
	}

	render() {
		const { envList, current } = this.state;

		return (
			<>
				<div className="input-group">
					<select className="form-control mr-3"
						onChange={ev => this.handleChanges(ev)}
						value={current}
						name="env"
					>
						{envList.map(
							({ name }, index) => (<option key={name} value={name}>{name.length ? name : 'Default'}</option>)
						)}
					</select>
					<button
						onClick={() => this.props.onToggle()}
						className="btn btn-secondary"
						type="button"
					>
					Criar novo ambiente
					</button>
				</div>
			</>
		)
	}
}

