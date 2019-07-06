import EnvPicker from './EnvPicker';
import EnvForm from './EnvForm';
import React from 'react'
import {RecordForm} from './RecordForm.js'
import {RecordTable} from './RecordTable.js'

export class Home extends React.Component {
	constructor() {
		super();
		this.state = {
			forceUpdate: null,
			createEnv: false,
			env: ''
		};
	}

	toggleEnvForm() {
		const { createEnv } = this.state;

		this.setState({
			createEnv: !createEnv
		});
	}

	onUpdate(){
		this.table.reloadTable(this.state.env);
		this.envPicker.reload(this.state.env);
	}

	onChangeEnv(env) {
		console.debug('c=Home, m=onChangeEnv, env=%s', env);
		this.setState(
			{ env },
			() => this.table.reloadTable(this.state.env)
		)
	}

	render(){
		console.debug('render state=%o', this.state)
		return (
			<div>
				<nav className="navbar navbar-inverse navbar-fixed-top" >
					<a className="navbar-brand" href="#">DNS Proxy Server</a>
					<button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav"
						aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
						<span className="navbar-toggler-icon"></span>
					</button>
					<div className="collapse navbar-collapse" id="navbarNav">
						<ul className="navbar-nav pull-right">
							<li className="nav-item active">
								<a className="nav-link" href="#">Home <span className="sr-only">(current)</span></a>
							</li>
							<li className="nav-item">
								<a className="nav-link" href="#">Features</a>
							</li>
							<li className="nav-item">
								<a className="nav-link" href="#">Pricing</a>
							</li>
							<li className="nav-item">
								<a className="nav-link disabled" href="#">Disabled</a>
							</li>
						</ul>
					</div>
				</nav>
				<div className="container">
					<div className="col-sm-12 col-md-7 col-lg-5 mb-4">
						<h3>Environments</h3>
						{this.state.createEnv
							? <EnvForm
									onCancel={() => this.toggleEnvForm()}
									onClick={() => this.toggleEnvForm()}
								/>
							: <EnvPicker
									onChange={env => this.onChangeEnv(env)}
									onToggle={() => this.toggleEnvForm()}
									ref={(it) => this.envPicker = it}
								/>
						}
					</div>

					<div className="col-sm-12">
						<h3>New Record</h3>
						<RecordForm env={this.state.env} onUpdate={(e) => this.onUpdate(e)} />
						<RecordTable env={this.state.env} ref={(it) => this.table = it} />
					</div>
				</div>
			</div>
		);
	}
}
