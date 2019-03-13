'use strict';

let login = document.querySelector("#login");
let signup = document.querySelector("#signup")

// login.addEventListener("click", window.location.href="home.html");
// signup.addEventListener("click", window.location.href="signup.html");

login.addEventListener("click", home);
signup.addEventListener("click", signUp);

function home() {
    const res = async () => {
        // const res = await fetch("abc")
        // const js = await res.json()
        // console.log(js)
        const response = await fetch("https://api.nehay.me/v1/users", {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
        if (response.status == 404) {
            console.log("You are not a user within our system/you are not authorized to use this system. Please make a new account or try again.")
        }
        console.log(response)
    }
   
    window.location.href="home.html";
}

function signUp() {
    window.location.href="signup.html";
}