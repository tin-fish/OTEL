# localtestenv
This docker container setup builds the following:-
* Splunk OTEL Collector (requiring the following env variables)
  * SPLUNK_REALM (for dest. signalfx)
  * SPLUNK_ACCESS_TOKEN
* Node Exporter (for ^ to scrape)
* OTLPumper sending metrics

The configuration of the OTel collector and the narrative emitted by the OTLPumper can be overridden using the yaml files.
