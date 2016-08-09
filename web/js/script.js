function validateForm() {

    if ($("#password").val() != $("#password1").val()) {
        console.log("passwords don't match")
        $("#result").html("Passwords Do Not Match")
    } else {
        console.log("passwords are good")
        $("#result").html("")
    }
}

function login() {
    $.ajax({
        type: 'POST',
        url: '/login',
        data: JSON.stringify({
            Username: $('#Username').val(),
            Pass: $('#password').val()
        }),
        dataType: 'json',
        async: false,
        success: function(data) {
            console.log("Posted Data");
        },
        error: function(xhr, textStatus, error) {
            console.log(xhr.statusText);
            console.log(textStatus);
            console.log(error);
        }
    });

}

function logout() {
    $.ajax({
        type: 'GET',
        url: '/logout',
        async: false,
        success: function(data) {
            alert("Myn")
                //  value = $("#inputChat").val();
        }
    });
}
