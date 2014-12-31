MCL Vehicle Tracker/Logger
=======================

myClubLink Vehicle tracking and logging project.

This project was created to track various vehicles. It can record vehicle positions along with speed , 
temperature + Arm SOC details and other metrics in real time. 

All metrics are sent to a backend server for logging via 3G (TCP) 
with radio to come (for remote locations where there is no 3G.

Clients of the system can login via a web page and observe the vehicle locations in real time via their browser. 
Clients can also obtain reporting information about the vehiles and other historical informations such as 
a particular route that was taken by a vehicle at a particular time and show second by second the metrics that were
pertainable.

Goals
-----
* To create something useful that is fun.
* To learn new technology. (Was new to AngularJS, Go , Websockets and I had to learn the Phidget SDK)
* To improve my current skills - My HTML5/CSS# skills need a lot of work (I'm not a designer) , I need to add more security. 
* To have a lot of fun !
* To monetize it if I finish to support my family ( Currently had hardware problems (overheating etc) delayed for now
* To create something that is extensible - multiple map vendor, different DBMS support  (currently only SQL Lite)

Constraints
-----------
This was done in my spare time. I have a full time job , a wife and 3 kids.
I was learning on the go. Expect bugs etc. Here be dragons !

Features
---------
* Real time tracking of vehicles.
* Multiple vehicle real time track and record trip.
* Basic reports - total KM travelled.
* Unlimited client login and observe  (according to licence)
* Option to support multiple mapping vendor (Google and some Leaflet done)

Tech Stack
----------
Linux
Nginx
Sql Lite
Javascript
Go (Backend)
AngularJS
Web Browsers (Anthing modern supporting CSS3 , HTML5 and web sockets)

Components
------------
The backend/server
********************
https://github.com/bradmccormack/MCLVehicleTrackerLogger/tree/master/backend/src/main

The client side code
*********************
https://github.com/bradmccormack/MCLVehicleTrackerLogger/tree/master/frontend

An example vehicle logger
***************************
https://github.com/bradmccormack/MCLVehicleTrackerLogger/tree/master/loggers


Utilities
*********

Replay
________
A replay facility to "replay" previous real trips back to the server in test mode for replicating and observing
previous routes. multiplevehicles.sh is provided as an example of showing replay in use for simulating real time 
updates of multiple vehicles.

https://github.com/bradmccormack/MCLVehicleTrackerLogger/tree/master/backend/src/main/misc

Schema creator
__________________
https://github.com/bradmccormack/MCLVehicleTrackerLogger/blob/master/backend/src/main/createdb.go
Will create the various databases for you.

OSM Importer
___________
https://github.com/bradmccormack/MCLVehicleTrackerLogger/blob/master/backend/src/main/osmreader.go

Will import OSM files from Open Street Map into the schema as dictated by createdb.go (please adapt)
Example location provided https://github.com/bradmccormack/MCLVehicleTrackerLogger/blob/master/backend/src/main/nowrasmall.osm


Screenshots to come!




