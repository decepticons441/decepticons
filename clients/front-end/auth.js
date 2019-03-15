'use strict';

let login = document.querySelector("#login");
let signup = document.querySelector("#signup")

// login.addEventListener("click", window.location.href="home.html");
// signup.addEventListener("click", window.location.href="signup.html");

login.addEventListener("click", home);
signup.addEventListener("click", signUp);
let sendData = {
    "email": document.querySelector("#email"),
    "password": document.querySelector("#password"),
};

function home() {
    const res = async () => {
        // const res = await fetch("abc")
        // const js = await res.json()
        // console.log(js)
        const userCheck = await fetch("https://api.nehay.me/v1/sessions", {
            method: "POST",
            body: JSON.stringify(sendData),
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
        if (userCheck.status == 404) {
            console.log("You are not a user within our system/you are not authorized to use this system. Please make a new account or try again.")
        }
        console.log(userCheck);

        const specificUser = await fetch("https://api.nehay.me/v1/users/" + userCheck.body.id, {
            method: "GET",
            headers: new Headers({
                "Content-Type": "application/json"
            })
        });
       
        if (specificUser.status == 404) {
            console.log("You are not a user within our system/you are not authorized to use this system. Please make a new account or try again.")
        }
        console.log(specificUser);

        var tokenArr = new Array();
        tokenArr = response.headers["Authorization"].split(" ");
        sessionStorage.setItem('bearer', tokenArr[1]);
        sessionStorage.setItem('user', specificUser.body);
    }
    
    window.location.href="home.html";
}

function signUp() {
    window.location.href="signup.html";
}