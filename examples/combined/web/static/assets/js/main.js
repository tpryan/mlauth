var auth_language = false;
var auth_vision = false;
var auth_speech = false;

var tokens = {};
tokens.language = "";
tokens.vision = "";
tokens.speech = "";

var auth = {};
auth.language = false;
auth.vision = false;
auth.speech = false;

document.addEventListener('DOMContentLoaded', function() {
    getKeys();
    document.querySelector("#text-language").addEventListener("keyup", checkLanguage);   
    document.querySelector("#file-speech").addEventListener("change", checkSpeech);   
    document.querySelector("#file-vision").addEventListener("change", checkVision);   
});





function getKeys(){

    var xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var keys = JSON.parse(this.responseText);
            setKeys(keys);
        } else if (this.status == 401) {
            console.log("Couldn't get keys");
        }
    };
    xmlhttp.onerror = function(){
        console.log("Couldn't get keys");
    };
    xmlhttp.open("GET", "/keys", true);
    xmlhttp.send();
    
}

function setKeys(keys){
    document.querySelector(".vision .key").innerHTML = keys.vision;
    document.querySelector(".speech .key").innerHTML = keys.speech;

    var lang = keys.language.split(",");
    document.querySelector(".language .key").innerHTML = lang[0];
    document.querySelector(".language .positive").innerHTML = lang[1];
}

function postAPI(callType, formData){
    var xmlhttp = new XMLHttpRequest();

    xmlhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var auth_resp = JSON.parse(this.responseText);
            if (auth_resp.auth){
                auth[callType] = true;
                tokens[callType] = auth_resp.token;
                passAuth("." + callType);
            } else {
                failAuth("." + callType);
            }
        } else if (this.status == 401) {
            failAuth("." + callType);
        }
    };
    xmlhttp.onerror = function(){
        failAuth("." + callType);
    };
    xmlhttp.open("POST", "/auth/" + callType , true);
    xmlhttp.send(formData);
}

function checkLanguage(e){

    if (e.keyCode != 13) {
        return;
    }

    var sentence = document.querySelector("#text-language").value;
    var formData = new FormData();
    formData.append("sentence", sentence);
    postAPI("language", formData);

}

function checkVision(e){
    var file = document.getElementById('file-vision').files[0];
    var formData = new FormData();
    formData.append("picture", file);
    postAPI("vision", formData);
}

function checkSpeech(e){
    var file = document.getElementById('file-speech').files[0];
    var formData = new FormData();
    formData.append("audio", file);
    postAPI("speech", formData);
}



function failAuth(target){
    document.querySelector(target).classList.add("item-failed");
    document.querySelector(target + " .alert").style.visibility = "visible";
    document.querySelector(target + " .icon").src = "assets/img/icon-lock-red.png";
}

function passAuth(target){
    var content = document.querySelector(target);
    content.classList.add("item-unlocked");
    content.classList.remove("item-failed");
    document.querySelector(target + " .icon").src = "assets/img/icon-unlock.png";
    document.querySelector(target + " .alert").style.visibility = "hidden";
    checkSecret();
}

function checkSecret(){
    if (auth.language && auth.vision && auth.speech){
        var formData = new FormData();
        formData.append("token_vision", tokens.vision);
        formData.append("token_speech", tokens.speech);
        formData.append("token_language", tokens.language);

        var xmlhttp = new XMLHttpRequest();
        xmlhttp.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                var token_resp = JSON.parse(this.responseText);
                if (token_resp.result){
                    writeSecret(token_resp.secret);
                } else {
                    failSecret();
                }
            } else if (this.status == 401) {
                failSecret();
            }
        };
        xmlhttp.onerror = function(err){
            failSecret();
        };
        xmlhttp.open("POST", "/auth/secret", true);
        xmlhttp.send(formData);
    }
}

function writeSecret(secret){
    document.querySelector(".secret-content").innerHTML = secret;
    document.querySelector(".secret").style.display = "block";
}
