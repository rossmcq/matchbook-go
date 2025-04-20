CREATE DATABASE matchbook;

CREATE SCHEMA IF NOT EXISTS football;

CREATE TABLE football.games (
	id VARCHAR(40) PRIMARY KEY,
	event_id VARCHAR ( 20 ) NOT NULL,
	market_id VARCHAR (40) NOT NULL,
	start_at TIMESTAMP NOT NULL,
	status VARCHAR ( 10 ) NOT NULL,
	home_team VARCHAR ( 20 ) NOT NULL,
	away_team VARCHAR ( 20 ) NOT NULL,
	description VARCHAR (90) ) ;

GRANT ALL ON TABLE football.games  TO matchbook_user;

CREATE TABLE football.match_odds (
	id VARCHAR(40) PRIMARY KEY,
	game_id VARCHAR ( 20 ) NOT NULL,
	insert_at TIMESTAMP NOT NULL,
	home_win_back_odds NUMERIC ( 8 ) NOT NULL,
	home_win_back_amount NUMERIC ( 16 ) NOT NULL,
	home_win_lay_odds NUMERIC ( 8 ) NOT NULL,
	home_win_lay_amount NUMERIC ( 16 ) NOT NULL,
	draw_back_odds NUMERIC ( 8 ) NOT NULL,
	draw_back_amount NUMERIC ( 16 ) NOT NULL,
	draw_lay_odds NUMERIC ( 8 ) NOT NULL,
	draw_lay_amount NUMERIC ( 16 ) NOT NULL,
	away_win_back_odds NUMERIC ( 8 ) NOT NULL,
	away_win_back_amount NUMERIC ( 16 ) NOT NULL,
	away_win_lay_odds NUMERIC ( 8 ) NOT NULL,
	away_win_lay_amount NUMERIC ( 16 ) NOT NULL,
	home_win_back_odds_second NUMERIC ( 8 ) NOT NULL,
	home_win_back_amount_second NUMERIC ( 16 ) NOT NULL,
	home_win_lay_odds_second NUMERIC ( 8 ) NOT NULL,
	home_win_lay_amount_second NUMERIC ( 16 ) NOT NULL,
	draw_back_odds_second NUMERIC ( 8 ) NOT NULL,
	draw_back_amount_second NUMERIC ( 16 ) NOT NULL,
	draw_lay_odds_second NUMERIC ( 8 ) NOT NULL,
	draw_lay_amount_second NUMERIC ( 16 ) NOT NULL,
	away_win_back_odds_second NUMERIC ( 8 ) NOT NULL,
	away_win_back_amount_second NUMERIC ( 16 ) NOT NULL,
	away_win_lay_odds_second NUMERIC ( 8 ) NOT NULL,
	away_win_lay_amount_second NUMERIC ( 16 ) NOT NULL ) ;

GRANT ALL ON TABLE football.match_odds  TO matchbook_user;
