    var ws = new WebSocket("ws://localhost/ws");
    var CurrentUser = ""
    ws.onopen = function() {
      CurrentUser = getIP()
        $("#ChatPanel").html("CONNECTED")
    };
    ws.onclose = function() {
        $("#ChatPanel").html("DISCONNECTED")
    };
    ws.onmessage = function(event) {
      var IP = getIP()
      console.log("CurrentUserIP: " + CurrentUser + "\tMessageIP: " + IP)
        if (CurrentUser == IP) {
            $('<div class="clearfix"><blockquote class="you pull-left">' + event.data + '</blockquote></div>"').appendTo('#chatbox');
        } else {
            $('<div class="clearfix"><blockquote class="me pull-right">' + event.data + '</blockquote></div>"').appendTo('#chatbox');
        }

    }

    function getIP() {
        var result = null;
        $.ajax({
            type: 'GET',
            url: '/getUser',
            async: false,
            success: function(data) {
                var obj = jQuery.parseJSON(data)
                ip = obj.IP
                result = ip
                    //  value = $("#inputChat").val();
            }
        });
        return result
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
          }
      });
    }

    $(document).ready(function() {
       $('.modal-trigger').leanModal();



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
                }
            });
        });

    });
