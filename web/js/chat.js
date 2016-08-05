    var ws = new WebSocket("ws://localhost/ws");
    var fromMe = false;
    ws.onopen = function() {
        $("#ChatPanel").html("CONNECTED")
    };
    ws.onclose = function() {
        $("#ChatPanel").html("DISCONNECTED")
    };
    ws.onmessage = function(event) {

      if(!fromMe){
        $('<div class="clearfix"><blockquote class="you pull-right">' +  event.data + '!</blockquote></div>"').appendTo('#chatbox');
      }
    }

    $(document).ready(function() {
        var $Username = $('#inputUsername');
        $('#inputChat').val("");
    // take what's the textbox and send it off
    $('#sendMsgBtn').click( function(event) {
      fromMe = true;
      ws.send($('#inputChat').val());
      $('#inputChat').val("");
    });
    fromMe = false;

      /*  $('#sendMsgBtn').click(function() {
            $.ajax({
                type: 'GET',
                url: '/getUser',
                success: function(data) {
                    console.log("Success", data);
                    value = $("#inputChat").val();
                    ws.send($Username.val() + ": " + value)
                    value = ""
                    $('<div class="clearfix"><blockquote class="me pull-left">' + $Username.val() + ": " + value + '</blockquote></div>"').appendTo('#chatbox');

                }
            });
        });*/

        $('#SetNameBtn').click(function() {
            $.ajax({
                type: 'POST',
                url: '/storeUser',
                data: JSON.stringify({
                    Username: $Username.val()
                }),
                dataType: 'json',
                success: function(data) {
                    console.log("Posted Data");
                }
            });
        });

    });
