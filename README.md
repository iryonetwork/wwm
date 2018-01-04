Iryo for Walk With Me
=====================

This project holds Iryo's first pilot project developed together with Walk With
Me foundation.

## Requirements

* docker
* docker-compose
* go (1.9+)
* govendor

## How to set up and work with the repository

```bash
# clone the repository
mkdir -p $GOPATH/src/github.com/iryo
cd $GOPATH/src/github.com/iryo
git clone git@gitlab.3fs.si:iryo/wwm.git
cd wwm

# prepare the repository (sync vendor folder)
make

# run tests
make test

# start everything up
make up
```

## Additional documentation

 * [WOW](docs/wow.md)
