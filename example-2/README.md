# kafka-tutorial

## Example 2: Persistent Storage and Confluent's Apache Kafka Golang Client

This example builds on the foundations of Example 1 by introducing **persistent storage** for Kafka and using **two Go microservices** to interact with the cluster. Instead of just a CLI-based demo, we create a more realistic environment where:

- Kafka's controllers and brokers use **PersistentVolumeClaims (PVCs)**.
- Go clients use **confluent-kafka-go** to produce and consume messages.
- Topics are created via **Kubernetes Custom Resources**, not manually.
- The setup is fully compatible with both local kind clusters and real-world Kubernetes deployments (e.g., EKS).

### 🧱 Architecture

- 3 Kafka controller pods
- 3 Kafka broker pods (each with persistent storage)
- 2 demo apps (producer + consumer)
- 1 Kafka topic managed through Strimzi CRD

### Prerequisites 
- docker
- kubectl
- kind
- Go >= 1.22

## Getting Started

1. Create cluster
   ```
   cd example-2
   ```

   ```
   $ kind create cluster --config ../kind-cluster.yaml

   Creating cluster "kind" ...
   ✓ Ensuring node image (kindest/node:v1.34.0) 🖼
   ✓ Preparing nodes 📦 📦 📦 📦 📦 📦 📦 📦 📦
   ✓ Configuring the external load balancer ⚖️
   ✓ Writing configuration 📜
   ✓ Starting control-plane 🕹️
   ✓ Installing CNI 🔌
   ✓ Installing StorageClass 💾
   ✓ Joining more control-plane nodes 🎮
   ✓ Joining worker nodes 🚜
   Set kubectl context to "kind-kind"
   You can now use your cluster with:

   $ kind get clusters
   kind

   $ kubectl config get-contexts                                                                   
   CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
            docker-desktop   docker-desktop   docker-desktop   
   *         kind-kind        kind-kind        kind-kind   

   $ kubectl cluster-info --context kind-kind

   Kubernetes control plane is running at https://127.0.0.1:59570
   CoreDNS is running at https://127.0.0.1:59570/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

   To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

   Have a nice day! 👋

   $ kubectl get nodes           
      NAME                  STATUS   ROLES           AGE   VERSION
      kind-control-plane    Ready    control-plane   48m   v1.34.0
      kind-control-plane2   Ready    control-plane   47m   v1.34.0
      kind-control-plane3   Ready    control-plane   47m   v1.34.0
      kind-worker           Ready    <none>          46m   v1.34.0
      kind-worker2          Ready    <none>          46m   v1.34.0
      kind-worker3          Ready    <none>          46m   v1.34.0
      kind-worker4          Ready    <none>          46m   v1.34.0
      kind-worker5          Ready    <none>          46m   v1.34.0
      kind-worker6          Ready    <none>          46m   v1.34.0
    ```

2. Create namespace
   
   ```
    $ kubectl create -f 0-namespace.yaml
   ```

3. Label worker nodes dedicated for Kafka workloads

   ```
   $ kubectl label node kind-worker kafka-nodes=true
   $ kubectl label node kind-worker2 kafka-nodes=true
   $ kubectl label node kind-worker3 kafka-nodes=true
   $ kubectl get nodes -l kafka-nodes=true
   ```

4. Create Strimzi Cluster Operator
   > Note: The operator manifest has been modified with a nodeSelector so that it also runs on the Kafka nodes. 
   ```
   $ kubectl create -f 1-strimzi-cluster-operator-0.48.0.yaml
   $ kubectl get pod -n kafka --watch

   NAME                                        READY   STATUS              RESTARTS   AGE
   strimzi-cluster-operator-64574988c8-6xhz9   0/1     ContainerCreating   0          9s
   strimzi-cluster-operator-64574988c8-6xhz9   0/1     Running             0          11s
   strimzi-cluster-operator-64574988c8-6xhz9   1/1     Running             0          43s
   ```
   
5. Deploy Kafka with Persistent Volumes
    ```
    kubectl apply -f 2-kafka-persistent.yaml -n kafka
    kubectl get pods -n kafka --watch -o wide
    kubectl get pvc -n kafka
    ```
6. Create Kafka Topic via CRD
    ```
    kubectl apply -f 3-topic.yaml -n kafka
    kubectl get kafkatopics -n kafka
    ```
7. Build & Load Go Services
   ```
   docker build -t kafka-producer-go:latest ./producer-service
   docker build -t kafka-consumer-go:latest ./consumer-service
   kind load docker-image kafka-producer-go:latest
   kind load docker-image kafka-consumer-go:latest
   ```
   
8. Deploy Go Producer & Consumer
   ```
   kubectl apply -f 4-producer.yaml
   kubectl apply -f 5-consumer.yaml
   ```

9.  Verify Message Flow
      ```
      kubectl logs -n kafka deploy/go-kafka-consumer -f
      ```  
10.  Clean up
      ```
      kubectl delete ns kafka
      ```