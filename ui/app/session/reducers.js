import { combineReducers } from 'redux';

function userinfo(state = {}, action) {
	if (action.type == 'SESSION_OK') {
		return { ...state, userinfo: action.userinfo };
	}
	return state;
}

const reducers = combineReducers({
	userinfo: userinfo,
});

export default reducers;
