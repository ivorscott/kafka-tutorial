# kafka-tutorial

1. Create cluster

   ```
   kind create cluster --config kind-cluster.yaml
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

   kubectl cluster-info --context kind-kind

   Have a nice day! 👋
    k get nodes           
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
    kubectl create -f manifests/0-namespace.yaml
    namespace/kafka created
   ```

3. Label worker nodes dedicated for Kafka workloads

   ```
   kubectl label node kind-worker kafka-nodes=true
   kubectl label node kind-worker2 kafka-nodes=true
   kubectl label node kind-worker3 kafka-nodes=true
   node/kind-worker labeled
   node/kind-worker2 labeled
   node/kind-worker3 labeled
   ```

4. Create Strimzi Cluster Operator
   > Note: To ensure the operator pod runs on a dedicated kafka worker node, the Deployment resource was modified with: 
   > 
   >  ``` 
   >   nodeSelector:
   >     kafka-nodes: "true"
   >  ```
   ```
   kubectl create -f manifests/1-strimzi-cluster-operator-0.48.0.yaml
   rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-entity-operator-delegation created
   rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-leader-election created
   serviceaccount/strimzi-cluster-operator created
   customresourcedefinition.apiextensions.k8s.io/kafkarebalances.kafka.strimzi.io created
   clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-leader-election created
   customresourcedefinition.apiextensions.k8s.io/kafkatopics.kafka.strimzi.io created
   customresourcedefinition.apiextensions.k8s.io/kafkaconnects.kafka.strimzi.io created
   clusterrole.rbac.authorization.k8s.io/strimzi-kafka-broker created
   customresourcedefinition.apiextensions.k8s.io/kafkamirrormaker2s.kafka.strimzi.io created
   clusterrolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-kafka-client-delegation created
   customresourcedefinition.apiextensions.k8s.io/kafkabridges.kafka.strimzi.io created
   clusterrolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-kafka-broker-delegation created
   rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator created
   clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-watched created
   clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-global created
   customresourcedefinition.apiextensions.k8s.io/kafkanodepools.kafka.strimzi.io created
   configmap/strimzi-cluster-operator created
   clusterrole.rbac.authorization.k8s.io/strimzi-entity-operator created
   customresourcedefinition.apiextensions.k8s.io/kafkas.kafka.strimzi.io created
   deployment.apps/strimzi-cluster-operator created
   clusterrole.rbac.authorization.k8s.io/strimzi-kafka-client created
   rolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator-watched created
   clusterrole.rbac.authorization.k8s.io/strimzi-cluster-operator-namespaced created
   customresourcedefinition.apiextensions.k8s.io/kafkaconnectors.kafka.strimzi.io created
   customresourcedefinition.apiextensions.k8s.io/kafkausers.kafka.strimzi.io created
   customresourcedefinition.apiextensions.k8s.io/strimzipodsets.core.strimzi.io created
   clusterrolebinding.rbac.authorization.k8s.io/strimzi-cluster-operator created

   kubectl get pod -n kafka --watch
   NAME                                        READY   STATUS              RESTARTS   AGE
   strimzi-cluster-operator-64574988c8-6xhz9   0/1     ContainerCreating   0          9s
   strimzi-cluster-operator-64574988c8-6xhz9   0/1     Running             0          11s
   strimzi-cluster-operator-64574988c8-6xhz9   1/1     Running             0          43s
   ```

5. Deploy Kafka workloads using ephemeral storage example. 

   > Note: This manifest defines two KafkaNodePools (3 controllers and 3 brokers) using ephemeral storage and pins them to the labeled Kafka worker nodes.

   ```
   kubectl apply -f manifests/2-kafka-ephemeral.yaml -n kafka
   kafkanodepool.kafka.strimzi.io/controller created
   kafkanodepool.kafka.strimzi.io/broker created
   kafka.kafka.strimzi.io/kafka-cluster unchanged
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