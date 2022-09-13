How long did this assignment take?
> 4h for the logic implementation + 1h for tests + 1h for cleanup

What was the hardest part?
> not to dive to deeply into details

Did you learn anything new?
> can't say so, but anyway it was a good peace of practice

Is there anything you would have liked to implement but didn't have the time to?
> unit tests, ci pipeline, swagger, improve auth, caching

What are the security holes (if any) in your system? If there are any, how would you fix them?
> as soon as we don't use refresh token here, the access token has big TTL.
> It can led us to the security issues if the token will be stolen.
> Also there is no mechanism of tokens invalidation

Do you feel that your skills were well tested?
> Yes