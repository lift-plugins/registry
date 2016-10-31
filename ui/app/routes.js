import React, { Component } from 'react';
import { Router, Route, browserHistory, IndexRoute } from 'react-router';
import { Provider } from 'react-redux';
import { syncHistoryWithStore } from 'react-router-redux';

import NotFound from './not-found';
import SearchPage from './search';
import Session from './session';
import store from './store';

// Google Analytics
import ga from 'react-ga';
ga.initialize('UA-000000-01');

function logPageView() {
	ga.pageview(window.location.pathname);
}

const history = syncHistoryWithStore(browserHistory, store);
history.listen(location => logPageView(location.pathname));

const routes =
<Provider store={store}>
	<Router history={history}>
		<Route path="/" component={Session}>
			<Route path="search" component={SearchPage} />
		</Route>
		<Route path="*" component={NotFound}/>
	</Router>
</Provider>;

export default routes;
