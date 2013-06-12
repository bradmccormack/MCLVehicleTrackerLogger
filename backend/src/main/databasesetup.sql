BEGIN EXCLUSIVE TRANSACTION;

DROP TABLE if exists GPSRECORDS;
DROP TABLE if exists ERRORS;
DROP TABLE if exists NETWORK;
DROP TABLE if exists Company;
DROP TABLE if exists User;
DROP TABLE if exists Settings;


CREATE TABLE GPSRecords (
	id integer primary key autoincrement,
	Message text,
	Latitude text not null,
	Longitude text not null,
	Speed integer not null,
	Heading float not null,
	Fix boolean not null,
	DateTime date not null default current_timestamp,
	BusID text not null
	);
	
CREATE TABLE Errors (
	id integer primary key autoincrement,
	GPSRecordID integer not null,
	Error text,
	DateTime date not null default current_timestamp,
	FOREIGN KEY (GPSRecordID) REFERENCES GPSrecords(id)
);

CREATE TABLE Network (
	id integer primary key autoincrement,
	GPSRecordID integer not null,
	Acknowledge boolean not null default 0,
	FOREIGN KEY (GPSRecordID) REFERENCES GPSRecords(id)
);

CREATE TABLE Company (
	ID integer primary key autoincrement,
	Name text not null,
	Expiry date not null default current_timestamp,
	MaxUsers integer not null default 0
);



CREATE TABLE User (
	ID integer primary key autoincrement,
	FirstName text not null,
	LastName text not null,
	CompanyID integer not null,
	Password text not null,
	AccessLevel integer not null default 0,
	FOREIGN KEY (CompanyID) REFERENCES Company(ID)
);



CREATE TABLE Settings (
	ID integer primary key autoincrement,
	UserID integer not null,
	MapAPI text not null default 'GoogleMaps',
	FOREIGN KEY (UserID) REFERENCES User(ID)	
);

INSERT INTO Company (Name, Expiry, MaxUsers) VALUES ('myClubLink', '', 1);
INSERT INTO User (FirstName, LastName, CompanyID, Password, AccessLevel) VALUES ('guest','user', 1, '', 0);
INSERT INTO Settings (UserID, MapAPI) VALUES (1, 'GoogleMaps');


COMMIT TRANSACTION;
