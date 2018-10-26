# SMACC backend

My experience in Go: 2-3 years.

Technically, every single line of code (apart from vendor/) was written by myself.

# Challenge

API that is capable of sending emails using specified providers. The failover functionality is a requirement.
I am not sure how advanced should be the failover part. I will code the easiest solution.

## Aspects

- Security: are there any obvious vulnerability?
    - **There is no way to validate the 'from' field.** Therefore anyone may send mails from you or me.
    - At this point, I keep API keys at the config. I intend to push the code to private repository, so it won't be public. Let's assume the secret keys are just the 'dev' ones. The production keys are supposed to be set via environment variables for the container. (Since this repository is public, I edited the git commits as to remove secret credentials from config.go)
- Does the README contain information on how to run it: `CGO_ENABLED=0 go build -o server && ./server`

### Bonus point:

- Scalability: **Yes.** Application is stateless and should scale well.
- Production-readiness: I tried to create the production-ready application. The logging and metrics are covered (somehow). Building a microcontainer is also added (`docker build ./ -t smacc-backend`).

## Documentation

At this moment there is no documentation of the REST API, but I think it's *nice to have*.
In case I will have more free time.. maybe I will add it.
