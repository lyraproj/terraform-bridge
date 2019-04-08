module github.com/lyraproj/terraform-bridge

go 1.12

require (
	cloud.google.com/go v0.36.0 // indirect
	github.com/Azure/azure-sdk-for-go v24.1.0+incompatible // indirect
	github.com/aws/aws-sdk-go v1.16.26 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/google/uuid v1.1.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/hashicorp/go-azure-helpers v0.0.0-20190129193224-166dfd221bb2 // indirect
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/terraform v0.11.11
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/lyraproj/issue v0.0.0-20190329160035-8bc10230f995
	github.com/lyraproj/pcore v0.0.0-20190408134742-7ef8f288585f
	github.com/lyraproj/servicesdk v0.0.0-20190408134916-985421696619
	github.com/marstr/guid v1.1.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/satori/uuid v1.2.0 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/terraform-providers/terraform-provider-aws v1.57.0
	github.com/terraform-providers/terraform-provider-azurerm v1.21.0
	github.com/terraform-providers/terraform-provider-github v1.3.0
	github.com/terraform-providers/terraform-provider-google v1.20.0
	github.com/terraform-providers/terraform-provider-kubernetes v1.5.0
	go.opencensus.io v0.19.0 // indirect
	golang.org/x/oauth2 v0.0.0-20190212230446-3e8b2be13635 // indirect
	k8s.io/apimachinery v0.0.0-20190126155707-0e6dcdd1b5ce // indirect
)

replace github.com/google/go-github => github.com/google/go-github v16.0.0+incompatible // Terraform GitHub provider requires this version
