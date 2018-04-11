#### Setup cloud instance

This assumes the database schema is already in place.

```
docker-compose run  cloudSymmetric ./symadmin --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 create-sym-tables
docker-compose run  cloudSymmetric ./dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 ../samples/initial.sql
docker-compose run  cloudSymmetric ./sym --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326
```

Database schema can be initialised / upgraded with:

```
docker-compose run  cloudSymmetric ./dbimport --engine cloud-f7e41e48-ec79-4c78-9db6-37c0c4f78326 --format XML --alter-case ../samples/schema.xml
```

#### Setup local instance

```
docker-compose run  localSymmetric ./sym --engine local-2d04b22e-1cc3-46b4-96dd-2bee5bad9ffa
```
