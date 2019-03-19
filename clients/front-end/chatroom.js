let home = document.getElementById('home-icon');
home.addEventListener("click", homeReroute);

function homeReroute() {
    window.location.href="home.html";
}

// https://www.w3schools.com/howto/tryit.asp?filename=tryhow_css_modal
var addModal = document.getElementById('add-pop-up');
var addBtn = document.getElementById("add-search-members");
var closeSpan = document.getElementsByClassName("close-add")[0];

addBtn.addEventListener("click", addModalPop);
closeSpan.addEventListener("click", closeModalAdd);

function addModalPop() {
    addModal.style.display = "block";
}

function closeModalAdd() {
    addModal.style.display = "none";
}
window.onload = function(event) {
    document.getElementById("channel-name").innerHTML = sessionStorage.getItem('currChannel').nameString;
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + sessionStorage.getItem('currChannel'), {
            method: "GET",
            headers: new Headers({
                "Authorization": sessionStorage.getItem('bearer'),
            })
        }).then((messages) => {
            if (messages.status == 401 || messages.status == 500) {
                console.log("Error when getting channels. try again.");
                console.log(messages);
            }
            console.log(messages);
            for (i = 1; i < messages.body; i++) {
                var messageDiv = document.createElement("div");
            
                var intro = document.createElement("div");
                var introInfo = document.createElement("p");

                fetch("https://api.nehay.me/v1/users/" + messages.body[i].creatorID, {
                    method: "GET",
                    headers: new Headers({
                        "Content-Type": "application/json",
                        "Authorization": sessionStorage.getItem('bearer')
                    })
                }).then((user) => {
                    if (user.status == 401 || user.status == 500) {
                        console.log("Error when getting a user. try again.")
                        console.log(user);
                    }
                    introInfo.innerHTML = user.body.userName + messages.body[i].createdAt;
                    var editButton = document.createElement("button");
                    editButton.className = "btn btn-primary";
                    editButton.setAttribute('id', 'editButton');
                    editButton.addEventListener("click", function() {
                        editMessage(messages.body[i].id, message.innerHTML);
                    });

                    var deleteButton = document.createElement("button");
                    deleteButton.className = "btn btn-primary";
                    deleteButton.setAttribute('id', 'deleteButton');
                    deleteButton.addEventListener("click", function() {
                        deleteMessage(messages.body[i].id);
                    });
                    intro.append(introInfo);
                    intro.append(editButton);
                    intro.append(deleteButton);

                    messageDiv.append(intro);
                    var message = document.createElement("input");
                    message.setAttribute('id', 'message');
                    message.innerHTML = messages.body[i].body;
                    messageDiv.append(message);
                })               
            }
        })
    }
    res();
}

function editMessage(messageID, message) {
    fetch("https://api.nehay.me/v1/messages/" + messageID, {
        method: "PATCH",
        body: {
            "body": message,
        },
        headers: new Headers({
            "Content-Type": "application/json",
            "Authorization": sessionStorage.getItem('bearer')
        })
    }).then((messages) => {
        if (messages.status == 403 || messages.status == 404 || messages.status == 500) {
            console.log("Error when updating message. try again.")
            console.log(messages);
        }
        console.log(messages);
    })
    // if (specificChannel.status == 500) {
    //     console.log(specificChannel);
    // }
}

function deleteMessage(messageID) {
    fetch("https://api.nehay.me/v1/messages/" + messageID, {
        method: "DELETE",
        headers: new Headers({
            "Authorization": sessionStorage.getItem('bearer')
        })
    }).then((messages) => {
        if (messages.status == 403 || messages.status == 404 || messages.status == 500) {
            console.log("Error when deleting message. try again.")
            console.log(messages);
        }
        console.log(messages);
    })
    // if (specificChannel.status == 500) {
    //     console.log(specificChannel);
    // }
    window.onload.call();
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == addModal) {
        addModal.style.display = "none";
    }
}

document.getElementById("buttonSearch").addEventListener("click", search);

var addMembers= document.getElementById("add-members-done");
addMembers.addEventListener("click", addNewMembersDone);

function addNewMembersDone() {
    window.location.href="chatroom.html";
}

function search() {
    const res = () => {
        fetch("https://api.nehay.me/v1/users?q=" + document.getElementById("search"), {
            method: "GET",
            headers: new Headers({
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((users) => {
            if (users.status == 401) {
                console.log("Error when searching for users. try again.")
                console.log(users);
            }
            console.log(users);
            var tableBody = document.querySelector("#results-users");
            for (i = 1; i < users.body; i++) {
                var row = tableBody.insertRow(i);
                // var photo = row.insertCell(0);
                // photo.innerHTML = users.body.photoURL;

                var username = row.insertCell(0);
                username.innerHTML = users.body[i].userName;

                var icon = row.insertCell(1);
                var plusIcon = document.createElement("icon");
                plusIcon.className = "fas fa-plus";
                plusIcon.setAttribute("id", "plus-icon");
                plusIcon.addEventListener("click", function() {
                    addMember(users.body[i].id);
                });
                icon.innerHTML = plusIcon;
            }
        })
    }

    res();
}

function addMember(memberID) {
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + sessionStorage.getItem('currChannel') + "/members", {
            method: "POST",
            body: {
                "id": memberID
            },
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((users) => {
            if (users.status == 401 || users.status == 403 || users.status == 404 || users.status == 500) {
                console.log("Error when adding members. try again.")
                console.log(users);
            }
            console.log(users);
            var tableBody = document.querySelector("#current-members");
            var row = tableBody.insertRow(i);
            // var photo = row.insertCell(0);
            // photo.innerHTML = users.body.photoURL;

            var username = row.insertCell(0);
            username.innerHTML = users.body.userName;

            var icon = row.insertCell(1);
            var minusIcon = document.createElement("icon");
            minusIcon.className = "fas fa-minus";
            minusIcon.setAttribute("id", "minus-icon");
            minusIcon.addEventListener("click", function () {
                removeMember(memberID);
            });
            icon.innerHTML = minusIcon;
        })
        
    }
    res();
}

function removeMember(memberID) {
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + sessionStorage.getItem('currChannel') + "/members", {
            method: "DELETE",
            body: {
                "id": memberID
            },
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((users) => {
            if (users.status == 401 || users.status == 403 || users.status == 404 || users.status == 500) {
                console.log("Error when adding members. try again.")
                console.log(users);
            }
            console.log(users);
        })
    }
    res();
}

// https://www.w3schools.com/howto/tryit.asp?filename=tryhow_css_modal
var modal = document.getElementById('close-pop-up');
var btn = document.getElementById("update-chat");
var span = document.getElementsByClassName("close-update")[0];

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

var updateChannel = document.getElementById("members-done");

updateChannel.addEventListener("click", chatroom);

function chatroom() {
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + sessionStorage.getItem('currChannel') , {
            method: "PATCH",
            body: {
                "nameString": document.getElementById("name"),
                "descriptionString": document.getElementById("description")
            },
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((channel) => {
            if (channel.status == 403 || channel.status == 404 || channel.status == 500) {
                console.log("Error when updating channel. try again.")
                console.log(channel);
            }
            console.log(channel);
        })
    }
    res();
    window.location.href="chatroom.html";
}



// https://www.w3schools.com/howto/tryit.asp?filename=tryhow_css_modal
var deleteCheckModal = document.getElementById('delete-warning');
var deleteBtn = document.getElementById("delete-chat");
var span = document.getElementsByClassName("close-update")[0];

deleteBtn.addEventListener("click", deleteWarningPop);
span.addEventListener("click", closeModal);

function deleteWarningPop() {
    deleteCheckModal.style.display = "block";
}

function closeModal() {
    deleteCheckModal.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == deleteCheckModal) {
        deleteCheckModal.style.display = "none";
    }
}

var deleteChannelYES = document.getElementById("members-yes-done");

deleteChannelYES.addEventListener("click", deleteRerouteYES);

function deleteRerouteYES() {
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + sessionStorage.getItem('currChannel') + "/members", {
            method: "DELETE",
            headers: new Headers({
                "Authorization": sessionStorage.getItem('bearer')
            })
        }).then((user) => {
            if (user.status == 401 || user.status == 403 || user.status == 404 || users.status == 500) {
                console.log("Error when deleting channel member. try again.")
                console.log(user);
            }
        })
    }
    res();
    window.location.href="home.html";
}

var deleteChannelNO = document.getElementById("members-no-done");
deleteChannelNO.addEventListener("click", deleteRerouteNO);

function deleteRerouteNO() {
    window.location.href="chatroom.html";
}

