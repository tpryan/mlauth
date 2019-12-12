# MLAuth
This is a collection of packages that use Google Cloud's ML APIs to act as a 
form of authentication. Basically, pass in a term or a condition and an ML-able 
artifact (iamge for vision, audio file for speech, etc..). Then the package 
uses the appropriate ML API to find the term in the artifact. 

## Requirements
* Create a Google Cloud Project
* Create env variables for project:
```` bash
export mlauth_project=[Name of proejct]
export mlauth_bucket=[Name of testing bucket] 
````
* Run init command
```` bash
make init
````
* Set the env variable for Google Cloud default credentials
```` bash
export GOOGLE_APPLICATION_CREDENTIALS=[Project Path]/creds/creds.json 
```` 

You can run `make test` afterwards to determine if everything worked. 


This will do the following:
* Create a bucket for testing
* Copy files for testing to bucket
* Create a service account for applications
* Create a key file so that we can use this service account in default 
credentials mode
* Activate Services in your project
    * Speech
    * Vision
    * Natural Language

This is not an official Google product. 