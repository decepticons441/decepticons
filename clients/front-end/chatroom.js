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

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
    if (event.target == addModal) {
        addModal.style.display = "none";
    }
}

var addMembers= document.getElementById("add-members");

addMembers.addEventListener("click", addNewMembers);

function addNewMembers() {
    window.location.href="chatroom.html";
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
    window.location.href="home.html";
}

var deleteChannelNO = document.getElementById("members-no-done");
deleteChannelNO.addEventListener("click", deleteRerouteNO);

function deleteRerouteNO() {
    window.location.href="chatroom.html";
}