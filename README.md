# Next Steps: Setup

> _**Quick Recap:** A reverse proxy server sits in-between the origin server and the network, receiving requests and sending responses on behalf of the origin server. It appears to the public as the origin server even though it's actually not, thus hiding intel on the origin server._

- Deploy this Go project in a server dedicated for caching (a server-side cache or reverse proxy cache).
- Start the server with the environment variable `ORIGIN_SERVER` set to your origin server's origin indicator (`protocol://host`)

  _Example:_ This will cache request's to developer.mozilla.org

  ```.env
  ORIGIN_SERVER_URL=https://developer.mozilla.org
  ```

The domain assigned to this proxy server is your public website domain.
