# `examples/development`

This example is for developing and testing this terraform provider. It uses this provider from a local path
that was built and configured by `make install`.

```shell
# (re)build this terraform provider, install it locally with `terraform init`
make init

# use `border0 login` to (re)generate a user access token from local api
make token

# run `terraform apply` against border0 api that's running locally
make apply

# run `terraform destroy` with locally running api
make destroy
```
