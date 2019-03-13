'use strict';

let login = document.querySelector("#login");
let signup = document.querySelector("#signup")

// login.addEventListener("click", window.location.href="home.html");
// signup.addEventListener("click", window.location.href="signup.html");

login.addEventListener("click", home);
signup.addEventListener("click", signUp);

function home() {
    window.location.href="home.html";
}

function signUp() {
    window.location.href="signup.html";
}