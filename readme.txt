Step-by-step to run application
1. go to directory mdntest
2. build application using command
    go build
3. run using command
    ./myapp

List endpoint
1. register
    <base_url>/register
    json body input example: 
    {
        "Username": "arie", 
        "Email": "arie@example.com", 
        "Password": "123456"
    }
2. login
    <base_url>/login
    json body input example: 
    {
        "username_or_email": "arie", 
        "password": "123456"
    }
3. change password
    <base_url>/change-password
    json body input example: 
    {
        "old_password": "123456", 
        "new_password": "1234567"
    }