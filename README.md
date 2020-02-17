# Jefe - A place to claim development environments

## Setup

### Requirements

- `mysql` DB
- [Github OAuth Application](https://github.com/settings/applications/new)


### Setup the DB

Once you have a `MySQL` instance, you can use the `setupdb.sql` script from the `/scripts/sql/` directory to setup the database.

You will have to configure Jefe with the `DSN` in order to be able to connect to your `MySQL` instance:

```
JEFE_DSN=user:pass@tcp(host:port)
```

### Setup an OAuth application

You'll need to create an [OAuth application on GitHub](https://github.com/settings/applications/new).

The "Authorization callback URL" must be the URL of your Jefe instance with `/oauth/redirect` appended:

`http://<my.jefe.com>/oauth/redirect`

You will be given a `Client ID` and a `Client Secret` for your new oauth application. The `client ID` and `secret` must then be configured on Jefe by setting the following env:

```
JEFE_GITHUB_CLIENT_ID: theClientID
JEFE_GITHUB_SECRET: theSecret 
```

Next you will need to provide a github organization Jefe will authenticate users against. That way you make sure that only users part of the given organization are authorized to use your Jefe instance.

```
JEFE_O_AUTH_ORG=my-github-org
```

### Create a session secret

Create a  32bit long session key which will be used encrypt and authenticate session data:

```
JEFE_SESSION_SECRET: mysecret 
```

## Deploy

You can deploy Jefe as a CloudFoundry application or as a container on any other Platform where you can deploy Golang apps or Containers.

### The CF way

1. `$ cf push jefe --no-start`
2. `$ cf set-env ENV_VAR <env-value>`
3. `$ cf set-env GO_INSTALL_PACKAGE_SPEC github.com/herrjulz/jefe/cmd/web`
3. Do this for all required environment variables
4. `$cf start jefe`

### Kubernetes Way

Deploy Eirini on Kubernetes and use the CF way ;)
