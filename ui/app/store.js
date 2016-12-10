import { routerReducer } from 'react-router-redux';
import { createStore, combineReducers, applyMiddleware } from 'redux';

import { sessionReducers } from './session/reducers';
import { apiReducer } from './api';

const reducers = combineReducers({
	session: sessionReducers,
	routing: routerReducer,
	api: apiReducer
});

const logger = store => next => action => {
	console.log('dispatching', action);
	let result = next(action);
	console.log('next state', store.getState());
	return result;
};

const loggerMiddleware = applyMiddleware(logger);

let store = createStore(reducers /*, retrieve state from localstorate*/, loggerMiddleware);
if (process.NODE_ENV = 'development') {
	store = createStore(reducers, window.devToolsExtension && window.devToolsExtension(), loggerMiddleware);
}
console.log(store.getState());

export default store;
