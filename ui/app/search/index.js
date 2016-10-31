import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Link } from 'react-router';


export default class SearchPage extends Component {
	render() {
		return (
			<section className="search_page row align-center align-middle">
				<div className="column">
					<header className="login_logo-header"></header>
					<Search/>
					<SearchResults/>
				</div>
			</section>
		);
	}
}
