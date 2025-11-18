# kafka-tutorial

## Example 1: Ephemeral Storage and Kafka CLI Tools

This example spins up a Strimzi-managed Kafka cluster on kind using **ephemeral storage**, then use the built-in **Kafka CLI tools** (`kafka-console-producer` / `kafka-console-consumer`) to send and receive messages.

### Prerequisites 
- docker
- kubectl
- kind

## Getting Started



1. Create cluster
   ```
   cd example-1
   ```

   ```
   kind create cluster --config ../kind-cluster.yaml

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
   ```
   ```
   kind get clusters
   kind
   ```
   ```
   kubectl config get-contexts                                                                   
   CURRENT   NAME             CLUSTER          AUTHINFO         NAMESPACE
            docker-desktop   docker-desktop   docker-desktop   
   *         kind-kind        kind-kind        kind-kind   
   ```

   ```
   kubectl cluster-info --context kind-kind

   Kubernetes control plane is running at https://127.0.0.1:59570
   CoreDNS is running at https://127.0.0.1:59570/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

   To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

   Have a nice day! 👋
   ```

   ```
   kubectl get nodes           
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
    kubectl create -f 0-namespace.yaml
    namespace/kafka created
   ```

3. Label worker nodes dedicated for Kafka workloads

   ```
   kubectl label node kind-worker kafka-nodes=true
   kubectl label node kind-worker2 kafka-nodes=true
   kubectl label node kind-worker3 kafka-nodes=true
   kubectl get nodes -l kafka-nodes=true 
   ```

4. Create Strimzi Cluster Operator
   > Note: To ensure the operator pod runs on a dedicated kafka worker node, the Deployment resource was modified with: 
   > 
   >  ``` 
   >   nodeSelector:
   >     kafka-nodes: "true"
   >  ```
   ```
   kubectl create -f 1-strimzi-cluster-operator-0.48.0.yaml
   kubectl get pod -n kafka --watch

   NAME                                        READY   STATUS              RESTARTS   AGE
   strimzi-cluster-operator-64574988c8-6xhz9   0/1     ContainerCreating   0          9s
   strimzi-cluster-operator-64574988c8-6xhz9   0/1     Running             0          11s
   strimzi-cluster-operator-64574988c8-6xhz9   1/1     Running             0          43s
   ```

5. Deploy Kafka workloads using ephemeral storage example. 

   > Note: This manifest defines two KafkaNodePools (3 controllers and 3 brokers) using ephemeral storage and pins them to the labeled Kafka worker nodes.

   ```
   kubectl apply -f 2-kafka-ephemeral.yaml -n kafka
   ```

6. View pods

   ```
   kubectl get pod -n kafka -o wide

   NAME                                        READY   STATUS    RESTARTS   AGE     IP           NODE           NOMINATED NODE   READINESS GATES
   kafka-cluster-broker-0                      1/1     Running   0          26s     10.244.8.5   kind-worker2   <none>           <none>
   kafka-cluster-broker-1                      1/1     Running   0          26s     10.244.3.4   kind-worker3   <none>           <none>
   kafka-cluster-broker-2                      1/1     Running   0          26s     10.244.9.7   kind-worker    <none>           <none>
   kafka-cluster-controller-3                  1/1     Running   0          26s     10.244.3.5   kind-worker3   <none>           <none>
   kafka-cluster-controller-4                  1/1     Running   0          26s     10.244.8.6   kind-worker2   <none>           <none>
   kafka-cluster-controller-5                  1/1     Running   0          26s     10.244.9.8   kind-worker    <none>           <none>
   strimzi-cluster-operator-579d777887-zj7zb   1/1     Running   0          3m22s   10.244.9.6   kind-worker    <none>           <none>
   ```

7. Create a Producer and send messages
   ```
   kubectl -n kafka run kafka-producer -ti --image=quay.io/strimzi/kafka:0.48.0-kafka-4.1.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server kafka-cluster-kafka-bootstrap:9092 --topic my-topic
   >Hello Strimzi!
   ```
   
8. Create a Consumer and receive a messages
   ```
   kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.48.0-kafka-4.1.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server kafka-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning
   Hello Strimzi!
   ```

9. Delete Strimzi custom resources
   ```
   kubectl -n kafka delete $(kubectl get strimzi -o name -n kafka)
   ```

10. Delete the operator
      ```
      kubectl -n kafka delete -f 1-strimzi-cluster-operator-0.48.0.yaml
      ```

11. Delete the namespace
      ```
      kubectl delete namespace kafka
      ```

## References
- https://kind.sigs.k8s.io/docs/user/quick-start/
- https://strimzi.io/quickstarts/

