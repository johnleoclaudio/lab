# Wed Nov 26 06:13:48 PM PST 2025

Since yesterday, my I cant access my domain `https://leoclaudio.dev`, usually it routes directly to my empty blog website. 
- I tried installing dokploy as i found a youtube tutorial how easy it is, but to my surprise it didn't fix my problem. 
- I tried deleting all of my DNS records and just add one A record `@` and tried rerunning my dokploy. Looking at the logs, there's a bunch of 403s. After a bit of tinkering, I tried adding a new record `www` in my domain. My dokploy-traefik container errors changed to 429! So this is something

# Wed Nov 26 06:18:08 PM PST 2025
- Laptop died. 429 error above was rate limit from acme. I'll wait a bit and see if i there's an improvement
- Still the same error, I'll try to uninstall dokploy. Found this quick uninstall [steps](https://docs.dokploy.com/docs/core/uninstall) from their site 
- Done installing dokploy, add my domain for the server domain. My domain is also propagated. It should work this time around. 
- No luck
