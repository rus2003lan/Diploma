<script src="/resources/labheader/js/labHeader.js"></script>


<script src="/resources/labheader/js/submitSolution.js"></script>


completedListeners = [];

(function () {
    let labHeaderWebSocket = undefined;
    function openWebSocket() {
        return new Promise(res => {
            if (labHeaderWebSocket) {
                res(labHeaderWebSocket);
                return;
            }

            let newWebSocket = new WebSocket(location.origin.replace("http", "ws") + "/academyLabHeader");

            newWebSocket.onopen = function (evt) {
                res(newWebSocket);
            };

            newWebSocket.onmessage = function (evt) {
                const labSolved = document.getElementById('notification-labsolved');
                const keepAliveMsg = evt.data === 'PONG';
                if (labSolved || keepAliveMsg) {
                    return;
                }
                document.getElementById("academyLabHeader").innerHTML = evt.data;
                animateLabHeader();

                for (const listener of completedListeners) {
                    listener();
                }
            };

            setInterval(() => {
                newWebSocket.send("PING");
            }, 5000)
        });
    }

    labHeaderWebSocket = openWebSocket();
})();

function animateLabHeader() {
    setTimeout(function() {
        const labSolved = document.getElementById('notification-labsolved');
        if (labSolved)
        {
            let cl = labSolved.classList;
            cl.remove('notification-labsolved-hidden');
            cl.add('notification-labsolved');
        }

    }, 500);
}


document.getElementById("submitSolution").addEventListener("click", function() {
    submitSolution(this.getAttribute("method"), this.getAttribute("path"), this.getAttribute("parameter"), this.getAttribute("data-answer-prompt"), this.getAttribute("data-submit-hint-callback"))
});

function submitSolution(method, path, parameter, answerPrompt, submitHintCallback) {
    var answer = prompt((answerPrompt || "Answer") + ":");
    if (answer && answer !== "") {
        var params = new URLSearchParams();
        params.append(parameter, answer);

        var XHR = new XMLHttpRequest();
        XHR.open(method, path);
        XHR.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
        XHR.addEventListener(
            "load",
            function verifyAnswer() {
                var correct = this.status === 200 && JSON.parse(this.responseText)['correct'];
                if (submitHintCallback) {
                    eval(submitHintCallback)(correct);
                } else if (!correct) {
                   alert("That answer is incorrect, please try again!");
                }
            });

        XHR.send(params.toString())
    }
}




