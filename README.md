# Alert Handler for Custom Metrics 


Alert Handler implements automated responses to telemetry-based alerts enabling the system to adapt to state change. 
It listens on port, waiting for incoming JSON packet describing alerts via webhook. On receipt of an alert it triggers an action or a user-configured script.

## Installation
You will need a go 1.11+ environment in order to build Alert Handler for Custom Metrics. Once you have a go environment, run 

``go build``

in the Alert Handler for Custom Metrics directory. The resulting executable, given default configuration, should be run in the same folder.
Note: For security reasons, Alert Handler for Custom Metrics should be owned and run by a non-priveledged Linux user.

## Configuration 

The file alert-handler-config.json contains the user configurable options. From this file you can add custom alerts, change the port number and URL path, and change the location of the scripts directory.
New configuration files are picked up at run time.
