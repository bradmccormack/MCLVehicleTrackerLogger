BEGIN EXCLUSIVE TRANSACTION;

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

CREATE TABLE Company (
  ID integer primary key autoincrement,
  Name text not null,
  Expiry DATE not null default current_timestamp,
  MaxUsers integer not null default 0,
);

CREATE TABLE Users (
  ID integer primary key autoincrement,
  CompanyID integer not null,
  Name text not null,
  Password text not null,
  AccessLevel int not null default 0,
  FOREIGN KEY (CompanyID) REFERENCES Company(ID)
);

COMMIT TRANSACTION;
