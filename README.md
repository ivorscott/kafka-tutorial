# kafka-tutorial

## Prerequisites 
- docker
- kubectl
- kind


1. [Ephemeral Storage and Kafka CLI Tools](/example-1/README.md)

   Spin up a Strimzi-managed Kafka cluster on kind using **ephemeral storage**, then use the built-in **Kafka CLI tools** (`kafka-console-producer` / `kafka-console-consumer`) to send and receive messages.  
      - Good for: first contact with Strimzi, HA layout, controllers vs brokers.

2. [Persistent Storage and Confluent's Apache Kafka Golang Client](/example-2/README.md)

   Run a Kafka cluster with **persistent volumes**, create Kafka topics via **CRDs**, and use two **Go microservices** (producer + consumer) built with `confluent-kafka-go` to exchange messages.  
      - Good for: realistic setup, PVCs, app integration, foundation for security/HPA/observability.
