## saiAuth
### Run in Docker
`make up`

### Run as standalone application
`microservices/saiAuth/build/sai-auth` 

### API
#### Register
- request:

 curl --location --request GET 'http://localhost:8800/register' \
 &emsp;    --header 'Token: SomeToken' \
 &emsp;    --header 'Content-Type: application/json' \
 &emsp;    --data-raw '{"key":"user","password":"12345"}'`

- response: '{\"Status\":\"Ok\"}'

#### Login
- request:

'curl --location --request GET 'http://localhost:8800/login' \
&emsp;    --header 'Token: SomeToken' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"key":"user","password":"12345"}''

- response:  '{"token":"3rwef2wef2ff23g2g","User":{"_id":"df22f23r435d","key":"user","roles":["User"]}}'

#### Access 
- request:

'curl --location --request GET 'http://localhost:8800/access' \
&emsp;    --header 'Token: 7ead9e6a0977a3bd33ffec382de1558c1ec139bf704ae19cc853094391afd145' \
&emsp;    --header 'Content-Type: application/json' \
&emsp;    --data-raw '{"collection":"users", "method": "get" }''
&emsp;    - response 

- response: 'true'


 
## Profiling
 `host:port/debug/pprof`