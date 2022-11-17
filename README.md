# Description
The authentication service has 3 main responsibilities. They are:

### 1. Create an Account
For a user to be created, the following information needs to be supplied:
* First Name
* Last Name
* Email
* Username
* Password

### 2. Authorize a User
A user can be authenticated by passing in either an email or a username as well as a password.

If the user has been verified a Paseto token will be created that can authenticate a user.

If the user has not been verified, an error will be sent back.

### 3. Verify a Paseto Token
The authentication service will also verify Paseto tokens. This will allow the broker service to verify incoming requests to see whether or not the request needs to be validated.