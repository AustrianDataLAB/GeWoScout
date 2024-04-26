# Terraform deployment
This is how you should perform the deployment of this architecture.

## Prerequisites
The terraform will build the applications and deploy them, that's why you also need to install those tools:
1. Azure `az cli`
2. Azure `func cli`
3. `npm` 
4. `go`

## Steps
1. Initialize the following environmental variables:
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

2. Perform `terraform init`:
```bash
terraform init -backend-config="resource_group_name=$TF_VAR_state_resource_group_name" -backend-config="storage_account_name=$TF_VAR_state_storage_account_name"
```

3. Perform `terraform apply`:
```bash
terraform apply -auto-approve
```

Don't forget to destroy the deployment after testing!
```bash
terraform destroy -auto-approve
```