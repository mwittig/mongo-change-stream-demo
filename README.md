# mongo-change-stream-demo

A prototype to demonstrate a topic-based PUSH/PUSH Message Queue implementation based on a MongoDB collection. 

Instead of using findAndModify() at the subscriber end, the prototype makes use of Change Streams. A filter 
is applied to demonstrate the possibility of establishing subscription topics. This way, the number of insert events
propagated to a single subscriber can be reduced. 

Documents are removed by the Mongo cluster by establishing a TTL index.

Note: To be able to utilize Change Streams, a MongoDB Cluster (replica set or sharded cluster) needs to be setup.
Thus, a docker compose file is included which provided a local MongoDB Cluster setup.

## Prerequisites

* Docker setup on your local machine.
* You may need to install docker compose to run the cluster. It may also work OOTB using your favorite IDE.
* You need to add some entries in /etc/hosts as follows
````
  127.0.0.1       mongo1
  127.0.0.1       mongo2  
  127.0.0.1       mongo3
````
## References

* MongoDB Change Streams - https://www.mongodb.com/docs/manual/changeStreams/
* MongoDB FindAndModify Command - https://www.mongodb.com/docs/manual/reference/command/findAndModify/
* MongoDB TTL Index - https://www.mongodb.com/docs/v5.0/tutorial/expire-data/
* MongoDB Message Queue tutorial  - https://learnmongodbthehardway.com/schema/queues/




