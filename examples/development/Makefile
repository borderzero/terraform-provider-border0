.PHONY: init token apply destroy

env ?= dev
api ?= http://localhost:8080/api/v1
portal ?= http://localhost:8081
tunnel ?=
capi ?=
token ?= $$(cat ~/.border0/token)

ifeq ($(env), dev)
	api = http://localhost:8080/api/v1
	portal = http://localhost:8081
else ifeq ($(env), staging)
	api = https://api.staging.border0.com/api/v1
	portal = https://portal.staging.border0.com
	tunnel = tunnel.staging.border0.com
	capi = capi.staging.border0.com:443
else ifeq ($(env), prod)
	api = https://api.border0.com/api/v1
	portal = https://portal.border0.com
	tunnel = tunnel.border0.com
	capi = capi.border0.com:443
else
	$(error env must be dev, staging or prod)
endif

init:
	(cd ../../; make install)
	rm -rf .terraform*
	terraform init

token:
	BORDER0_API=$(api) BORDER0_WEB_URL=$(portal) border0 login

apply:
	@TF_VAR_token=$(token) TF_VAR_api_url=$(api) terraform apply

destroy:
	@TF_VAR_token=$(token) TF_VAR_api_url=$(api) terraform destroy

connector:
	@BORDER0_API=$(api) \
		BORDER0_WEB_URL=$(portal) \
		BORDER0_TUNNEL=$(tunnel) \
		BORDER0_CONNECTOR_SERVER=$(capi) \
		BORDER0_TOKEN=$(token) \
		border0 connector start --v2
