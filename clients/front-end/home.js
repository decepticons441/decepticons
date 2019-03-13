// 'use strict';

// let login = document.querySelector("#login-button");
// let signup = document.querySelector("#signup-button")

// login.addEventListener("click", window.location.href="home.html");
// signup.addEventListener("click", window.location.href="signup.html");


// https://www.w3schools.com/howto/tryit.asp?filename=tryhow_css_modal
var modal = document.getElementById('pop-up');
var btn = document.getElementById("create-chat");
var span = document.getElementsByClassName("close")[0];

btn.addEventListener("click", openModal);
span.addEventListener("click", closeModal);

function openModal() {
    modal.style.display = "block";
}

function closeModal() {
    modal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

var createChannel = document.getElementById("create-channel");

createChannel.addEventListener("click", chatroom);
let sendData = {
    "email": document.querySelector("#email"),
    "password": document.querySelector("#password"),
};


function chatroom() {
    const response = await fetch("https://api.nehay.me/v1/sessions", {
            method: "POST",
            body: JSON.stringify(sendData),
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
        if (response.status == 404) {
            console.log("You are not a user within our system/you are not authorized to use this system. Please make a new account or try again.")
        }
    window.location.href="chatroom.html";
}