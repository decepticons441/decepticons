'use strict';

const home = () => {
    var sendData = {
        "email": document.getElementById("email").value,
        "password": document.getElementById("password").value,
        "passwordConf": document.getElementById("passwordConf").value,
        "userName": document.getElementById("userName").value,
        "firstName": document.getElementById("firstName").value,
        "lastName": document.getElementById("lastName").value,
    };
    console.log(sendData);
    fetch("https://api.nehay.me/v1/users", {
        method: "POST",
        headers: new Headers({
            "Content-Type": "application/json"
        }),
        body: JSON.stringify(sendData),
    })
    .then((response) => {
        if (response.status == 400 || response.status == 415 || response.status == 500) {
            console.log("Error when making new user. try again.")
            console.log(response)
        }
        console.log(response);
        var tokenArr = new Array();
        tokenArr = response.headers.get("Authorization").split(" ");
        sessionStorage.setItem('bearer', tokenArr[1]);
        sessionStorage.setItem('user', JSON.stringify(response));
        // sessionStorage.setItem('user', response.body);
        window.location.href="home.html";
    })
}

document.getElementById("create-account").addEventListener("click", (e) => { 
    console.log("hello")
    e.preventDefault(); 
    home(); 
})

document.getElementById("login").addEventListener("click", (e) => { 
    e.preventDefault(); 
    window.location.href="auth.html";
})