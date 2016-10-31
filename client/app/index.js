import ReactDOM from 'react-dom';
import routes from './routes';
import store from './store';

import './stylesheets/app.scss';

import api from './api';

const containerEl = document.getElementById("app");

// We need to wait for our API client to load before attempting to render all the components.
store.subscribe(() => {
	let state = store.getState();

	if (!state.api || !state.api.ready) {
		return;
	}

	console.info('Booting application...');

	ReactDOM.render(routes, containerEl);

	if (module.hot) {
		module.hot.accept('./store', () => {
			store.replaceReducer(require('./store') /*.default if you use Babel 6+ */ );
		});

		module.hot.accept('./routes', function() {
			ReactDOM.unmountComponentAtNode(containerEl);
			let routes = require('./routes').default;
			ReactDOM.render(routes, containerEl);
		});
	}
});
