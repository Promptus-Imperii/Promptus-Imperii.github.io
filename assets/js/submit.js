const URL = "https://api.svpromptusimperii.nl/api/signup"
// const URL = "http://localhost:3000/api/signup"

// your form
const formElement = document.getElementById("signup-form");

const errorsElement = document.getElementById("errorsDiv");


// attach event listener after DOM has loaded to prevent issues with browsers other than Firefox.
document.addEventListener('DOMContentLoaded', function () {
    formElement.addEventListener("submit", submitSignup, true);
})

function submitSignup(e) {

    e.preventDefault();

    let formData = new FormData(e.target);

    let data = {};
    formData.forEach(function(value, key){
        data[key] = value;
    });

    if (data.altcha == null) {
        return
    }
    console.log(data)
    fetch(URL, {
        method: "POST",
        body: JSON.stringify(data),
    }).then((response) => {
        console.log(response)
        if (response.status === 200) {
            window.location.href = "/signup_success/";
            console.log("Great success")
            return
        }
        response.json().then((body) => {
            const errors = body["Errors"]
            if (errors !== null) {
                const errorsTextElement = document.getElementById("errorsText");
                errorsTextElement.innerHTML = ""; // Clear previous errors
                // Remove Tailwind "hidden" class
                if (errorsElement.classList.contains('hidden')) {
                    errorsElement.classList.remove('hidden');
                }
                errors.forEach((error) => {
                    // Make string look nice
                    error = error.charAt(0).toUpperCase() + error.slice(1);
                    error = error + '.'
                    var p = document.createElement("p");
                    p.textContent = error;
                    errorsTextElement.appendChild(p);
                });
                return
            }
            displayGenericError(body)
        }).catch(displayGenericError)
    }).catch(displayGenericError)

}

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