// getting dom elements
var divSelectRoom = document.getElementById("selectRoom");
var divConsultingRoom = document.getElementById("consultingRoom");
var inputRoomNumber = document.getElementById("roomNumber");
var btnGoRoom = document.getElementById("goRoom");
var localVideo = document.getElementById("localVideo");
var remoteVideo = document.getElementById("remoteVideo");

// variables
var roomNumber;
var localStream;
var remoteStream;
var rtcPeerConnection;
var iceServers = {
    'iceServers': [
        { 'urls': 'stun:stun.services.mozilla.com' },
        { 'urls': 'stun:stun.l.google.com:19302' }
    ]
}
var streamConstraints = { audio: true, video: true };
var isCaller;

// Create socket
var socket = io();

btnGoRoom.onclick = function () {
    if (inputRoomNumber.value === '') {
        alert("Please type a room number")
    } else {
        roomNumber = inputRoomNumber.value;
        socket.emit('create or join', roomNumber);
        divSelectRoom.style = "display: none;";
        divConsultingRoom.style = "display: block;";
    }
};

// message handlers
socket.on('created', function (room) {
    navigator.mediaDevices.getUserMedia(streamConstraints).then(function (stream) {
        localStream = stream;
        localVideo.srcObject = stream;
        isCaller = true;
    }).catch(function (err) {
        console.log('An error ocurred when accessing media devices', err);
    });
});

socket.on('joined', function (room) {
    navigator.mediaDevices.getUserMedia(streamConstraints).then(function (stream) {
        localStream = stream;
        localVideo.srcObject = stream;
        socket.emit('ready', roomNumber);
    }).catch(function (err) {
        console.log('An error ocurred when accessing media devices', err);
    });
});

socket.on('candidate', function (event) {
    var candidate = new RTCIceCandidate({
        sdpMLineIndex: event.label,
        candidate: event.candidate
    });
    rtcPeerConnection.addIceCandidate(candidate);
});

socket.on('ready', function () {
    if (isCaller) {
        rtcPeerConnection = new RTCPeerConnection(iceServers);
        rtcPeerConnection.onicecandidate = onIceCandidate;
        rtcPeerConnection.ontrack = onAddStream;
        rtcPeerConnection.addTrack(localStream.getTracks()[0], localStream);
        rtcPeerConnection.addTrack(localStream.getTracks()[1], localStream);
        rtcPeerConnection.createOffer()
            .then(sessionDescription => {
                rtcPeerConnection.setLocalDescription(sessionDescription);
                socket.emit('offer', {
                    type: 'offer',
                    sdp: sessionDescription,
                    room: roomNumber
                });
            })
            .catch(error => {
                console.log(error)
            })
    }
});

socket.on('offer', function (event) {
    if (!isCaller) {
        rtcPeerConnection = new RTCPeerConnection(iceServers);
        rtcPeerConnection.onicecandidate = onIceCandidate;
        rtcPeerConnection.ontrack = onAddStream;
        rtcPeerConnection.addTrack(localStream.getTracks()[0], localStream);
        rtcPeerConnection.addTrack(localStream.getTracks()[1], localStream);
        rtcPeerConnection.setRemoteDescription(new RTCSessionDescription(event));
        rtcPeerConnection.createAnswer()
            .then(sessionDescription => {
                rtcPeerConnection.setLocalDescription(sessionDescription);
                socket.emit('answer', {
                    type: 'answer',
                    sdp: sessionDescription,
                    room: roomNumber
                });
            })
            .catch(error => {
                console.log(error)
            })
    }
});

socket.on('answer', function (event) {
    rtcPeerConnection.setRemoteDescription(new RTCSessionDescription(event));
})

// handler functions
function onIceCandidate(event) {
    if (event.candidate) {
        console.log('sending ice candidate');
        socket.emit('candidate', {
            type: 'candidate',
            label: event.candidate.sdpMLineIndex,
            id: event.candidate.sdpMid,
            candidate: event.candidate.candidate,
            room: roomNumber
        })
    }
}

function onAddStream(event) {
    remoteVideo.srcObject = event.streams[0];
    remoteStream = event.stream;
}

let home = document.getElementById('home-icon');
home.addEventListener("click", homeReroute);

function homeReroute() {
    window.location.href="home.html";
}

// https://www.w3schools.com/howto/tryit.asp?filename=tryhow_css_modal
var addModal = document.getElementById('add-pop-up');
var addBtn = document.getElementById("add-search-members");
var closeSpan = document.getElementsByClassName("close-add")[0];

document.getElementById("add-icon").addEventListener("click", function() {
    console.log("i'm here");
    var message = document.getElementById('write-message').value;
    addMessage(message);
});

function addMessage(message) {
    fetch("https://api.nehay.me/v1/channels/" + JSON.parse(sessionStorage.getItem('currChannel')).id, {
        method: "POST",
        body: JSON.stringify({
            "body": message,
        }),
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
        window.onload.call();
    })
}

addBtn.addEventListener("click", addModalPop);
closeSpan.addEventListener("click", closeModalAdd);

function addModalPop() {
    addModal.style.display = "block";
}

function closeModalAdd() {
    addModal.style.display = "none";
}
window.onload = function(event) {
    document.getElementById("channel-name").innerHTML = JSON.parse(sessionStorage.getItem('currChannel')).nameString;
    // document.getElementById("message-chat").innerHTML = "";
    const res = () => {
        fetch("https://api.nehay.me/v1/channels/" + JSON.parse(sessionStorage.getItem('currChannel')).id, {
            method: "GET",
            headers: new Headers({
                "Authorization": sessionStorage.getItem('bearer'),
            })
        }).then((messages) => {
            if (messages.status == 401 || messages.status == 500) {
                console.log("Error when getting channels. try again.");
                console.log(messages);
            }
            return messages.json();
        }).then((messages) => {
            console.log(messages);
            for (i = 0; i < messages.length; i++) {
                var messageDiv = document.createElement("div");
            
                var intro = document.createElement("div");
                var introInfo = document.createElement("p");

                introInfo.innerHTML = "CreatorID #" + messages[i].creatorID + " " + messages[i].createdAt;
                if (messages[i].creatorID == JSON.parse(sessionStorage.getItem('user')).id) {
                    var editButton = document.createElement("button");
                    editButton.className = "btn btn-primary";
                    editButton.setAttribute('id', 'editButton');
                    editButton.addEventListener("click", function() {
                        editMessage(messages.body[i].id, messages.body[i].body);
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
                } else {
                    intro.append(introInfo);
                }

                messageDiv.append(intro);
                var message = document.createElement("p");
                message.setAttribute('id', 'message');
                // console.log(messages[i].body);
                message.innerHTML = messages[i].body;
                messageDiv.append(message); 
                document.getElementById("message-chat").append(messageDiv);   
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
    window.onload.reload();
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
        fetch("https://api.nehay.me/v1/channels/" + JSON.parse(sessionStorage.getItem('currChannel')).id + "/members", {
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
        fetch("https://api.nehay.me/v1/channels/" + JSON.parse(sessionStorage.getItem('currChannel')).id + "/members", {
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
        fetch("https://api.nehay.me/v1/channels/" + JSON.parse(sessionStorage.getItem('currChannel')).id , {
            method: "PATCH",
            body: {
                "nameString": document.getElementById("name").value,
                "descriptionString": document.getElementById("description").value
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
        fetch("https://api.nehay.me/v1/channels/" + JSON.parse(sessionStorage.getItem('currChannel')).id + "/members", {
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

