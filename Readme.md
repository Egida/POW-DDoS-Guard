# Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge

# Usage

To demonstrate the work of application use:
```shell
> make run
```

This script does following:
* Creates docker network
* Builds server and client images
* Runs both services. Client will print the info about number of POW iterations and the received quote and will stop executing after
* Stops server
* Cleans up docker images and network