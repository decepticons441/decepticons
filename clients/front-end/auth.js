'use strict';

// import { BeginSession } from '../../servers/gateway/sessions/session.go';
// var token = BeginSession();

const response = () => {
    let sendData = {
        "email": document.querySelector("#email").value,
        "password": document.querySelector("#password").value,
    };

    console.log(sendData);

    fetch("https://api.nehay.me/v1/sessions", {
        method: "POST",
        body: JSON.stringify(sendData),
        headers: new Headers({
            "Content-Type": "application/json",
        })
    // }).then((response) => { 
    //     if (response.status == 400 || response.status == 401 || response.status == 415 || response.status == 500) {
    //         console.log("Error when getting specific user. try again.")
    //         console.log(response);
    //     }
    //     // console.log(response);
    //     return response.json();
    }).then(userCheck => {
        if (userCheck.status == 400 || userCheck.status == 401 || userCheck.status == 415 || userCheck.status == 500) {
            console.log("Error when getting specific user. try again.")
            console.log(userCheck);
        }

        var tokenArr = new Array();
        console.log(userCheck.headers.get("Authorization"));
        tokenArr = userCheck.headers.get("Authorization").split(" ");
        console.log(tokenArr[1]);
        sessionStorage.setItem('bearer', tokenArr[1]);
        // console.log(userCheck.body);
        sessionStorage.setItem('user', JSON.stringify(userCheck.json()));
    })
}

document.getElementById("login").addEventListener("click", (e) => { 
    e.preventDefault(); 
    response();
    window.location.href="home.html"; 
})

document.getElementById("signup").addEventListener("click", (e) => { 
    e.preventDefault(); 
    window.location.href="signup.html";
})