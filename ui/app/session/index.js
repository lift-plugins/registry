import React, { Component, PropTypes } from 'react';
import { connect } from 'react-redux';
import { sessionRequest, sessionOK, sessionError } from './actions';
import store from '../store';
import classNames from 'classnames';

// mapStateToProps maps state concerned to the Session component to its
// internal properties. This function gets executed each time the store is mutated.
const mapStateToProps = (state) => {
	let props = {};
	if (state.api) {
		props.apiReady = state.api.ready;
	}

	return props;
};

const Session = (props) => {
	let className = 'app-loading';
	if (props.apiReady) {
		className = '';
	}
	return (
			<div className={className}>
				{props.children}
			</div>
	);
};

const SessionContainer = connect(mapStateToProps)(Session);
export default SessionContainer;
