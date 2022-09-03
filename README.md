# Nats JWT Setup

How to quickly get a NATS JWT setup working.

```bash
NATS_SERVER_ADDR=localhost:4222
NATS_ENV_NAME=local
NATS_USER_ACCOUNT_NAME=APP
NATS_USER_NAME=user

source ./nsc.rc
nsc add operator --generate-signing-key --sys --name ${NATS_ENV_NAME}
nsc edit operator --require-signing-keys --service-url "nats://${NATS_SERVER_ADDR}" --account-jwt-server-url "nats://${NATS_SERVER_ADDR}"
nsc generate config --nats-resolver --sys-account SYS > ./conf/resolver.conf
```

Modify the `./conf/resolver.conf` to:

* set `resolver { dir = "./data/jwt" }` instead of `resover { dir = "./jwt" }`

Modify teh `./conf/nats.conf` to:

* append with `include resolver.conf`

Server can now be started with:

```bash
nats-server ./conf/nats.conf
```

You should no longer need to generate/change the operator/sys account unless a newer version of NATS requires the permissions to change.

## Setup an account

```bash
nsc add account ${NATS_USER_ACCOUNT_NAME}
nsc edit account ${NATS_USER_ACCOUNT_NAME} --sk generate
nsc add user --account ${NATS_USER_ACCOUNT_NAME} ${NATS_USER_NAME}

nsc push -a ${NATS_USER_ACCOUNT_NAME}
```

## Allow use of the users

```bash
nats context save "${NATS_ENV_NAME}-admin" --nsc "nsc://${NATS_ENV_NAME}/SYS/sys"
nats context save "${NATS_ENV_NAME}-${NATS_USER_NAME}" --nsc "nsc://${NATS_ENV_NAME}/${NATS_USER_ACCOUNT_NAME}/${NATS_USER_NAME}"
nats context select "${NATS_ENV_NAME}-admin"
nats server ls
```

## Use raw NKEY/JWT

Get the values from the `./keys/creds/${NATS_ENV_NAME}/${NATS_USER_ACCOUNT_NAME}/${NATS_USER_NAME}.creds`.

Open up `main.go` and 

* replace `jwt` string value with the JWT from the file.
* replace `seed` string value with the user nkey seed from the file.

Idealistically these two values would be pulled from vault, and not from a const.
