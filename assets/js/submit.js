
function submitSignup(e) {

    e.preventDefault();

    const URL = "http://localhost:8080/signup"

    let formData = new FormData(e.target);

    let data = {};

    data.member = {};
    data.emergency_contact = {};

    data.legalfirstnames = formData.get("full_name");

    data.member.firstname = formData.get("nickname");
    data.member.infix = formData.get("infix");
    data.member.lastname = formData.get("surname");
    data.member.phone = formData.get("phone");

    let bday = formData.get("birthday");

    data.date_of_birth = new Date(bday).toISOString().replace(/[.]\d+[Z]/g, "Z");

    data.address = formData.get("address");
    data.postal_code = formData.get("zipcode");
    data.city = formData.get("city_village");
    data.email = formData.get("mail");
    data.course = formData.get("education");
    data.cohort = formData.get("cohort_year");

    data.emergency_contact.firstname = formData.get("emergency_contact_first_name");
    data.emergency_contact.infix = formData.get("emergency_contact_infix");
    data.emergency_contact.lastname = formData.get("emergency_contact_surname");
    data.emergency_contact.phone = formData.get("emergency_contact_phone_number");

    data.iban = formData.get("iban");
    data.account_holder = formData.get("card_holder");

    fetch(URL, {
        method: "POST",
        body: JSON.stringify(data),
        mode: "no-cors"
    }).then((a) => {
        // TODO: redirect to success page
        console.log(a.body);
        console.log("very good very nice");
    });
}

document.getElementById("signup-form").addEventListener("submit", submitSignup);