'use strict';

let createAccount = document.querySelector("#create-account");
let login = document.querySelector("#login")

createAccount.addEventListener("click", home);
login.addEventListener("click", auth);

let sendData = {
    "email": document.querySelector("#email"),
    "password": document.querySelector("#password"),
    "passwordConf": document.querySelector("#passwordConf"),
    "userName": document.querySelector("#userName"),
    "firstName": document.querySelector("#firstName"),
    "lastName": document.querySelector("#lastName"),
};

function home() {
    const response = await fetch("https://api.nehay.me/v1/users", {
        method: "POST",
        body: JSON.stringify(sendData),
        headers: new Headers({
            "Content-Type": "application/json"
        })
    });
    if (response.status == 400) {
        console.log("Error when making new user. try again.")
    }
    console.log(response)
    window.location.href="home.html";
}

function auth() {
    window.location.href="auth.html"
}