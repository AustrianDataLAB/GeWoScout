# Terraform deployment
This is how you should perform the deployment of this architecture. Because of some limitations this can only work in an environment which supports bash-like variables (linux, macOS, WSL...).

## Prerequisites
The terraform will build the applications and deploy them, that's why you also need to install those tools:
1. Azure `az cli`
2. Azure `func cli`
3. `npm` 
4. `go`

## Steps
1. Build swagger
 - go to the `backend` dir
 - run: `go install github.com/swaggo/swag/cmd/swag@latest` and `swag init -g handler.go`

2. Initialize the following environmental variables:
 - ``TF_VAR_state_resource_group_name``
 - ``TF_VAR_state_storage_account_name``
 - ``TF_VAR_infra_resource_group_name``
 - ``TF_VAR_infra_subscription_id``
 - ``TF_VAR_infra_tenant_id``
 - ``ARM_CLIENT_ID``
 - ``ARM_CLIENT_SECRET``
 - ``ARM_SUBSCRIPTION_ID``
 - ``ARM_TENANT_ID``
 - ``TF_VAR_arm_client_id``
 - ``TF_VAR_arm_client_secret``
 - ``TF_VAR_arm_tenant_id``

3. Perform `terraform init`:
```bash
terraform init -backend-config="resource_group_name=$TF_VAR_state_resource_group_name" -backend-config="storage_account_name=$TF_VAR_state_storage_account_name"
```

4. Perform `terraform apply`:
```bash
terraform apply -auto-approve
```

Don't forget to destroy the deployment after testing!
```bash
terraform destroy -auto-approve
```

## SWA deployment (https://learn.microsoft.com/en-us/azure/static-web-apps/static-web-apps-cli-deploy)
- install SWA cli `npm install -g @azure/static-web-apps-cli`
- get the deployment token from Azure Portal (https://learn.microsoft.com/en-us/azure/static-web-apps/static-web-apps-cli-deploy#deployment-token)
- store the token in a variable (SWA_CLI_DEPLOYMENT_TOKEN)
- build the frontend (the output files should be in the `dist` directory)
- make sure that you are in <project-dir>/frontend/vue-gewoscout and run `swa deploy --app-location ./dist --env <environment-name>` (<environment-name> is `Development` or `Production`)

## Link the backend to a SWA environment
- deploy the backend (if you perform a terraform deployment `terrraform apply -auto-approve -target null_resource.backend_env`, you will get a file: `env.sh` file with the necessary backend ID and location)
- login to azure using your account `az login` or using a service principal:
```az login --service-principal -u ${arm_client_id} -p ${arm_client_secret} --tenant ${arm_tenant_id}```
- if linking the production environment:
    - unlink the old Production backend:
    ```az staticwebapp backends unlink --name gewofrontend --resource-group rg-management-brotholomew```
    - link a new Production backend:
    ```az staticwebapp backends link --backend-resource-id ${BACKEND_FUNCTION_ID} --name gewofrontend --resource-group rg-management-brotholomew --backend-region ${PROJECT_LOCATION}```
- if linking the development environment:
    - unlink the old Production backend:
    ```az staticwebapp backends unlink --name gewofrontend --resource-group rg-management-brotholomew --environment-name Development```
    - link a new Production backend:
    ```az staticwebapp backends link --backend-resource-id ${BACKEND_FUNCTION_ID} --name gewofrontend --resource-group rg-management-brotholomew --backend-region ${PROJECT_LOCATION} --environment-name Development```

please note that `${BACKEND_FUNCTION_ID}` is the fully qualified id of the resource, e.g.: `/subscriptions/e31c37ff-9b82-4f6b-8337-51314cc300ff/resourceGroups/rg-management-gewopro/providers/Microsoft.Web/sites/funcapp-backend-gewoscout-13uwpx` and the location will probably be: `westeurope`
