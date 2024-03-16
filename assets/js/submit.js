const URL = "http://localhost:8080/signup" // FIXME: prod URL

// your form
var form = document.getElementById("signup-form");

// // attach event listener
form.addEventListener("submit", submitSignup, true);

function submitSignup(e) {

    e.preventDefault();

    let formData = new FormData(e.target);

    var data = {};
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
        mode: "no-cors"
    }).then((a) => {
        // TODO: redirect to success page
        console.log(a.body);
    });
}