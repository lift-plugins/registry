export function sessionRequest() {
	return {
		type: 'SESSION_REQUEST'
	};
}

export function sessionOK(userinfo) {
	return {
		type: 'SESSION_OK',
		userinfo: userinfo,
	};
}

export function sessionError(error) {
	return {
		type: 'SESSION_ERROR',
		error: error
	};
}
