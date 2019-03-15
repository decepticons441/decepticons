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
let sendData = {
    "id": sessionStorage.getItem('user').id
}

window.onload = function (e) {
    const res = async () => {
        const specificChannelsBasedOnUser = await fetch("https://api.nehay.me/v1/channels/members", {
            method: "GET",
            body: JSON.stringify(sendData),
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
        if (specificChannelsBasedOnUser.status == 500) {
            console.log(specificChannelsBasedOnUser);
        }
        var tableBody = document.querySelector("#table");
        for (i = 1; i < specificChannelsBasedOnUser.body; i++) {
            var row = tableBody.insertRow(i);

            var channelName = row.insertCell(0);
            channelName.innerHTML = specificChannelsBasedOnUser[i].nameString;
            var scope = row.insertCell(1);
            var privateBool = specificChannelsBasedOnUser[i].privateBool;
            if (privateBool) {
                scope.innerHTML = "private";
            } else {
                scope.innerHTML = "public";
            }
        }

    }
}

var createChannel = document.getElementById("create-channel");

createChannel.addEventListener("click", chatroom);

function chatroom() {
    
    window.location.href="chatroom.html";
}