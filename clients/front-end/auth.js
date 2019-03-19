'use strict';

// import { BeginSession } from '../../servers/gateway/sessions/session.go';
// var token = BeginSession();

const response = () => {
    let sendData = {
        "email": document.querySelector("#email").value,
        "password": document.querySelector("#password").value,
    };

    fetch("https://api.nehay.me/v1/sessions", {
        method: "POST",
        body: JSON.stringify(sendData),
        headers: new Headers({
            "Content-Type": "application/json",
            "Authorization": sessionStorage.getItem('bearer')       // won't work
        })
    }).then(userCheck => {
        if (userCheck.status == 400 || userCheck.status == 401 || userCheck.status == 415 || userCheck.status == 500) {
            console.log("Error when getting specific user. try again.")
            console.log(userCheck);
        }
        console.log(userCheck);

        fetch("https://api.nehay.me/v1/users/" + userCheck.body.id, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json",
                "Authorization": sessionStorage.getItem('bearer')       // won't work
            })
        }).then(specificUser => {
            if (specificUser.status == 400 || specificUser.status == 401 || specificUser.status == 500) {
                console.log("Error when getting specific user. try again.")
                console.log(specificUser);
            }
            var tokenArr = new Array();
            tokenArr = response.headers["Authorization"].split(" ");
            sessionStorage.setItem('bearer', tokenArr[1]);
            console.log(specificUser.body);
            sessionStorage.setItem('user', specificUser.body);
        })
    })
}

document.getElementById("login").addEventListener("click", (e) => { 
    e.preventDefault(); 
    response();
    // window.location.href="home.html"; 
})

document.getElementById("signup").addEventListener("click", (e) => { 
    e.preventDefault(); 
    window.location.href="signup.html";
})