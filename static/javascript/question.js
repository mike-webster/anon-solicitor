const sleep = (milliseconds) => {
    return new Promise(resolve => setTimeout(resolve, milliseconds))
}

function hideAnswer(id) {
    showSuccess(id);
    sleep(2000).then(()=>{
        // fade out the div
        var s = document.getElementById("answer" + id).style;
        s.opacity = 1;
        (function fade(){(s.opacity-=.1)<0?s.display="none":setTimeout(fade,40)})();
    });
}

function showSuccess(id) {
    // replace the question with a success message but keep the same width
    // until the card fades out.
    answer = document.getElementById("answer" + id);
    width = "width: " + answer.offsetWidth + "px"; 
    answer.innerHTML = '<p class="success-para" style="' + width  + ';">Success!</p>';
    answer.setAttribute("style",width);
}

function setErrorBackground(id) {
    answer = document.getElementById("answer" + id);
    answer.style.backgroundColor = '#f2e1e1';
    answer.style.borderColor = '#e23f3f';
}

function showErrorBanner() {
    header = document.getElementsByClassName("header")[0];
    banner = document.createElement("div");
    banner.innerHTML = "<p class='error-banner'>An error occurred - please refresh the page.</p>";
    header.parentNode.insertBefore(banner, header.nextSibling);
}

function hideValidationError(id) {
    answer = document.getElementById("answer" + id);
    verrors = document.querySelectorAll('p[class="error-para"]');
    for (i = 0; i < verrors.length; i++) {
        if (verrors[i].parentNode.id == answer.id) {
            verrors[i].style.display = "none";
        }
    }
}

function hideAllQuestions() {
    var divs = document.querySelectorAll('div[id^="answer"]'), i;

    for (i = 0; i < divs.length; ++i) {
      divs[i].style.display = "none";
    }
}

function addValidationMessage(id, message) {
    answerDiv = document.getElementById("answer" + id);
    answer = document.getElementById("feedback" + id);
    newDiv = document.createElement("p");
    newDiv.className = "error-para";
    newDiv.innerHTML = message;
    answerDiv.insertBefore(newDiv, answerDiv.lastChild);
}

function makeRequest(action, path, body) {
    console.log("making request");
    var xhr = new XMLHttpRequest();
    xhr.open(action, path, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(body));
    xhr.onprogress = function(req, ev) {
        if (this.status == 200) {
            console.log("200 response");
            hideAnswer(questionID);
        } else if (this.status == 400) {
            setErrorBackground(questionID);
            body = JSON.parse(this.responseText);
            addValidationMessage(questionID, body.Content);
            console.log("400 response");
        } else if (this.status == 401) {
            // todo: hide form - show "unauth"
            console.log("401 response");
        } else if (this.status == 403) {
            // todo: idk... this shouldn't happen.
            console.log("403 response");
        } else {
            showErrorBanner();
            hideAllQuestions();
            console.log("unknown response: " + this.status);
            // TODO: show what went wrong
            var data = JSON.parse(this.responseText);
            console.log(data);
        }

    }
}

function answerQuestion(event) {
    console.log("click");
    route = "/v1/answers/{eid}/{qid}"

    e = event || window.event;
    questionID = e.id.replace("question", "")
    eid = document.getElementById("eventid").value;
    if (!eid) {
        // if we don't have an eid we're fucked
        console.log("no eid");
        return;
    }

    hideValidationError(questionID);

    path = route.replace("{eid}", eid).replace("{qid}", questionID);

    // attempt to get the feedback text by id - feedback + id
    feedback = document.getElementById("feedback" + questionID);
    if (!feedback) {
        console.log("no feedback found");
    } else {
        if (!feedback.value) {
            // todo: show validation error
            console.log("blank message");
            setErrorBackground(questionID);
            addValidationMessage(questionID, "please provide a value");
        } else {
            body = { content: feedback.value };
            makeRequest("POST", path, body);
        }
    }

    // if no element, check for the radio buttons
    radios = document.querySelectorAll('input[id^="question' + questionID + '-"]');
    if (!radios) {
        console.log("no radios");
    } else {
        console.log("radios found");
        value = "";
        for(var i = 0; i < radios.length; i++){
            console.log("radio value: ", radios[i].value);
            if (radios[i].checked) {
                value += radios[i].value;
            }
        }
        if (value.length > 0) {
            body = { content: value };
            makeRequest("POST", path, body);
        }

        return
    }
}