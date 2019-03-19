// 'use strict';

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

window.onload = function () {
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/members", {
            method: "GET",
            headers: new Headers({
                "Authorization": sessionStorage.getItem('bearer')
            })
        // }).then((response) => {
        //     if (response.status == 401 || response.status == 500) {
        //         console.log("Error when getting channels specific to user. try again.")
        //         console.log(response);
        //     }
        //     return response.json();
        }).then((channels) => {
            if (channels.status == 401 || channels.status == 500) {
                console.log("Error when getting channels specific to user. try again.")
                console.log(channels);
            }
            console.log(channels);
            var tableBody = document.querySelector("#table");
            for (i = 1; i < channels.body; i++) {
                var row = tableBody.insertRow(i);
                row.addEventListener("click", function() {
                    specificChatroom(channels[i].id);
                });

                var channelID = row.insertCell(0);
                channelID.innerHTML = channels[i].id;
                var channelName = row.insertCell(1);
                channelName.innerHTML = channels[i].nameString;
                var scope = row.insertCell(2);
                var privateBool = channels[i].privateBool;
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


const chatroom = () => {
    var button = "";
    if (document.getElementById("public").checked){
        button = false;
    } else {
        button = true;
    }
    let sendData2 = {
        "nameString": document.getElementById("name").value,
        "descriptionString": document.getElementById("description").value,
        "privateBool": button,
    };
    console.log(sendData2);
    fetch("https://api.nehay.me/v1/channels", {
        method: "POST",
        body: JSON.stringify(sendData2),
        headers: new Headers({
            "Content-Type": "application/json",
            "Authorization": sessionStorage.getItem('bearer')
        })
    }).then(response => {
        if (response.status == 403 || response.status == 500) {
            console.log("Error when adding channel. try again.")
            console.log(response);
            // throw new error with catch
        }
        return response.json();
    }).then((channel) => {
        console.log(channel);
        sessionStorage.setItem('currChannel', JSON.stringify(channel));
        window.location.reload();
        window.location.href="chatroom.html";
    })
}
    

function specificChatroom(id) {
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + id, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((response) => {
            if (response.status == 403 || response.status == 500) {
                console.log("Error when adding channel. try again.")
                console.log(response);
            }
            return response.json();
        }).then((channel) => {
            console.log(channel);
            sessionStorage.setItem('currChannel', JSON.stringify(channel));
            window.onload.call();
            window.location.href="chatroom.html";
        })
    }
    res();
}

document.getElementById("general").addEventListener("click", function() {
    // console.log("i'm here");
    specificChatroom(1);
});

document.getElementById("create-channel").addEventListener("click", function() {
    chatroom();
});