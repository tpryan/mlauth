# Copyright 2019 Google Inc. All Rights Reserved.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#      http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BASEDIR = $(shell pwd)
PROJECT = $(mlauth_project)
BUCKET = $(mlauth_bucket)

env:
	gcloud config set project $(PROJECT)

init: env bucket serviceaccount services key

serviceaccount:
	gcloud iam service-accounts create mlauth

key:
	gcloud iam service-accounts keys create creds/creds.json \
  --iam-account mlauth@$(PROJECT).iam.gserviceaccount.com

bucket:
	-gsutil mb gs://$(BUCKET)
	gsutil -m cp $(BASEDIR)/vision/testdata/* gs://$(BUCKET)/vision 
	gsutil -m cp $(BASEDIR)/speech/testdata/* gs://$(BUCKET)/speech 
	gsutil -m acl -r ch -u AllUsers:R gs://$(BUCKET)/

services:
	gcloud services enable vision.googleapis.com	
	gcloud services enable speech.googleapis.com
	gcloud services enable language.googleapis.com

clean: 
	-gsutil rm -rf gs://$(BUCKET)/*
	-gsutil rb gs://$(BUCKET)
	-gcloud iam service-accounts delete mlauth@$(PROJECT).iam.gserviceaccount.com -q

test:
	cd language && go test
	cd speech && go test
	cd vision && go test