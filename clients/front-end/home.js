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
    document.getElementById("container").style.filter = "blur(5px)";
}

function closeModal() {
    modal.style.display = "none";
    document.getElementById("container").style.filter = "blur(0px)";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

window.onload = function (e) {
    let sendData = {
        "id": sessionStorage.getItem('user').id
    }
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + sessionStorage.getItem('currChannel') + "/members", {
            method: "GET",
            body: JSON.stringify(sendData),
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((specificChannelsBasedOnUser) => {
            if (specificChannelsBasedOnUser.status == 401 || specificChannelsBasedOnUser.status == 500) {
                console.log("Error when getting channels. try again.")
                console.log(specificChannelsBasedOnUser);
            }
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
        })
    }
    res();
}

document.getElementById("general").addEventListener("click", function() {
    console.log("i'm here");
    specificChatroom(-1);
});
document.getElementById("create-channel").addEventListener("click", function() {
    chatroom();
});

function chatroom() {
    let sendData2 = {
        "nameString": document.getElementById("name"),
        "descriptionString": document.getElementById("description"),
        "privateBool": document.getElementById("public-privacy"),
    };

    const res = () => {
        fetch("https://api.nehay.me/v1/channels", {
            method: "POST",
            body: JSON.stringify(sendData2),
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((channel) => {
            if (channel.status == 403 || channel.status == 500) {
                console.log("Error when adding channel. try again.")
                console.log(channel);
            }
            sessionStorage.setItem('currChannel', channel.body.id);
        })
    }
    res();
    window.location.href="chatroom.html";
}

function specificChatroom(id) {
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + id, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((channel) => {
            if (channel.status == 403 || channel.status == 500) {
                console.log("Error when adding channel. try again.")
                console.log(channel);
            }
            sessionStorage.setItem('currChannel', channel.body.id);
        })
        if (specificChannel.status == 500) {
            console.log(specificChannel);
        }
        sessionStorage.setItem('currChannel', specificChannel);
    }
    res();
    window.location.href="chatroom.html";
}