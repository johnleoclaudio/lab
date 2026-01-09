# JWT Notes

- JWT can't be changed
- JWTs are not encrypted so dont store sensitive info
- JWT are stateless.. Server doesn't need to keep track of which users are logged in
- JWTs can't be revoked 

### Access Token
- stateless
- Short-lived (15m-24h)
- Irrevocable

### Refresh Tokens
- stateful 
- Long-lived (24h-60d)
- Revocable

