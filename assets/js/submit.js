const URL = "http://localhost:8080/signup" // FIXME: prod URL

// your form
const formElement = document.getElementById("signup-form");

const errorsElement = document.getElementById("errorsDiv");


// // attach event listener
formElement.addEventListener("submit", submitSignup, true);

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
            // TODO: Redirect to success page
            console.log("Great success")
            return
        }
        if (response.status === 400) {
            response.json().then((body) => {
                const errors = body["Errors"]
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
            });
        }
    })
}