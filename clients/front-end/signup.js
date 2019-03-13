'use strict';

let createAccount = document.querySelector("#create-account");
let login = document.querySelector("#login")

createAccount.addEventListener("click", home);
login.addEventListener("click", auth);

function home() {
    window.location.href="home.html";
}

function auth() {
    window.location.href="auth.html"
}