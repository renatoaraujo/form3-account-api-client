Form3 Take Home Exercise
=====
----
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

@TODO

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

@TODO
