.PHONY: init token apply destroy

api ?= http://localhost:8080/api/v1
portal ?= http://localhost:8081
token ?= $$(cat ~/.border0/token)

init:
	(cd ../../; make install)
	rm -rf .terraform*
	terraform init

token:
	BORDER0_API=$(api) BORDER0_WEB_URL=$(portal) border0 login

apply:
	TF_VAR_token=$(token) TF_VAR_api_url=$(api) terraform apply

destroy:
	TF_VAR_token=$(token) TF_VAR_api_url=$(api) terraform destroy
