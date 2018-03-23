#### Setup cloud instance

```
docker-compose run  cloudSymmetric ./dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 --format XML --alter-case ../samples/create-sample.xml
docker-compose run  cloudSymmetric ./symadmin --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 create-sym-tables
docker-compose run  cloudSymmetric ./dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 ../samples/initial.sql
docker-compose run  cloudSymmetric ./sym --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 --debug
```

#### Setup

#### Setup local instance

```
docker-compose run  localSymmetric ./sym --engine local-2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa --debug
```
