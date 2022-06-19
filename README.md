# SMACC backend

Experience in Go: 2-3 years.

Author implemented all code in this repository apart from code placed in vendor/.

# Challenge

API capable of sending emails using online mail providers. One of the requirements is failover.
Interviewer didn't provide detailed requirements for the failover, therefore solution contains the easiest solution --
iteration over providers until one of them succeeds.

## Aspects

- Security consierations:
    - **There is no way to validate the 'from' field.** Recipient isn't verified and could represent anyone, even CEO of a FANG company.
    - Configuration provides API keys, but there isn't any secure way of deployment.
- Does the README contain information on how to run it: `CGO_ENABLED=0 go build -o server && ./server`


### Bonus point:

- Scalability: **Yes.** Application is stateless and should scale well.
- Production-readiness: best-offert to achieve, for instance logging and metrics cover the observability part.
    Containerization of the microcontainer is achievable by the command: `docker build ./ -t smacc-backend`.

## Documentation

At this moment there is no documentation of the REST API.
Consider this point as *nice to have*.

