### Start / Restart:
Bash: `make restart`

### Start / Restart in docker:
Bash: `make docker`

## Roles
### Create role:
```json
{
  "method": "create_role",
  "data": {
    "type": "sto",
    "sto_id": 10,
    "permissions": [
      {
        "microservice": "go-auth",
        "method": "test",
        "required_params": [
          {
            "param": "sto_id",
            "values": ["10"],
            "all": false
          },
          {
            "param":  "user.internal_id",
            "values": ["$.internal_id"],
            "all":    false
          },
          {
            "param":  "user.first_name",
            "values": ["$.data.first_name"],
            "all":    false
          }
        ],
        "restricted_params": [
          {
            "param": "param1",
            "values": null,
            "all": true
          }
        ]
      }
    ],
    "data": {
      "alias": "super-admin",
      "name": "Super Admin",
      "any_other_field": "blabla"
    }
  }
}
```
`"param":  "user.internal_id"`: path to parameter in the method  
`"values": ["$.internal_id"]`: path to parameter in the user object

### Update role:
```json
{
  "method": "update_roles",
  "data": {
    "Select": {"sto_id":  10},
    "Data": {
      "type": "sto",
      "sto_id": 10,
      "permissions": [
        {
          "microservice": "go-auth",
          "method": "test",
          "required_params": [],
          "restricted_params": []
        }
      ],
      "data": {
        "alias": "super-admin",
        "name": "Super Admin 2",
        "any_other_field": "blabla"
      }
    }
  }
}
```

### Get roles:
```json
{
  "method": "get_roles",
  "data": {
     "sto_id": "1"
  }
}
```

### Delete roles:
```json
{
  "method": "delete_roles",
  "data": {
     "sto_id": "1" 
  }
}
```

### Attach role:
```json
{
  "method": "attach_role",
  "data": {
    "user_id": "19fc7d6f-c03b-4d0b-97d9-8660362c8930",
    "role_id": "de1538cd-24f0-43cd-b264-c5f6eb6a1e46"
  }
}
```

### Dettach role:
```json
{
  "method": "detach_role",
  "data": {
    "user_id": "19fc7d6f-c03b-4d0b-97d9-8660362c8930",
    "role_id": "de1538cd-24f0-43cd-b264-c5f6eb6a1e46"
  }
}
```

## Users
### Sign Up | Create user
```json
{
  "method": "sign_up",
  "data": {
    "email": "john@example.com",
    "phone": "+11234567890",
    "otp_code": "1234",
    "password": "securePassword123",
    "data" : {
      "first_name": "John",
      "last_name":  "Doe"
    }    
  }
}
```

### Update user
```json
{
  "method": "update_user",
  "data": {
    "Select": {
      "email": "currentEmail@example.com"
    },
    "Data": {
      "email": "newemail@example.com",
      "phone": "1234567890",
      "password": "newPassword",
      "data" : {
        "first_name": "John",
        "last_name":  "Doe"
      }
    }
  }
}
```

### Get users
```json
{
  "method": "get_users",
  "data": {
    "internal_id": "19fc7d6f-c03b-4d0b-97d9-8660362c8930"
  }
}
```

### Delete users
```json
{
  "method": "delete_users",
  "data": {
    "internal_id": "19fc7d6f-c03b-4d0b-97d9-8660362c8930"
  }
}
```

## Auth
### Send OTP code:
```json
{
  "method": "send_verify_code",
  "data": {
    "phone": "+1234567890",
    "template": "test",
    "variables": {
        "var1": "value1"
    }
  }
}
```

### Send fake OTP code (no SMS is sent and the OTP code is set to "1111"):
```json
{
  "method": "send_verify_code",
  "data": {
    "phone": "+1234567890",
    "fake": "mstfiqalx"
  }
}
```

### Sign In
```json
{
  "method": "sign_in",
  "data": {
    "login": "username_or_email",
    "password": "yourpassword"
  }
}
```

### Check permission example
```json
{
  "method": "test_cred",
  "data":{
    "sto_id":      "",
    "param2":      "",
    "internal_id": "",
    "first_name":  "",
    "token":       "$token"
  }
}
```

### Check OTP code
```json
{
    "method":"check_otp_code",
    "data":{
        "otp_code": "1111",
        "phone": "+1234567890"
    }
}
```

### Send OTP code for reset password (no SMS is sent and the OTP code is set to "1111")
```json
{
  "method": "send_reset_password_verify_code",
  "data": {
    "phone": "+1234567890",
    "fake": "mstfiqalx"
  }
}
```
### Send OTP code for reset password
```json
{
  "method": "send_reset_password_verify_code",
  "data": {
    "phone": "+1234567890",
    "template": "test",
    "variables": {
      "var1": "value1"
    }
  }
}
```
### Restore password
```json
{
  "method": "restore_password",
  "data": {
    "phone": "+1234567890",
    "otp_code": "123456",
    "password": "11115"
  }
}
```

`data`: method parameters expect token  
`$token`: token got from user sign_in method 

## Error codes:
| Error Code | Description                                                                            |
|:-----------|----------------------------------------------------------------------------------------|
| DFE_01     | Invalid data format error. The request data format does not match the expected format. |
| RFE_02     | Restricted field error. The request data includes a restricted field.                  |
| VLE_03     | Validation error. The request data fails validation rules.                             |
| UEE_04     | User exists error. A user with the provided email or phone already exists.             |
| OPE_05     | OTP error. The provided OTP code is invalid or expired.                                |
| SVE_06     | Server error. An unexpected server error occurred.                                     |


