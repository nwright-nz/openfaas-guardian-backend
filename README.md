# Guardian Backend for OpenFaaS

This is a very basic proof of concept for adding a backend to the OpenFaaS project (https://github.com/alexellis/faas).  
  
Please note: This is very basic and obviously shouldnt be run anywhere NEAR production :)  

The aim of this project is to extend the OpenFaaS project to be able to create 'serverless' functions using Garden and RunC, with the idea that this will be extending to running on Cloud Foundry.
 
At the moment, the following operations are supported:
* Creating new functions
* Invoking functions
* tracking invocation count
* Deleting functions

The following has not been implemented yet:

* Prometheus alerting
* Autoscaling containers based on alerts


Once all the work has been completed, the aim is that this will be able to be deployed to Cloud Foundry.  

## Install Instructions

### Install Bosh Lite
The instructions to install Bosh Lite for your environment are available here: https://bosh.io/docs/bosh-lite.html

Once this has been installed and configured, move onto the next step.

### Deploy Garden-RunC release
The Garden-RunC repository has the required instructions (https://github.com/cloudfoundry/garden-runc-release)

In summary:

Clone it:
```
git clone https://github.com/cloudfoundry/garden-runc-release

cd garden-runc-release

git submodule update --init --recursive
```

Update the cloud config:
```
bosh -e vbox update-cloud-config ./manifests/cloud-config-lite.yml
```

Deploy the release:
```
./scripts/deploy-lite.sh
```

Once this has deployed you can get the garden-runc server with the following command (assuming you used the virtual box CPI, otherwise substitute your environment):
```
bosh -e vbox vms
```

If you use the cloud-config in the garden-runc repo, your server will probably have the ip : 10.244.0.2. The port will default to 7777.

To test all is ok, install the gaol command line utility (https://github.com/contraband/gaol) and create a basic container:

```
gaol -t <<server_name>>:7777 create -n my-container
```

## Deploy OpenFaaS backend

run the following command to get the repository:
```
go get github.com/nwright-nz/openfaas-guardian-backend
```

Run the gateway:

```
cd $GOPATH/src/github.com/nwright-nz/openfaas-guardian-backend
go run server.go 
```

Access localhost:8080 to see the OpenFaaS UI.

Please let me know if there are any issues with this, I am quite fresh to go, and things are still quite messy.


