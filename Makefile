GOCMD=go


test-apollo:
	export APOLLO_HOST=""
	export APOLLO_SCHEME="https"
	export APOLLO_PORT="443"
	export APOLLO_NAMESPACE="application"
	export APOLLO_CLUSTER="default"
	export APOLLO_USERNAME=""
	export APOLLO_PASSWORD=""
	$(GOCMD) test  ./richtesting/apollotesting/apollotesting_test.go