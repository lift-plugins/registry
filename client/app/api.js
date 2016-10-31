import Swagger from 'swagger-client';
import store from './store';

let options = {
	success: () => {
		console.info("API spec loaded successfully");
		store.dispatch({
			type: 'API_READY'
		});
	},

	failure: (error) => {
		store.dispatch({
			type: 'API_LOAD_FAILED',
			error: error
		});
		console.error("Unable to load swagger spec.");
		console.error(error);
	}
};

// For some strange reason we cannot use ES6 syntax here, babel will leave the reducer undefined
export function apiReducer(state = {}, action) {
	if (action.type == 'API_READY') {
		return { ...state, ready: true };
	}
	return state;
}

export default new Swagger('/lib/api.swagger.json', options);
