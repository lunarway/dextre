# dextre

## Table of contents
- [What is dextre?](#what-is-dextre)
- [Why did we build dextre?](#why-did-we-build-dextre)
- [Why the name dextre?](#why-dextre)
- [Prerequisuites?](#prerequisuites)
- [Commands](#commands)

## What is dextre?
`dextre` is a CLI for automating and securing draining of kubernetes nodes and restarts of pods in a kubernetes cluster. `dextre` is built for automating some processes around managing a kops kubernetes cluster, and currently only supporting AWS.

## Why did we build dextre?
When running a Kops kubernetes cluster, you can use the upgrading procedures of Kops to upgrade and roll your cluster. However, there are a number of things that this procedure is not taking care of for you, which potentially can result in downtime of your services if you are not aware. The rolling-update feature uses the kubernetes drain method for gracefully evicting pods on a node that needs to be rolled, however, this often result in many pods being rescheduled at the same time, leading to many pods fighting for bandwidth on the other nodes. To cope with this, we have manually been moving pods around during upgrades, however, this is no longer feasible, hence this tool. 

## Why the name dextre?
Definition of dextre by WikiPedia
>Dextre, also known as the Special Purpose Dexterous Manipulator (SPDM), is a two armed robot, or telemanipulator, which is part of the Mobile Servicing System on the International Space Station (ISS), and does repairs otherwise requiring spacewalks.
This description fits well with the purpose of `dextre` which is to automate tasks for your kubernetes and aws environments.

## Prerequisuites
In order to use `dextre`, you need:
* AWS credentials configured
* Kubeconfig configured

## Commands
`dextre` currently provides the following commands:

### drain

```
$ dextre drain --node <node_name>
```
`dextre drain` will go through the following process:
* List the pods to be terminated, for you to insepct
* Cordon the node, so that other pods won't be scheduled on the node
* Evict one pod at a time and wait for a new one to be started (regular pods first, system pods last)
* (Optional) Terminate the instance in an aws scaling group for a new one to be spun up.

### roll 

#### roll all pods with a given label
```
$ dextre roll pods --label app=service --namespace default
```

#### roll all nodes with a given role and/or label
```
$ dextre roll nodes --role node --label type=some
```
This command will run the drain command for each node and wait for a new node to be ready before continuing. 
