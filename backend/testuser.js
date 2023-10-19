async function postData(url = "", data = {}) {
    // Default options are marked with *
    const response = await fetch(url, {
      method: "POST", // *GET, POST, PUT, DELETE, etc.
      mode: "cors", // no-cors, *cors, same-origin
      cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
      credentials: "same-origin", // include, *same-origin, omit
      headers: {
        "Content-Type": "application/json",
        // 'Content-Type': 'application/x-www-form-urlencoded',
      },
      redirect: "follow", // manual, *follow, error
      referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
      body: JSON.stringify(data), // body data type must match "Content-Type" header
    });
    return response.json(); // parses JSON response into native JavaScript objects
  }

postData("http://127.0.0.1:8080/signup",
    {
        legalfirstnames: "boben b",
        member: {
            firstname: "bob",
            infix: "de",
            lastname: "tak",
            phone: "+31612345678"
        },
        date_of_birth: "2000-10-12T00:00:00Z",
        address: "Lovensdijkstaat 16",
        postal_code: "4793RR",
        city: "Breda",
        email: "jandevries@example.org",
        course: "TI",
        cohort: "2022/2023",
        emergency_contact: {
            firstname: "greetje",
            infix: "de",
            lastname: "vries",
            phone: "+31687654321"
        },
        iban: "NL13KNAB121223232345",
        account_holder: "B. B. de Tak"
    }
).then(x => console.log(x));