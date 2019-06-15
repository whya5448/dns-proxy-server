import React from 'react'
import $ from 'jquery';

export default class EnvPicker extends React.Component {
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

		return $.ajax({
			url: '/env/active',
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				name: env
			})
			}).then((...args) => console.log(args));
	}

	componentDidMount() {
		this.reload();
	}

	handleChanges(ev) {
		const current = ev.target.options[ev.target.selectedIndex].value;
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
					({ name }, index) => (<option key={`${index}${name}`} value={name}>{name.length ? name : 'Default'}</option>)
				)}
			</select>
		)
	}
}


