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
		$.ajax({
			url: '/env/active',
			type: 'PUT',
			dataType: 'json',
			error: response => {
				if (response.status === 200) {
					defer.resolve();
					return;
				}

				defer.reject();
			},
			success: () => defer.resolve(),
			data: {
				env
			}
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
			<select className="form-control" name="env" value={current} onChange={ev => this.handleChanges(ev)}>
				{envList.map(
					({ name }, index) => (<option key={name} value={name}>{name.length ? name : 'Default'}</option>)
				)}
			</select>
		)
	}
}

