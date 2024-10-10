const URL = "https://api.svpromptusimperii.nl/api/email"
// const URL = "http://localhost:3000/api/email"

// your form
const formElement = document.getElementById("email-form");

const errorsElement = document.getElementById("errorsDiv");

document.querySelector('#altcha').addEventListener('statechange', (ev) => {
    // state can be: unverified, verifying, verified, error
    console.log('state:', ev.detail.state);
    if (ev.detail.state === 'verified') {
      fetch(URL, {
        method: "POST",
        // Send altcha payload
        body: JSON.stringify({
            altcha: ev.detail.payload
        }),
    }).then((response) => {
        response.text().then((body) => {
            const emailBlock = document.getElementById("email-block");
            emailBlock.innerHTML = ""; // Clear form
            var p = document.createElement("p");
            // Add email adress to it
            p.textContent = body;
            emailBlock.appendChild(p);
        }).catch(displayGenericError)
    }).catch(displayGenericError)
    }
  });

function displayGenericError(error) {
    console.log(error)
    const errorsTextElement = document.getElementById("errorsText");
    errorsTextElement.innerHTML = ""; // Clear previous errors
    // Remove Tailwind "hidden" class
    if (errorsElement.classList.contains('hidden')) {
        errorsElement.classList.remove('hidden');
    }
    var p = document.createElement("p");
    p.textContent = "Er is iets foutgegaan. Probeer het opnieuw of neem contact met ons op.";
    errorsTextElement.appendChild(p);
}