    var ws = new WebSocket("ws://localhost/ws");
    var CurrentUser = ""
    ws.onopen = function() {
        $("#ChatPanel").html("CONNECTED")
    };
    ws.onclose = function() {
        $("#ChatPanel").html("DISCONNECTED")
    };
    ws.onmessage = function(event) {

        console.log(CurrentUser + ": " + getUsername())
        if (CurrentUser == getUsername()) {
            $('<div class="clearfix"><blockquote class="you pull-left">' + event.data + '!</blockquote></div>"').appendTo('#chatbox');
        } else {
            $('<div class="clearfix"><blockquote class="me pull-right">' + event.data + '!</blockquote></div>"').appendTo('#chatbox');
        }

    }

    function getUsername() {
        var result = null;
        $.ajax({
            type: 'GET',
            url: '/getUser',
            async: false,
            success: function(data) {
                var obj = jQuery.parseJSON(data)
                console.log("Success", obj.IP);
                Username = obj.Username
                result = Username
                    //  value = $("#inputChat").val();
            }
        });
        return result
    }
    $(document).ready(function() {
        var $Username = $('#inputUsername');
        $('#inputChat').val("");
        // take what's the textbox and send it off
        $('#sendMsgBtn').click(function(event) {
            ws.send($('#inputChat').val());
            $('#inputChat').val("");
        });


        $('#SetNameBtn').click(function() {
            $.ajax({
                type: 'POST',
                url: '/storeUser',
                data: JSON.stringify({
                    Username: $Username.val()
                }),
                dataType: 'json',
                success: function(data) {
                    CurrentUser = $Username.val()
                    console.log("Posted Data");
                }
            });
        });

    });
