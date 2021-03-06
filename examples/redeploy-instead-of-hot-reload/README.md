# Redeploy Instead of Hot Reload example

This example shows you how to develop a small go application with automated repeated deploying with devspace on devspace.cloud. For a more detailed documentation, take a look at https://devspace.cloud/docs

# Step 0: Prerequisites

## For local 
1. Install minikube (no docker required, since devspace uses the built in minikube docker daemon)

## Other
1. Make sure docker is installed
2. Replace `yourusername/devspace` with your desired image name in `devspace.yaml` and `kube/deployment.yaml`

# Step 1: Develop the application

1. Run `devspace dev` to start the application.

The command does several things in this order:
- Build the docker image
- Create needed image pull secrets
- Deploy the kubectl manifests (you can find and change them in `kube/`)
- Attach to the pod and print its output

You should see the following output:
```
[info]   Loaded config from devspace-configs.yaml
[info]   Using space redeploy                       
[info]   Skip building image 'default'         
[info]   Deploying devspace-default with kubectl
deployment.extensions/devspace unchanged           
[done] √ Finished deploying devspace-default       
[info]   The Space is now reachable via ingress on this URL: https://yourusername.devspace.host
[info]   Will now try to print the logs of a running devspace pod...
[info]   Printing logs of pod devspace-59c4d868f8-j2pg5/default...
Hello World!
Hello World!
Hello World!
Hello World!
```
2. Change the message in main.go and you should see the container reloading 
```
[info]   Change detected, will reload in 2 seconds
[info]   Building image 'dscr.io/yourusername/devspace' with engine 'docker'
[done] √ Authentication successful (dscr.io)
Sending build context to Docker daemon  12.29kB
Step 1/6 : FROM golang:1.11
 ---> 901414995ecd
[...]
[info]   Image pushed to registry (dscr.io)
[done] √ Done processing image 'dscr.io/yourusername/devspace'
[info]   Deploying devspace-default with kubectl
deployment.extensions/devspace configured          
[done] √ Finished deploying devspace-default       
[info]   The Space is now reachable via ingress on this URL: https://yourusername.devspace.host
[info]   Will now try to print the logs of a running devspace pod...
[info]   Printing logs of pod devspace-5d5cdb7554-lmxsg/default...
Hello World devspace!
```

# Troubleshooting 

If you experience problems during deploy or want to check if there are any issues within your deployed application devspace provides useful commands for you:
- `devspace analyze` analyzes the namespace and checks for warning events and failed pods / containers
- `devspace enter` open a terminal to a kubernetes pod (the same as running `kubectl exec ...`)
- `devspace logs` shows the logs of a devspace (the same as running `kubectl logs ...`)
- `devspace purge` delete the deployed application

See https://devspace.cloud/docs for more advanced documentation
