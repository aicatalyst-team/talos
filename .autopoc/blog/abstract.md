# Blog Abstract: Deploying Ory Talos on OpenShift

## Thesis
Ory Talos, an open-source API key management server, deploys cleanly on OpenShift with a UBI-based container image, validating that Go microservices with embedded databases are excellent candidates for platform security infrastructure.

## Target Audience
Platform engineers and DevOps teams managing AI infrastructure who need API key management for inference endpoints, model registries, and agent tool APIs.

## Blog Type
Red Hat Developer Blog

## Key Points (3 max)
1. Talos compiles to a single static binary with embedded SQLite, making containerization trivial on UBI
2. The full API key lifecycle (create, verify, revoke) works reliably on OpenShift with proper init containers for DB migration
3. Built-in Prometheus metrics and health endpoints make it Kubernetes-native from day one

## Products/Projects
- Red Hat OpenShift AI
- Open Data Hub
- Ory Talos

## CTA
Try deploying your own API key management layer on OpenShift using the AutoPoC pipeline.

## Proposed Sections
1. What is Ory Talos?
2. Why API key management matters for AI platforms
3. Containerizing a Go binary on UBI
4. Deploying to OpenShift with init containers
5. Testing the API key lifecycle
6. What we learned
7. Try it yourself
