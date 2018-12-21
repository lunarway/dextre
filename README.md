# dextre

## Table of contents
- [What is dextre?](#what-is-dextre)
- [Why the name dextre?](#why-dextre)
- [Prerequisuites?](#prerequisuites)
- [Commands](#commands)

## What is dextre?
`dextre` is a CLI for automating and securing draining of kubernetes nodes and restarts of pods in a kubernetes cluster.

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
$ dextre drain --help
Usage:
  dextre drain [flags]

Flags:
      --grace-period duration   pod grace-period (default 30s)
  -h, --help                    help for drain
      --node string             The node that lander should drain in a safe manner (required)

Global Flags:
      --kubeconfig string   kubeconfig file
```
`dextre drain` will go through the following process:

* List the pods to be terminated, for you to insepct
* Cordon the node, so that other pods won't be scheduled on the node
* Evict one pod at a time and wait for a new one to be started (regular pods first, system pods last)
* (Optional) Terminate the instance in an aws scaling group for a new one to be spun up.

### restart
```
$ dextre restart --help
Usage:
  dextre restart [flags]

Flags:
      --grace-period duration   pod grace-period (default 10s)
  -h, --help                    help for restart
      --label string            The labels that should be restarted on the form: type=service
      --namespace string        The namespace to search for pods

Global Flags:
      --kubeconfig string   kubeconfig file
```
`dextre restart` will restart pods by a given label and namespace. 
