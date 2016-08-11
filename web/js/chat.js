var CurrentUser = ""
$(document).ready(function() {
    checkLogin()
    CurrentUser = getUser()
    $("#loggedinAs").html(CurrentUser)
    $('.modal-trigger').leanModal();


    var $Username = $('#inputUsername');
    $('#inputChat').val("");
    // take what's the textbox and send it off
    $('#sendMsgBtn').click(function(event) {
        ws.send($('#inputChat').val());
        $('#inputChat').val("");
    });



});

var ws = new WebSocket("ws://192.168.71.1/ws");
var CurrentUser = ""
ws.onopen = function() {
    $("#ChatPanel").html("CONNECTED")
};
ws.onclose = function() {
    $("#ChatPanel").html("DISCONNECTED")
};
ws.onmessage = function(event) {
    console.log("Current User: " + CurrentUser)
    if (CurrentUser == getUser()) {
        $('<div class="clearfix"><blockquote class="you pull-left">' + event.data + '</blockquote></div>"').appendTo('#chatbox');
     } else {
         $('<div class="clearfix"><blockquote class="me pull-right">' + event.data + '</blockquote></div>"').appendTo('#chatbox');
     }

}


function getUser() {
    var result = null;
    $.ajax({
        type: 'GET',
        url: '/getUser',
        async: false,
        success: function(data) {
            var obj = jQuery.parseJSON(data)
            Username = obj.Username
                //  value = $("#inputChat").val();
        }
    });
    return Username
}

function CreateRoom() {
    $.ajax({
        type: 'POST',
        url: '/CreateRoom',
        data: JSON.stringify({
            RoomName: $('#RoomName').val()
        }),
        dataType: 'json',
        success: function(data) {
            console.log("Posted Data");
            $('#modal1').closeModal();
        }
    });
}

function logout() {
    $.ajax({
        type: 'GET',
        url: '/logout',
        async: false,
        success: function() {
            window.location = "index.html"
        }
    });
}
