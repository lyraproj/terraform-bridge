<p align="center"><img src="docs/media/lyrabanner.png" alt="Lyra"></p>

## What is the Lyra Terraform Bridge?

[Lyra](https://github.com/lyraproj/lyra) is an open source workflow engine for provisioning and managing cloud native infrastructure. Using infrastructure as code, Lyra enables you to declaratively provision and manage public cloud, private cloud, and other API-backed resources as well as orchestrate imperative actions.

The Lyra Terraform Bridge makes [Terraform providers](https://github.com/terraform-providers) useable in Lyra.

## Getting Started

If you want to make use of Terraform providers that have already been integrated you don't need to use this repo directly. When you build Lyra it will incorporate all this content automatically. If you wish to contribute by integrating new providers or improving the bridge, then you can do the following:

### Build
The project requires [Go](https://golang.org/doc/install) 1.11 or higher, and [go modules](https://github.com/golang/go/wiki/Modules) to be on.

Build the project using make:

	make

When no targets are specified, the build will lint, test, compile and sanity-check every plugin.

### Refreshing content

Refresh all currently integrated providers like this:

	make generate

### Integrating new content

To follow!

## Contributing
We'd love to get contributions from you! For a quick guide, take a look at our guide to [contributing](CONTRIBUTING.md).
