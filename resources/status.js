var updateLog = function(log) {
    editor.setValue(log);
};

var updateStatus = function(status) {
    $("#status").html(status)
    if(status != "In progress") {
	// This function relies on a variable that relies on this function
	// TODO? cleanup
	clearInterval(logUpdater);
    }
}

var updateCmd = function(cmd) {
    $("#command").html(cmd)
}

var logUpdater = setInterval(function() {
    $.ajax({
	url: "/info"
    }).done(function(data) {
	updateLog(data["log"]);
	updateStatus(data["status"]);
	updateCmd(data["command"])
    });
}, 1000);

var editor;

$(function() {
    editor = CodeMirror(document.getElementById("log"), {
	lineNumbers: true,
	styleActiveLine: true,
	matchBrackets: true,
	readOnly: true
    });
});
