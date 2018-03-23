# Symmetric replication

## About the tool

> SymmetricDS is open source database replication software that focuses on features and cross platform compatibility. (taken from https://www.symmetricds.org/)

Symmetric an utility that is able to replicate data modifications between multiple nodes. Any of those nodes can act as a master node. Instead of a proprietary replication protoc limited to a single database Symmetric can sync data between wide variety of databases. Additionally we can limit scope of replicated data with custom rules to limit which data should be replicated to specific nodes.

Internally Symmetric is a java based toolset that observes a specific tables inside a database and periodically sends data to defined set of nodes and also receives data from external nodes.

## Our setup

Two different groups set up inside symmetric:

* **cloud**
    Receives all data and syncs data to other nodes. Represents the database running in the cloud.
* **local**
    Contains data for patients linked to current location. Represents all external nodes.

Inside those groups we have nodes that are represented with currently manually assigned unique identifiers (`uuidgen`). To each of those nodes we need to assign `group`, `engine.name ($GROUP-$UUID)` and `external.id ($UUID)`. Each of the nodes can connect to a different database with their own custom database connection details. Same IDs are used inside the application to distinguish between different locations. IDs used in the development environment are listed in the main [README](../README.md).

To sync our data we require three different replication scenarios:

* **local_2_cloud**
  Copies all changes on a `local` node to the cloud node.
* **cloud_2_select_local**
  Tables `patients`, `connections` and `locations` all contain a `patient_id` column. For each row symmetric queries the `locations` table to check if patient is linked to a specific location to decide if data should be replication to that node.
* **cloud_2_local**
  All local nodes receive all data. Used for replicating `locations` table to all nodes.

Each of those nodes are defined inside the `./services/symmetric/engines` folder.

### Symmetric elements

#### Groups

## Operations

Symmetric requires some manual operations to enable all required operations. All files required to setup the environment are located inside `./services/symmetric`.

#### Initial "master" setup

```bash
# Import the initial database structure
# (reads the XML document and converts it to database specific SQL
# statements)
docker-compose run cloudSymmetric ./dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 --format XML --alter-case ../samples/create-sample.xml

# Create symmetric tables
# (We need to manually create them to be able to insert our own rules in
# the next step even though symmetric would create them automatically on
# first run)
docker-compose run cloudSymmetric ./symadmin --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 create-sym-tables

# Insert replication rules
# (Mind that this document also clears all preexisting configuration)
docker-compose run cloudSymmetric ./dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 ../samples/initial.sql

# Start the cloud node
make up/cloudSymmetric
```

#### Local nodes

```bash
# Register the local node on the cloud node
docker-compose run  cloudSymmetric ./symadmin open-registration --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 local 2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa

# Trigger the full reload of the local node
docker-compose run cloudSymmetric ./symadmin reload-node --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa

# Start the node
make up/localSymmetric
```

