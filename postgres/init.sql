CREATE DATABASE matchbook;

CREATE SCHEMA IF NOT EXISTS football;

CREATE TABLE football.games (
	id VARCHAR(40) PRIMARY KEY,
	event_id VARCHAR ( 20 ) NOT NULL,
	market_id VARCHAR (40) NOT NULL,
	description VARCHAR (40) ) ;
