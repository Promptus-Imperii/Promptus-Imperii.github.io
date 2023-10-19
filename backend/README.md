# backend

## requirements
- go >= 1.21.0
- node >= 18.0.0 (optional)


## data
the backend accepts a json schema from the signup page in the following format

```json
{
    "legalfirstnames": "Johannes Hendrikus",
    "member": {
        "firstname": "Jan",
        "infix": "de",
        "lastname": "Vries",
        "phone": "+31612345678"
    },
    "date_of_birth": "2000-10-12T00:00:00Z",
    "address": "Lovensdijkstraat 16",
    "postal_code": "4793RR",
    "city": "Breda",
    "email": "jandevries@example.org",
    "course": "TI",
    "cohort": "2022/2023",
    "emergency_contact": {
        "firstname": "Greetje",
        "infix": "de",
        "lastname": "Vries",
        "phone": "+31687654321"
    },
    "iban": "NL19KNAB123456345678",
    "account_holder": "J. H. de Vries"
}
```
> [!IMPORTANT]
> member.firstname is always the name a potential members wishes to be called by. (roepnaam)

it will then validate the phone numbers, postal code and IBAN.

the server returns errors sequentially for each field that is malformatted, and assumes at least some frontend validation has been done

### testuser.js
this JS script requires node >= v18.0.0  