Form3 Take Home Exercise
=====

### Hello there :wave:!
I just want to give a brief introduction and explain my background using Go.
My name is Renato, I am from Brazil, currently living in Berlin Germany and working for HelloFresh for about 1 year and 4 months.
My background is in PHP, 11 years working with it and only having a close contact with Go during my current stay in HelloFresh.
In HelloFresh we use mostly Go but also Python, Kotlin and a bit of PHP. So even 16 months of experience doesn't mean constantly
coding in Go.

But anyway, I hope that you enjoy the code.

I am looking forward for your feedback! 

Cheers :beers:

## Documentation

[![Test](https://github.com/renatoaraujo/form3-interview-accountapi/actions/workflows/test.yml/badge.svg)](https://github.com/renatoaraujo/form3-interview-accountapi/actions/workflows/test.yml)

In order to use this library in your project you just need to import it using the following command:
```bash
$ go get github.com/renatoaraujo/form3-account-api-client@1.0.0
```

#### Accounts

Import the module 

```go
import "renatoaraujo/form3-account-api-client/accounts"
```

To create, fetch or delete an account resource you need to initiate the client with the base uri and the request timeout

```go
httpClient, err := httputils.NewClient("https://api.form3.tech", 10)


accountClient := accounts.NewClient(httpClient)
```

And finally just call action

```go
// generates a valid accounts.AccountData{} 
accountData := &accounts.AccountData{}

// create resource sending the account data and it will return an accounts.AccountData{} or an error
created, err := accountClient.CreateResource(accountData)

// generates an uuid for the account id
accountID, _ := uuid.Parse("f199fe08-90b4-4756-9c1f-3a2352ea4933")

// fetch resource and it will return an accounts.AccountData{} or an error
fetched, err := accountClient.FetchResource(accountID)

// and finally delete a resource, and it will return an error or nil
err := accountClient.DeleteResource(accountID)

```

## Testing

To test the package you can just up the containers with the following command 
```bash
$ docker compose up
```

or you can use the make command available

```bash
$ make tests
```

## Final notes

### Project structure

I imagine this library would be specifically for [Organization](https://api-docs.form3.tech/api.html#organisation), 
so I chose to keep the project structure as a reflection of the documentation, as I believe it makes it easier for a 
developer to inspect if they wanted to know what's going on under the hood.

Based on that I would structure packages with names like:
- [Accounts](https://api-docs.form3.tech/api.html#organisation-accounts) (in the case implemented in this exercise)
- [Accounts Notifications](https://api-docs.form3.tech/api.html#organisation-account-identifications)
- [Account Events](https://api-docs.form3.tech/api.html#organisation-account-events)
- etc...

The packages are completely independent, they don't even have knowledge of a http client implementation, 
which I opted to leave it as external to simplify the solution and create cleaner, easier-to-change code and also to 
become easier to test.

You can find the http client implementation inside `httputils` package.

### Integration tests

Integration tests are simple, you can find it in the `/integration_tests` directory.
Considering this library being the main one for the Organization session, if in the future the Account Identifications 
context was inserted, this would be the place where the integration tests would be with a filename probably 
`account_identifications_test.go`.

The integration tests, obviously, depends on a functional API so checks if the API is available will happen and 
in case of unavailable host the tests will be skipped.

Once I realised the failure on GitHub Actions job I started a "fancy" solution for spinning up Docker on the fly using
[dockertest](https://github.com/ory/dockertest) library but in the end of the day I decide to just skip the tests case the host is 
unavailable, because, this is what I would expect from a real life scenario. We don't want to force this dependency for
services using this library.

### Continuous Integration

I decided to play around with GitHub Actions a bit so that I could provide a simple pipeline, but I just added one job 
for unit tests and setup the badge to add in this README file.

In a real life scenario I would probably add some support for the integration tests, together with an 
automatic semver release based in the release tags.

I will probably do it in the future for a learning experience, still in this repository.
