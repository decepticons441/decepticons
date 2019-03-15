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
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        });
        if (specificChannelsBasedOnUser.status == 500) {
            console.log(specificChannelsBasedOnUser);
        }
        document.getElementById("general").addEventListener("click", function() {
            specificChatroom();
        });

        var tableBody = document.querySelector("#table");
        for (i = 1; i < specificChannelsBasedOnUser.body; i++) {
            var row = tableBody.insertRow(i);
            row.addEventListener("click", function() {
                specificChatroom(specificChannelsBasedOnUser[i].id);
            });

            var channelID = row.insertCell(0);
            channelID.innerHTML = specificChannelsBasedOnUser[i].id;
            var channelName = row.insertCell(1);
            channelName.innerHTML = specificChannelsBasedOnUser[i].nameString;
            var scope = row.insertCell(2);
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
let sendData2 = {
    "nameString": document.getElementById("name"),
    "descriptionString": document.getElementById("description"),
    "privateBool": document.getElementById("public-privacy"),
};

function chatroom() {
    const res = async () => {
        const response = await fetch("https://api.nehay.me/v1/channels", {
            method: "POST",
            body: JSON.stringify(sendData2),
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        });
        if (response.status == 400) {
            console.log("Error when making new user. try again.")
        }
        sessionStorage.setItem('currChannel', response.body.id);
    }
    window.location.href="chatroom.html";
}

function specificChatroom(id) {
    const res = async () => {
        const specificChannel = await fetch("https://api.nehay.me/v1/channels/" + id, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        });
        if (specificChannel.status == 500) {
            console.log(specificChannel);
        }
        sessionStorage.setItem('currChannel', specificChannel);
    }
    window.location.href="chatroom.html";
}